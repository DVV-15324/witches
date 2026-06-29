[Tiếng Việt](./README.vi.md) | [English](./README.md)

<div align="center">

<img src="../logo/logo.png" alt="Witches Logo" width="200"/>

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

## Features
- JWT Authentication
- Swagger Auto-Generate
- Database Migration
- Hash Utilities
- UID Masking
- Clean Architecture

---

# Quick Start

## 1. Initialize Project

### Install Witches
```bash
go install github.com/DVV-15324/witches@latest
```

### Witches version
```bash
witches version
```

### Create new project
```bash
witches create example --db=mysql
```

#### Output: `example/witches.env`
```env
APP_PORT=YOUR_PORT_APP
DB_PASSWORD=YOUR_PASSWORD
DB_NAME=YOUR_DATABASE
DB_HOST=DB_HOST
DB_PORT=YOUR_PORT_APP
REDIS_PORT=REDIS_PORT
DB_DRIVER=mysql
```

```bash
cd example
```
## 2. Configure Database
Edit `witches.env`:
```env
APP_PORT=3000
DB_PASSWORD=123
DB_NAME=test
DB_HOST=localhost
DB_PORT=3307
REDIS_PORT=1504
DB_DRIVER=mysql
```

## 3. Start Database
```bash
witches database docker-up
```

#### Output:
```text
APP_PORT=3000
DB_PASSWORD=123
DB_NAME=test
DB_HOST=localhost
DB_PORT=3307
REDIS_PORT=1504
DB_DRIVER=mysql
DB_URL=mysql://root:123@tcp(localhost:3307)/test?charset=utf8mb4&parseTime=True&loc=Local
```

## 4. Install Dependencies & Init Templates

### Init Templates
```bash
witches init
```

### Install Dependencies
```bash
witches install
```

## 5. Run Migrations

**Create migration files:**

**Up migration** (`./migrate/migrations/1_init.up.sql`):
```sql
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Down migration** (`./migrate/migrations/1_init.down.sql`):
```sql
DROP TABLE IF EXISTS users;
```

**Run:**
```bash
witches migrate up
```

**Rollback:**
```bash
witches migrate down
```

---

## Migration Commands

| Command | Description |
|---------|-------------|
| `witches migrate up` | Apply 1 pending migrations |
| `witches migrate down` | Rollback 1 migrations |
| `witches migrate version` | Show current migration version |
| `witches migrate force <version>` | Force set migration version |
| `witches migrate drop` | Drop all database tables |

## 6. Run Application

### Run normally
```bash
witches run
```

---

## Clean Architecture of this project

<div align="center">
  <img src="../image/arc.png" alt="Clean Architecture" width="400"/>
</div>

---

## Project Structure

```text
├───cmd
│   └───server
│       ├───config
│       └───routers
├───internal
│   ├───dto
│   │   ├───auth
│   │   │   ├───request
│   │   │   └───response
│   │   └───user
│   │       ├───request
│   │       └───response
│   ├───entity
│   │   ├───auth
│   │   └───user
│   ├───handler
│   │   ├───auth
│   │   └───user
│   ├───mapping
│   ├───middleware
│   ├───repository
│   │   ├───auth
│   │   └───user
│   ├───usecase
│   │   ├───auth
│   │   └───user
│   └───utils
├───logs
├───migrate
│   └───migrations
└───swagger
```

---

## Clean Architecture

This project follows **Clean Architecture** by Robert C. Martin (Uncle Bob).

### Folder Mapping to Clean Architecture

| Directory | Layer | Role |
|-----------|-------|------|
| `internal/entity/` | **Entities** | Core entities and business rules. No dependency on frameworks or databases. |
| `internal/usecase/` | **Use Cases** | Application-specific business rules. Orchestrates flow and communicates with repositories. |
| `internal/handler/`<br>`internal/dto/` | **Interface Adapters** | Converts data between external world and application. Handlers receive requests, call use cases, return responses. |
| `internal/repository/` | **Interface Adapters** | Repository implementations for database-entity conversion. |
| `internal/middleware/` | **Frameworks & Drivers** | Framework-dependent middleware: JWT, CORS, Logging, Recovery, etc. |
| `cmd/server/` | **Frameworks & Drivers** | Entry point, initializes dependencies, routers, server. |
| `pkg/` | **Frameworks & Drivers** | Shared components: database connection, Redis, logging, utilities. |

---

## Supported Databases

### SQL
- PostgreSQL
- MySQL
- MSSQL

### NoSQL
- Redis

---

## Best Practices

- Always write both `up` and `down` migrations
- Test migrations in development before production
- Never edit applied migrations - create new ones
- Backup database before running in production
- Run migrations independently from application runtime

---

## Technologies

| Component | Technology |
|-----------|------------|
| HTTP Framework | Gin |
| Logger | Zap |
| Migration | Golang-Migrate |
| Cache | Redis |
| Swagger | EasyJSON, Cobra |
| Database | PostgreSQL / MySQL / MSSQL |

---

## License

MIT

---

## Support

Please open an issue on GitHub for questions or support.
```
