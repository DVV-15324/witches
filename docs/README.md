
[Tiếng Việt](./README.vi.md) | [English](./README.md)

<div align="center">

<img src="../logo/logo.png" alt="Witches Logo" width="250"/>

### Fast & Scalable Golang Backend

<p>
  REST API built with <b>Go</b>, designed for performance,
  clean architecture, and modern backend development.
</p>

<p>
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go">
  <img src="https://img.shields.io/badge/Gin-Web_Framework-008ECF?style=for-the-badge">
  <img src="https://img.shields.io/badge/SQL-Database-orange?style=for-the-badge&logo=sql">
  <img src="https://img.shields.io/badge/Swagger-API_Docs-green?style=for-the-badge">
  <img src="https://img.shields.io/badge/GORM-ORM-25A162?style=for-the-badge&logo=gorm">
</p>

</div>

---

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Database & Redis](#database--redis)
- [Migrations](#migrations)
- [Running the Application](#running-the-application)
- [Project Structure](#project-structure)
- [Clean Architecture](#clean-architecture)
- [Technologies](#technologies)
- [Best Practices](#best-practices)
- [License](#license)

---

## Features

| Feature | Description |
|---------|-------------|
| **JWT Verification** | Secure token-based verification with Refresh Token |
| **Swagger Auto-Generate** | Auto-generated API documentation |
| **Database Migration** | Version control for database schema |
| **Hash Utilities** | Secure password hashing |
| **UID Masking** | Hide sensitive user IDs |
| **Clean Architecture** | Scalable and maintainable code structure |
| **Redis Support** | Built-in caching |
| **Rate Limiting** | Protect your APIs |
| **GORM Support** | Powerful ORM for database operations |

---

## Quick Start

### 1. Installation

```bash
# Install Witches CLI
go install github.com/DVV-15324/witches@latest

# Verify installation
witches version
```

### 2. Create New Project

```bash
# Create a new project
witches create example

# Navigate to project directory
cd example
```

#### Generated: `witches.env`

```env
# SERVER CONFIGURATION

APP_PORT=8080
APP_HOST=localhost

# DATABASE CONFIGURATION

DB_DRIVER=%s
DB_USER=root
DB_HOST=localhost
DB_PORT=3306
DB_NAME=your_database
DB_PASSWORD=your_password


# REDIS CONFIGURATION

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=


# TOKEN EXPIRATION

ACCESS_TOKEN_TTL=900
REFRESH_TOKEN_TTL=604800


# SESSION SETTINGS

SESSION_TTL=604800
IDLE_TIMEOUT=1800


# BLACKLIST SETTINGS

REVOKED_TTL=300


# RATE LIMIT

RATE_LIMIT_PERIOD=60
RATE_LIMIT_MAX=100 
```

---

## Configuration

### Environment Variables Reference

| Variable | Type | Description |
|----------|------|-------------|
| `APP_PORT` | string | Application server port |
| `APP_HOST` | string | Application server host |
| `DB_DRIVER` | string | Database driver (mysql, postgres, mssql) |
| `DB_HOST` | string | Database server host |
| `DB_PORT` | string | Database server port |
| `DB_PASSWORD` | string | Database password |
| `DB_NAME` | string | Database name |
| `REDIS_HOST` | string | Redis server host |
| `REDIS_PORT` | string | Redis server port |
| `REDIS_PASSWORD` | string | Redis password |
| `ACCESS_TOKEN_TTL` | int64 | Access token expiration (seconds) |
| `REFRESH_TOKEN_TTL` | int64 | Refresh token expiration (seconds) |
| `SESSION_TTL` | int64 | Session expiration (seconds) |
| `REVOKED_TTL` | int64 | Revoked token cache TTL (seconds) |
| `IDLE_TIMEOUT` | int64 | HTTP idle timeout (seconds) |
| `RATE_LIMIT_PERIOD` | int64 | Rate limit time window (seconds) |
| `RATE_LIMIT_MAX` | int64 | Max requests per period |

### Edit Configuration

Edit `witches.env` with your actual settings:


```env
# SERVER CONFIGURATION

APP_PORT=8080
APP_HOST=localhost

# DATABASE CONFIGURATION

DB_DRIVER=mysql
DB_USER=root
DB_HOST=localhost
DB_PORT=3307
DB_NAME=test
DB_PASSWORD=123


# REDIS CONFIGURATION

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=


# TOKEN EXPIRATION

ACCESS_TOKEN_TTL=900
REFRESH_TOKEN_TTL=604800


# SESSION SETTINGS

SESSION_TTL=604800
IDLE_TIMEOUT=1800


# BLACKLIST SETTINGS

REVOKED_TTL=300


# RATE LIMIT

RATE_LIMIT_PERIOD=60
RATE_LIMIT_MAX=100 
```
---

## Database & Redis

### Use Existing Your Database

> **Supported:** MySQL, PostgreSQL, MSSQL

```bash
witches database up
```

This command will:
- Auto-generate `DB_URL` for MySQL and PostgreSQL
- Prompt you to manually configure `DB_URL` for MSSQL

#### Output Example:

```text
# DATABASE CONFIGURATION
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_NAME=test
DB_PASSWORD=123
DB_URL=root:123@tcp(localhost:3307)/test?charset=utf8mb4&parseTime=True&loc=Local

# REDIS CONFIGURATION
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

---

## Migrations

### 1. Initialize & Install Dependencies

```bash
# Generate templates
witches init

# Install dependencies
witches install
```

### 2. Migration Files

**Up Migration** (`./migrate/migrations/1_init.up.sql`):
**Down Migration** (`./migrate/migrations/1_init.down.sql`):

# Install dependencies

```bash
witches migrate up
```

### 3. Migration Commands

| Command | Description |
|---------|-------------|
| `witches migrate up` | Apply 1 pending migration |
| `witches migrate down` | Rollback 1 migration |
| `witches migrate version` | Show current migration version |
| `witches migrate force <version>` | Force set migration version |
| `witches migrate drop` | Drop all tables ⚠️ DANGEROUS |

> **Note:** Migration commands require Docker to be installed and running.

---

## Running the Application

```bash
# Start the application
witches run

# Server running on: http://localhost:8080
```

---

## Project Structure

```
.
├── cmd/
│   └── server/              # Application entry point
│       ├── config/          # Configuration loading
│       └── routers/         # Route definitions
├── internal/
│   ├── dto/                 # Data Transfer Objects
│   │   ├── auth/
│   │   │   ├── request/
│   │   │   └── response/
│   │   └── user/
│   │       ├── request/
│   │       └── response/
│   ├── entity/              # Business entities
│   │   ├── auth/
│   │   └── user/
│   ├── handler/             # HTTP handlers
│   │   ├── auth/
│   │   └── user/
│   ├── mapping/             # DTO ↔ Entity mappers
│   ├── middleware/          # HTTP middleware
│   ├── repository/          # Data access layer
│   │   ├── auth/
│   │   └── user/
│   ├── usecase/             # Business logic
│   │   ├── auth/
│   │   └── user/
│   └── utils/               # Utility functions
├── logs/                    # Application logs
├── migrate/                 # Database migrations
│   └── migrations/
├───pkg
│   └───redis                # Redis client
└── swagger/                 # Swagger documentation

```

---

## Clean Architecture

This project follows **Clean Architecture** by Robert C. Martin (Uncle Bob).

### Layer Mapping

| Directory | Layer | Responsibility |
|-----------|-------|----------------|
| `internal/entity/` | **Entities** | Core business rules, independent of frameworks |
| `internal/usecase/` | **Use Cases** | Application-specific business rules |
| `internal/handler/`<br>`internal/dto/` | **Interface Adapters** | Convert data between external and internal layers |
| `internal/repository/` | **Interface Adapters** | Database operations, data persistence |
| `internal/middleware/` | **Frameworks & Drivers** | Framework-dependent components |
| `cmd/server/` | **Frameworks & Drivers** | Application entry point, dependency injection |
| `pkg/` | **Frameworks & Drivers** | Shared utilities (DB, Redis, logging) |

<div align="center">
  <img src="../image/arc.png" alt="Clean Architecture" width="400"/>
</div>

---

## Technologies

| Component | Technology |
|-----------|------------|
| **HTTP Framework** | [Gin](https://github.com/gin-gonic/gin) |
| **ORM** | [GORM](https://gorm.io/) |
| **Logger** | [Zap](https://github.com/uber-go/zap) |
| **Migration** | [Golang-Migrate](https://github.com/golang-migrate/migrate) |
| **Cache** | [Redis](https://redis.io/) |
| **Swagger** | [EasyJSON](https://github.com/mailru/easyjson) |
| **CLI** | [Cobra](https://github.com/spf13/cobra) |
| **Database** | PostgreSQL / MySQL / MSSQL |
...

All the open-source contributors who make these tools possible ❤️
---

## Best Practices

- Always write both `up` and `down` migrations
- Test migrations in development before production
- Never edit applied migrations - create new ones
- Backup database before running migrations on production
- Keep migrations independent of application code
- Use environment variables for configuration

---

## License

MIT License

---

## Support

- Open an [issue](https://github.com/DVV-15324/witches/issues)

---

<div align="center">
  <br>
  <div><b> Open Source · Non-Profit</b></div>
  This project is completely open source and non-profit.<br>
  Built with passion for the Go community.
  Made by <a href="https://github.com/DVV-15324">DVV-15324</a>
  <br><br>
</div>