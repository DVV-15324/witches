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
- Swagger Documentation
- Database Migration
- Docker Support
- Hash Utilities
- Sitemap Support
- UID Masking
- Clean Architecture (Coming Soon)
- Rate Limiting (Coming Soon)
- Redis Cache (Coming Soon)
- Multi-Database Support (Coming Soon)

---

### Clean Architecture of this project

<div align="center">

<img src="../image/arc.png"/>

</div>

---

### UML (Unified Modeling Language)

<div align="center">
<img src="../image/uml.png"/>
</div>

---

### Clean Architecture Structure of this project

```text
├───cmd
│   └───server
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
├───migrate
│   └───migrations
└───pkg
    ├───logs
    ├───postgres
    └───redis
```

---

### Clean Architecture by Robert C. Martin (Uncle Bob)
https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

---

### Mapping Folder Structure to Clean Architecture by Robert C. Martin (Uncle Bob)

| Project Directory                         | Clean Architecture Layer    | Role                                                                                                                                                                                           |
| ----------------------------------------- | --------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `internal/entity/`                        | **Entities**                | Contains core entities and business rules. This is the most independent layer, with no dependency on frameworks, databases, or any external technology.                                        |
| `internal/usecase/`                       | **Use Cases**               | Contains application-specific business rules (use cases). Orchestrates the business flow, works with entities, and uses necessary interfaces to access data or external services.               |
| `internal/handler/`<br>`internal/dto/`    | **Interface Adapters**      | Converts data between the outside world and the application. Handlers receive requests, call use cases, and return responses. DTOs are used for data exchange between layers.                 |
| `internal/repository/`                    | **Interface Adapters**      | Contains repository implementations, responsible for converting data between database and entities, while implementing interfaces required by use cases.                                       |
| `internal/middleware/`                    | **Frameworks & Drivers**    | Contains framework-dependent middleware such as JWT Authentication, CORS, Logging, Recovery, Rate Limiting, etc.                                                                               |
| `cmd/server/`                             | **Frameworks & Drivers**    | Application entry point. Initializes dependencies, configures routers, middleware, database connections, and starts the HTTP server.                                                            |
| `pkg/postgres/`<br>`pkg/redis/`<br>`pkg/logs/` | **Frameworks & Drivers**    | Contains shared infrastructure components such as database connections, Redis client, logging, and other technical utilities.                                                                  |
| `migrate/migrations/`                     | **Infrastructure**          | Contains migration files for database schema management. This is an infrastructure component and does not belong to the application's business layer.                                          |

---

### General Flow of Processing

```text
Client
   │
   ▼
Handler
   │
   ▼
UseCase
   │
   ▼
Repository
   │
   ▼
PostgreSQL / Redis
```

And the most important rule of Clean Architecture:

```text
Frameworks & Drivers
        ↓
Interface Adapters
        ↓
Use Cases
        ↓
Entities
```

Inner layers must not depend on outer layers. This makes the application easier to maintain, test, and change technologies in the future.

---

### Supported Databases

### SQL
- PostgreSQL
- MySQL
- MSSQL

### NoSQL
- Redis


## Quick Start

### 1. Initialize Project
---

## Install Project
```text
go install https://github.com/DVV-15324/witches.git
```

---

```text
witches init example --db=mysql
cd example
```

### 2. Install Dependencies
```text
witches install
```

### 3. Configure Database
Edit `witches.env`:
```text
DB_PASSWORD=Your!StrongPassword123
DATABASE=test_db
DB_PROFILE=mysql
```

### 4. Start Database
```text
witches database docker-up
```

#### Output:
```text
DB_PASSWORD=Your!StrongPassword123
DATABASE=test_db
DB_PROFILE=mysql
DB_USER=root
DB_URL=mysql://root:Your!StrongPassword123@tcp(localhost:1502)/test_db?charset=utf8mb4&parseTime=True&loc=Local
```

### 5. Run Migrations

