[Vietnamese](./README.vi.md) | [English](./README.md)

<div align="center">

<img src="../logo/logo.png" alt="witches Logo" width="250"/>

### Go Backend CLI Tool

<p>
  REST API scaffolding for <b>Go</b>, designed for simplicity,
  low-resource usage, and practical development.
</p>

<p>
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go">
  <img src="https://img.shields.io/badge/Gin-Web_Framework-008ECF?style=for-the-badge">
  <img src="https://img.shields.io/badge/GORM-ORM-25A162?style=for-the-badge&logo=gorm">
</p>

</div>

---

## ⚠️ Status

**Not ready for use.**  
This is a personal project. Use at your own risk.  
May be ready in 30 years. Or never.

---

## Features

| Feature | Description |
|---------|-------------|
| **JWT Auth** | Token-based authentication with Refresh Token |
| **Swagger** | Auto-generated API docs |
| **Migration** | Version control for database schema |
| **Hash Utilities** | Password hashing |
| **Redis Support** | Optional caching |
| **Rate Limiting** | Protect your APIs |
| **GORM Support** | ORM for database operations |

---

## Quick Start

### 1. Installation

```bash
go install github.com/DVV-15324/witches@latest
```

### 2. Create New Project

```bash
witches create example
cd example
```

### 3. Configure

Edit `witches.env` with your settings.

### 4. Run

```bash
witches run
```

---

## Configuration

See `witches.env` for all available options.

---

## Database

Supported: MySQL, PostgreSQL, MSSQL

```bash
witches database generate
```

---

## Migrations

```bash
witches migrate up
witches migrate down
```

---

## Project Structure

```
├── cmd/           # Entrypoint
├── internal/       # Private code
│   ├── shared/     # Shared utilities
│   └── *-service/  # Domain modules
├── migrate/        # SQL migrations
└── pkg/            # Reusable packages
```

---

## Technologies

- **HTTP:** Gin
- **ORM:** GORM
- **Migration:** Golang-Migrate
- **Cache:** Redis
- **CLI:** Cobra

---

## License

MIT License

---

<div align="center">
  <br>
  Made by <a href="https://github.com/DVV-15324">DVV-15324</a>
  <br><br>
</div>

