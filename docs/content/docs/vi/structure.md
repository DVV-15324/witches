---
title: "Kiến trúc"
weight: 4
bookCollapseSection: true
---

## Ánh xạ Clean Architecture

| Thư mục | Layer | Vai trò |
|---------|-------|---------|
| **`cmd/server`** | **Delivery/Interface Adapters** | Entry point của ứng dụng, khởi tạo server, config, routes |
| **`internal/`** | **Business Logic** | Toàn bộ logic nghiệp vụ (domain) |
| **`pkg/`** | **Infrastructure** | Các thành phần dùng chung (Redis, database, external services) |
| **`migrate/`** | **Infrastructure** | Database migrations |
| **`logs/`** | **Infrastructure** | Log files |
| **`swagger/`** | **Delivery/Interface Adapters** | API documentation |

---

### Ánh xạ từng tầng:

| Thư mục | Tầng Clean Architecture | Chức năng |
|---------|------------------------|-----------|
| `model/` | **Entities** | Định nghĩa cấu trúc dữ liệu nghiệp vụ (domain models) |
| `usecase/` | **Use Cases** | Chứa logic nghiệp vụ, orchestrate các operation |
| `handler/` | **Interface Adapters - Controllers** | Xử lý HTTP requests, gọi usecase |
| `dto/request/` | **Interface Adapters - DTO** | Validation input từ client |
| `dto/response/` | **Interface Adapters - DTO** | Format dữ liệu trả về client |
| `mapping/` | **Interface Adapters - Mapper** | Chuyển đổi giữa DTO và Domain và  Model |
| `repository/` | **Interface Adapters - Gateway** | Interface cho database operations |

---

## Luồng dữ liệu Clean Architecture

```
HTTP Request
    ↓
Delivery Layer (cmd/server/routers)
    ↓
Handler (internal/{domain}/handler)
    ↓
DTO → Domain (mapping/)
    ↓
UseCase (usecase/)
    ↓
Repository Interface (repository/)
    ↓
    ↓
Domain → Model (mapping/)
    ↓
(pkg/, database, redis)
```
---
## Kiến trúc tổng quát Clean Architecture của Witches

![arc](/images/arc.png)