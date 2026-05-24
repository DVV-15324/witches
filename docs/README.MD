<div align="center">

<img src="../logo/logo.png" alt="Witches Logo" width="200"/>

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

## Features

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

## Supported Databases

### SQL
- PostgreSQL
- MySQL
- MSSQL

### NoSQL
- Redis

---

## Clone Project
```text
git clone https://github.com/DVV-15324/witches.git
cd witches
```
---

## Configuration

Edit `witches.env` file:

`# Database Type: postgres, mysql, mssql`
```text
DB_TYPE=mysql
```

`# Database Credentials`
```text
DB_USER=root
DB_PASSWORD=Your!StrongPassword123
DATABASE=test_db
```
`# Connection URL`
```text
DB_URL=mysql://root:Your!StrongPassword123@tcp(localhost:1502)/test_db?charset=utf8mb4&parseTime=True&loc=Local
```
### Available Database Types
- `postgres`
- `mysql`
- `mssql`

---

## Quick Start

### 1. Initialize Project
```text
witches init example
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
witches database up
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
### Create file:
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
### Run:

```text
witches migrate up
```

## Migration Commands

| Command | Description |
|---|---|
| witches migrate up | Apply all pending migrations |
| witches migrate up 1 | Apply 1 pending migration |
| witches migrate down | Rollback all migrations |
| witches migrate down 1 | Rollback 1 migration |
| witches migrate version | Show current migration version |
| witches migrate force <version> | Force set migration version |
| witches migrate drop | Drop all database tables |

---

### Best Practices
- Always write both up and down migrations
- Test migrations in development before production
- Never edit applied migrations - create new ones
- Backup database before running in production
- Run migrations separately from application runtime

### Troubleshooting

**Error: "Dirty database version"**
`witches migrate force <last_clean_version>`

**Error: "file does not exist"**
- Check if migration files exist in migrate/migrations/
- Verify Docker Desktop file sharing settings (Windows)
---

### 6. Start Application
```text
witches run
```
### Or with Swagger and easyjson generation:
```text
witches run -all
```
### Or using Go directly:
```text
go run .
```
---
## Project Structure
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
image/roadmap.png

### Refresh Token Flow
image/fresher_token.png

---

## Development Notes

- Database lifecycle is managed through Docker Compose
- Migration should run independently from the runtime application
- witches.env contains all application configuration
- Unit tests are applied across core functions and modules
- For Windows Docker users: ensure drive is shared in Docker Desktop settings

---

## License

MIT

---

## Support

For issues or questions, please open an issue on GitHub.