[🇻🇳 Tiếng Việt](./docs/README.vi.md) | [🇬🇧 English](./docs/README.md)
<div align="center">

<img src="../logo/logo.png" alt="Witches Logo" width="200"/>

### Backend Golang Nhanh & Mở Rộng

<p>
  REST API được xây dựng với <b>Go</b>, thiết kế để tối ưu hiệu năng,
  Kiến trúc Clean Architecture và phát triển backend hiện đại.
</p>

<p>
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go">
  <img src="https://img.shields.io/badge/Gin-Web_Framework-008ECF?style=for-the-badge">
  <img src="https://img.shields.io/badge/SQL-Database-orange?style=for-the-badge&logo=sql">
  <img src="https://img.shields.io/badge/Swagger-API_Docs-green?style=for-the-badge">
</p>

</div>

---

## Tính năng
- Xác thực JWT
- Tạo tài liệu Swagger tự động
- Quản lý Migration Database
- Tiện ích băm (Hash Utilities)
- Che giấu UID (UID Masking)
- Kiến trúc Clean Architecture

---

# Bắt đầu nhanh

## 1. Khởi tạo dự án

### Cài đặt Witches
```bash
go install github.com/DVV-15324/witches@latest
```

### Tạo dự án mới
```bash
witches create example --db=mysql
```

#### Kết quả: `example/witches.env`
```env
APP_PORT=8080
DB_PASSWORD=your_password
DB_NAME=your_database
DB_HOST=localhost
DB_PORT=3306
REDIS_PORT=6379
DB_DRIVER=mysql
```

```bash
cd example
```

## 2. Cài đặt dependencies và khởi tạo templates

### Khởi tạo templates
```bash
witches init
```

### Cài đặt dependencies
```bash
witches install
```

## 3. Cấu hình Database
Sửa file `witches.env`:
```env
APP_PORT=3000
DB_PASSWORD=123
DB_NAME=test
DB_HOST=localhost
DB_PORT=3307
REDIS_PORT=1504
DB_DRIVER=mysql
```

## 4. Khởi động Database
```bash
witches database docker-up
```

#### Kết quả:
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

## 5. Chạy Migration

**Tạo file migration:**

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

**Chạy:**
```bash
witches migrate docker-up
```

**Rollback:**
```bash
witches migrate docker-down
```

---

## Lệnh Migration

| Lệnh | Mô tả |
|------|-------|
| `witches migrate up` | Áp dụng tất cả migration đang chờ |
| `witches migrate up 1` | Áp dụng 1 migration đang chờ |
| `witches migrate down` | Rollback tất cả migration |
| `witches migrate down 1` | Rollback 1 migration |
| `witches migrate version` | Hiển thị phiên bản migration hiện tại |
| `witches migrate force <version>` | Gán phiên bản migration |
| `witches migrate drop` | Xóa tất cả bảng trong database |

## 6. Chạy ứng dụng

### Chạy thông thường
```bash
witches run
```

---

## Clean Architecture của sự án

<div align="center">
  <img src="../image/arc.png" alt="Clean Architecture" width="400"/>
</div>

---

## Cấu trúc dự án

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

## Kiến trúc Clean Architecture

Dự án được xây dựng dựa trên **Clean Architecture** của Robert C. Martin (Uncle Bob).

### Ánh xạ thư mục với Clean Architecture

| Thư mục | Layer | Vai trò |
|---------|-------|----------|
| `internal/entity/` | **Entities** | Chứa các entity và quy tắc nghiệp vụ cốt lõi, không phụ thuộc vào framework hay database. |
| `internal/usecase/` | **Use Cases** | Chứa các quy tắc nghiệp vụ cụ thể, điều phối luồng xử lý và giao tiếp với repository. |
| `internal/handler/`<br>`internal/dto/` | **Interface Adapters** | Chuyển đổi dữ liệu giữa bên ngoài và ứng dụng. Handler nhận request, gọi use case và trả response. |
| `internal/repository/` | **Interface Adapters** | Triển khai repository, chuyển đổi dữ liệu giữa database và entity. |
| `internal/middleware/` | **Frameworks & Drivers** | Chứa middleware phụ thuộc framework: JWT, CORS, Logging, Recovery,... |
| `cmd/server/` | **Frameworks & Drivers** | Điểm vào ứng dụng, khởi tạo dependencies, router và start server. |
| `pkg/` | **Frameworks & Drivers** | Các thành phần chia sẻ: database connection, Redis, logging, utilities. |

---

## Hỗ trợ Database

### SQL
- PostgreSQL
- MySQL
- MSSQL

### NoSQL
- Redis

---

## Best Practices

- Luôn viết cả migration `up` và `down`
- Test migration trên môi trường dev trước khi chạy production
- Không sửa migration đã áp dụng - hãy tạo migration mới
- Backup database trước khi chạy migration trên production
- Chạy migration độc lập với ứng dụng

---

## Công nghệ sử dụng

| Thành phần | Công nghệ |
|-----------|-----------|
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

## Hỗ trợ

Vui lòng mở issue trên GitHub nếu có câu hỏi hoặc vấn đề cần hỗ trợ.
```