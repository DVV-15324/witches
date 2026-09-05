[Vietnamese](./README.vi.md) | [English](./README.md)

<div align="center">

<img src="../logo/logo.png" alt="witches Logo" width="250"/>

### Go Backend CLI Tool
<p>
    REST API scaffolding for <b>Go</b>, built for <b>simplicity</b>, <b>low-resource usage</b>, 
    and <b>practical development</b>.
  </p>
  
[![Security Scan](https://github.com/DVV-15324/witches/actions/workflows/security.yml/badge.svg)](https://github.com/DVV-15324/witches/actions/workflows/security.yml)
[![codecov](https://codecov.io/gh/DVV-15324/witches/branch/main/graph/badge.svg)](https://app.codecov.io/github/DVV-15324/witches)

<p>
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go">
  <img src="https://img.shields.io/badge/Gin-Web_Framework-008ECF?style=for-the-badge">
  <img src="https://img.shields.io/badge/GORM-ORM-25A162?style=for-the-badge&logo=gorm">
  <img src="https://img.shields.io/badge/Redis-Cache-DC382D?style=for-the-badge&logo=redis">
  <img src="https://img.shields.io/badge/Cobra-CLI-4B32C3?style=for-the-badge">
  <img src="https://img.shields.io/badge/Swagger-API_Docs-85EA2D?style=for-the-badge&logo=swagger">
</p>

</div>

---

## ⚠️ Status

> **Not ready for production use.**  
> This is a personal project. Use at your own risk.

---

## Features

| Feature | Description |
|---------|-------------|
| **Swagger** | Auto-generated API documentation |
| **Migration** | Version control for database schema |
| **Redis Support** | Optional caching |
| **Rate Limiting** | Protect your APIs from abuse |
| **GORM Support** | ORM for database operations |
| **EasyJSON** | Automate JSON serialization code generation |
| **Metrics** | metrics, timing, and tracing |
| **Logging** | Structured logging with Zap |
| **CLI** | Cobra-based command-line interface |
| **Slow Query** | Database query tracing and slow query logging |

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

Edit `witches.env` file with your database credentials.

### 4. Generate database URL from env

```bash
witches database generate
```

### 5. Run database migrations

```bash
witches migrate up
```

### 6. Init & Install Dependencies

```bash
witches init      # Generate template files
witches install   # Install dependencies
```

### 7. Run the Server

```bash
witches run
```
---

## Commands

| Command | Description |
|---------|-------------|
| `witches create <name>` | Create a new project |
| `witches init captain` | Generate template files and structure |
| `witches init member` | Generate template files and structure |
| `witches install` | Install dependencies |
| `witches run` | Start the server |
| `witches migrate up` | Apply all migrations |
| `witches migrate down` | Rollback last migration |
| `witches database generate` | Generate database URL |
|`witches add <your_module>` | Add a new module |
|`witches rm  <your_module>`| Remove a module |
|`witches link <module> <repo-url>` | Link module from GitHub repo |

---

## License

MIT License - see the [LICENSE](../LICENSE) file for details.

---

<div align="center">
  <br>
  Made with ❤️ by <a href="https://github.com/DVV-15324">DVV-15324</a>
  <br><br>
</div>
