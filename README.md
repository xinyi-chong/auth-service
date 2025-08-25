# Auth Service

The **Auth Service** is a microservice responsible for handling user authentication and authorization. It provides APIs for user registration, login, password management, and token-based authentication using JWT. This service is designed to be lightweight, modular, and easily integrable into larger systems.

---

## 🚀 Features

- User registration and login
- Password hashing
- JWT-based authentication (access & refresh tokens)
- Custom error handling
- Postgres database support via GORM
- Swagger API documentation

---

## 📂 Project Structure

```text
auth-service/
├── cmd/                # Application entrypoints
├── configs/            # Configuration management (YAML)
├── internal/           # Core business logic
│   ├── api/            # API routes and handlers
│   ├── config/         # Internal configuration loading
│   ├── middleware/     # Gin middlewares (auth, logging, recovery, etc.)
│   ├── shared/         # Shared utilities and helpers
│   ├── user/           # User domain logic
│   └── auth/           # Auth domain (DTOs, services, controllers)
├── db/                 # Database management
│   └── migrations/     # Database schema migrations
├── pkg/                # Reusable packages
│   ├── error/          # Centralized error handling
│   ├── filters/        # Query filters and utilities
│   ├── i18n/           # Internationalization support
│   ├── jwt/            # JWT utilities
│   ├── logger/         # Structured logging
│   ├── redis/          # Redis cache integration
│   └── success/        # Standardized success responses
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

---
