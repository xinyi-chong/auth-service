BEGIN;

CREATE
EXTENSION IF NOT EXISTS "pg_uuidv7";

CREATE SCHEMA auth;

CREATE TABLE auth.users
(
    id                   UUID PRIMARY KEY     DEFAULT uuid_generate_v7(),
    username             VARCHAR(50) UNIQUE,
    email                VARCHAR(255) UNIQUE,
    email_verified       BOOLEAN     NOT NULL DEFAULT FALSE,
    password_hash        TEXT, -- Nullable for social logins
    last_login           TIMESTAMPTZ,
    is_active            BOOLEAN     NOT NULL DEFAULT TRUE,
    password_changed_at  TIMESTAMPTZ,
    account_locked_until TIMESTAMPTZ,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ,
    CONSTRAINT valid_email CHECK (
        email IS NULL OR email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
) , CONSTRAINT email_or_username_required CHECK (
    (email IS NOT NULL) OR (username IS NOT NULL)
)
);

CREATE INDEX idx_users_email ON auth.users (LOWER(email)) WHERE email IS NOT NULL;
CREATE INDEX idx_users_active ON auth.users (id) WHERE deleted_at IS NULL;

CREATE TABLE auth.auth_methods
(
    id          UUID PRIMARY KEY     DEFAULT uuid_generate_v7(),
    user_id     UUID        NOT NULL REFERENCES auth.users (id) ON DELETE CASCADE,
    method_type VARCHAR(20) NOT NULL CHECK (
        method_type IN ('password', 'oauth2_google', 'oauth2_github', 'oauth2_apple', 'webauthn')
        ),
    provider_id VARCHAR(255), -- For OAuth
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, method_type, provider_id)
);

CREATE INDEX idx_auth_methods_user ON auth.auth_methods (user_id);
CREATE INDEX idx_auth_methods_user_type ON auth.auth_methods (user_id, method_type);

CREATE TABLE auth.security_logs
(
    id      UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID        REFERENCES auth.users (id) ON DELETE SET NULL,
    action  VARCHAR(50) NOT NULL CHECK (
        action IN ('login', 'login_failed', 'logout', 'token_refresh',
        'password_change', 'email_change')
) ,
    status VARCHAR(20) NOT NULL CHECK (
        status IN ('success', 'failed', 'revoked', 'expired')
    ),
    ip_address INET,
    user_agent TEXT,
    device_fingerprint VARCHAR(255),
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE TABLE auth.security_logs_default PARTITION OF auth.security_logs DEFAULT;

CREATE INDEX idx_security_logs_user ON auth.security_logs (user_id);
CREATE INDEX idx_security_logs_action ON auth.security_logs (action);
CREATE INDEX idx_security_logs_created ON auth.security_logs (created_at);

CREATE TABLE auth.oauth_clients
(
    id                 UUID PRIMARY KEY      DEFAULT uuid_generate_v7(),
    client_id          VARCHAR(100) NOT NULL,
    client_secret_hash TEXT         NOT NULL,
    name               VARCHAR(255) NOT NULL,
    description        TEXT,
    redirect_uris      TEXT[] NOT NULL,
    scopes             VARCHAR(255)[] NOT NULL DEFAULT '{}'::VARCHAR[],
    is_confidential    BOOLEAN      NOT NULL DEFAULT TRUE,
    is_active          BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (client_id)
);

CREATE INDEX idx_oauth_clients_active ON auth.oauth_clients (id) WHERE is_active = TRUE;

CREATE TABLE auth.roles
(
    id          UUID PRIMARY KEY     DEFAULT uuid_generate_v7(),
    name        VARCHAR(50) NOT NULL,
    description TEXT,
    is_system   BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (name)
);

CREATE TABLE auth.permissions
(
    id          UUID PRIMARY KEY      DEFAULT uuid_generate_v7(),
    code        VARCHAR(100) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (code)
);

CREATE TABLE auth.role_permissions
(
    role_id       UUID        NOT NULL REFERENCES auth.roles (id) ON DELETE CASCADE,
    permission_id UUID        NOT NULL REFERENCES auth.permissions (id) ON DELETE CASCADE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE auth.user_roles
(
    user_id    UUID        NOT NULL REFERENCES auth.users (id) ON DELETE CASCADE,
    role_id    UUID        NOT NULL REFERENCES auth.roles (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

COMMIT;