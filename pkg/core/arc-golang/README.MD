https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

---
### Clean Architecture

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
###  Bản đồ ánh xạ Folder Structure với Clean Architecture

| Thư mục trong dự án                            | Layer trong Clean Architecture | Vai trò                                                                                                                                                                      |
| ---------------------------------------------- | ------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `internal/entity/`                             | **Entities**                   | Chứa các entity và business rules cốt lõi của hệ thống. Đây là tầng độc lập nhất, không phụ thuộc vào framework, database hay bất kỳ công nghệ bên ngoài nào.                |
| `internal/usecase/`                            | **Use Cases**                  | Chứa các application business rules (use cases). Điều phối luồng xử lý nghiệp vụ, làm việc với entity và các interface cần thiết để truy cập dữ liệu hoặc dịch vụ bên ngoài. |
| `internal/handler/`<br>`internal/dto/`         | **Interface Adapters**         | Chuyển đổi dữ liệu giữa thế giới bên ngoài và ứng dụng. Handler tiếp nhận request, gọi use case và trả response. DTO được sử dụng để trao đổi dữ liệu giữa các tầng.         |
| `internal/repository/`                         | **Interface Adapters**         | Chứa các repository implementation, chịu trách nhiệm chuyển đổi dữ liệu giữa database và entity, đồng thời hiện thực các interface mà use case yêu cầu.                      |
| `internal/middleware/`                         | **Frameworks & Drivers**       | Chứa các middleware phụ thuộc vào framework như JWT Authentication, CORS, Logging, Recovery, Rate Limiting,...                                                               |
| `cmd/server/`                                  | **Frameworks & Drivers**       | Điểm khởi động ứng dụng. Khởi tạo các dependency, cấu hình router, middleware, database connection và chạy HTTP server.                                                      |
| `pkg/postgres/`<br>`pkg/redis/`<br>`pkg/logs/` | **Frameworks & Drivers**       | Chứa các thành phần hạ tầng dùng chung như database connection, Redis client, logging và các tiện ích kỹ thuật khác.                                                         |
| `migrate/migrations/`                          | **Infrastructure**             | Chứa các file migration phục vụ quản lý schema database. Đây là thành phần hạ tầng và không thuộc business layer của ứng dụng.                                               |

---

### Luồng xử lý tổng quát

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

Và quy tắc quan trọng nhất của Clean Architecture:

```text
Frameworks & Drivers
        ↓
Interface Adapters
        ↓
Use Cases
        ↓
Entities
```

Các tầng bên trong không được phụ thuộc vào các tầng bên ngoài. Điều này giúp ứng dụng dễ bảo trì, dễ kiểm thử và dễ thay đổi công nghệ trong tương lai.