**Create file:**

**Up migration** (`./migrate/migrations/1_init.up.sql`):
```text
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Down migration** (`./migrate/migrations/1_init.down.sql`):
```text
DROP TABLE IF EXISTS users;
```

**Run:**
```text
witches migrate docker-up
```

**Down:**
```text
witches migrate docker-down
```

---

## Migration Commands

| Command | Description |
|---------|-------------|
| `witches migrate up` | Apply all pending migrations |
| `witches migrate up 1` | Apply 1 pending migration |
| `witches migrate down` | Rollback all migrations |
| `witches migrate down 1` | Rollback 1 migration |
| `witches migrate version` | Show current migration version |
| `witches migrate force <version>` | Force set migration version |
| `witches migrate drop` | Drop all database tables |

---

### Best Practices
- Always write both up and down migrations
- Test migrations in development before production
- Never edit applied migrations - create new ones
- Backup database before running in production
- Run migrations separately from application runtime

### Troubleshooting

**Error: "Dirty database version"**
```text
witches migrate force <last_clean_version>
```

**Error: "file does not exist"**
- Check if migration files exist in `migrate/migrations/`
- Verify Docker Desktop file sharing settings (Windows)

---

### 6. Start Application
```text
witches run
```

**Or with Swagger and easyjson generation:**
```text
witches run -all
```

**Or using Go directly:**
```text
go run .
```

---

## Project Structure
```text
├───cmd
│   ├───cmd_database
│   ├───cmd_migrate
│   ├───cmd_run
│   └───cmd_utils
├───docs
├───example
│   └───migrate
│       └───migrations
├───image
├───logo
├───pkg
│   └───core
│       ├───arc-golang
│       │   ├───cmd
│       │   │   └───server
│       │   │       └───routers
│       │   ├───internal
│       │   │   ├───dto
│       │   │   │   ├───auth
│       │   │   │   │   ├───request
│       │   │   │   │   └───response
│       │   │   │   └───user
│       │   │   │       ├───request
│       │   │   │       └───response
│       │   │   ├───entity
│       │   │   │   ├───auth
│       │   │   │   └───user
│       │   │   ├───handler
│       │   │   │   ├───auth
│       │   │   │   └───user
│       │   │   ├───mapping
│       │   │   ├───middleware
│       │   │   ├───repository
│       │   │   │   ├───auth
│       │   │   │   └───user
│       │   │   ├───usecase
│       │   │   │   ├───auth
│       │   │   │   └───user
│       │   │   └───utils
│       │   ├───logs
│       │   ├───migrate
│       │   │   └───migrations
│       │   └───pkg
│       │       ├───postgres
│       │       └───redis
│       ├───database
│       │   ├───connect
│       │   │   ├───nosql
│       │   │   ├───sql
│       │   │   └───storage
│       │   ├───migrate
│       │   ├───mssql-db
│       │   ├───mysql-db
│       │   ├───postgres-db
│       │   └───redis-db
│       ├───gen_easyjson
│       │   └───generator
│       ├───handle_swagger
│       ├───jwt
│       ├───response_logger
│       │   └───logger
│       └───utils
└───test
    └───swagger
```

---

## Technologies

| Component | Technology |
|-----------|------------|
| HTTP Framework | Gin |
| Logger | Zap |
| Migration | Golang-Migrate |
| Cache | Redis |
| Documentation | Swagger |
| Database | PostgreSQL / MySQL / MSSQL |

---

## Images

### Roadmap
`image/roadmap.png`

### Refresh Token Flow
`image/fresher_token.png`

---

## Development Notes

- Database lifecycle is managed through Docker Compose
- Migration should run independently from the runtime application
- `witches.env` contains all application configuration
- Unit tests are applied across core functions and modules
- For Windows Docker users: ensure drive is shared in Docker Desktop settings

---

## License

MIT

---

## Support

For issues or questions, please open an issue on GitHub.