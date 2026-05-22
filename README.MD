<div align="center">

<img src="./logo/logo.png" alt="Witches Logo" width="200"/>

# Witches API

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
</p>

</div>

---

# Features

- JWT Authentication
- Swagger Documentation
- Database Migration
- Clean Architecture
- Docker Support
- Rate Limiting
- Redis Cache
- Multi Database Support
- Unit Testing

---

# Supported Databases

## SQL
- PostgreSQL
- MySQL
- MSSQL

## NoSQL
- Redis

---

# Installation

## Windows

Install `make` using Scoop:

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
scoop install make
```

## Clone Project

```bash
git clone https://github.com/DVV-15324/witches.git
cd witches
```

---

# Configuration

Edit:

```env
witches.env
```

## Database Type

Supported values:

```env
DB_TYPE=postgres
```

Available:
- postgres
- mysql
- mssql

---

# Start Database

```bash
witches database up
```

The command will return a database connection URL.

Example:

```env
DB_URL=postgres://localhost:5432/test_db
```

Update `DB_URL` inside:

```env
witches.env
```

---

# Database Migration

## Migrate Up

```bash
witches migrate up
```

## Migrate Down

```bash
witches migrate down
```

Migration is used to manage SQL database schemas.

Recommended:
- Run migrations in a separate terminal
- Keep migrations isolated from runtime application processes

---

# Run Application

```bash
witches run
```

or:

```bash
go run .
```

---

# Project Structure

```text
internal/
└── core/
    ├── cmd/
    ├── database/
    │   ├── redis/
    │   ├── postgres/
    │   ├── mysql/
    │   ├── mssql/
    │   └── migration/
    ├── handle/
    ├── logger/
    ├── response/
    ├── sitemap/
    └── swagger/
```

---

# Technologies

| Component | Technology |
|---|---|
| HTTP Framework | Gin |
| Logger | Zap |
| Migration | Golang-Migrate |
| Cache | Redis |
| Documentation | Swagger |
| Database | PostgreSQL / MySQL / MSSQL |

---

# Images

## Roadmap

```text
image/roadmap.png
```

## Refresh Token Flow

```text
image/fresher_token.png
```

---

# Development Notes

- Database lifecycle is managed through Docker Compose
- Migration should run independently from the runtime application
- `witches.env` contains all application configuration
- Unit tests are applied across core functions and modules

---

# License

MIT