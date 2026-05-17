<div align="center">

<img src="./logo/logo.png" alt="Witches Logo" width="200"/>

# Witches API

### Fast & Scalable Golang Backend

<p>
  REST API built with <b>Go</b>, designed for performance,
  clean architecture, and modern backend development.
</p>

<p>
  <img src="https://img.shields.io/badge/Go-1.23-00ADD8?style=for-the-badge&logo=go">
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

---

## Quick Start

```bash
git clone https://github.com/DVV-15324/witches.git

cd witches

go run main.go
0. cmd

1. database
    Link: core-v\internal\core\database
    1.1. Nosql:
        - Redis
    1.2. Sql:
        - PostgreSQL
        - MySQL
        - MSSQL
    1.3. Migration:
        - Golang-Migrate

2. handle
    Link: core-v\internal\core\handle
    - Gin:

3. logger
    Link: core-v\internal\core\logger
    - Zap

4. reponse
    Link: core-v\internal\core\reponse

5. sitemap
    Link: core-v\internal\core\reponse

6. wagger
    Link: core-v\internal\core\swagger

7. image
    - RoadMap
        Link: core-v\image\roadmap.png
    - Luồng
        Link: core-v\image\fresher_token.png

8. Unit test
    - Được dùng trên mọi hàm và chức năng
