DO $$
BEGIN

  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'auth_service') THEN
    EXECUTE format('CREATE ROLE auth_service WITH LOGIN PASSWORD %L',
      current_setting('app.db_auth_service_pwd'));

    ALTER ROLE auth_service SET search_path = auth;
    ALTER ROLE auth_service CONNECTION LIMIT 30;

    COMMENT ON ROLE auth_service IS 'Application service account with write access';
  END IF;

  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'auth_ro') THEN
    EXECUTE format('CREATE ROLE auth_ro WITH LOGIN PASSWORD %L NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT NOREPLICATION',
      current_setting('app.db_auth_ro_pwd'));

    ALTER ROLE auth_ro SET search_path = auth;

    COMMENT ON ROLE auth_ro IS 'Read-only monitoring account';
  END IF;

  GRANT USAGE, CREATE ON SCHEMA auth TO auth_service;
  GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA auth TO auth_service;
  GRANT USAGE ON SCHEMA auth TO auth_ro;
  GRANT SELECT ON ALL TABLES IN SCHEMA auth TO auth_ro;

  RAISE NOTICE 'Roles initialized successfully';
  EXCEPTION WHEN others THEN RAISE EXCEPTION 'Role initialization failed: %', SQLERRM;

END
$$;