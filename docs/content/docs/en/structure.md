---
title: "Structure"
weight: 4
bookCollapseSection: true
---

## Clean Architecture Mapping

| Directory | Layer | Role |
|-----------|-------|------|
| **`cmd/server`** | **Delivery / Interface Adapters** | Application entry point, server initialization, config, routes |
| **`internal/`** | **Business Logic** | All business logic (domain) |
| **`pkg/`** | **Infrastructure** | Shared components (Redis, database, external services) |
| **`migrate/`** | **Infrastructure** | Database migrations |
| **`logs/`** | **Infrastructure** | Log files |
| **`swagger/`** | **Delivery / Interface Adapters** | API documentation |

---

### Layer Mapping per Domain:

| Directory | Clean Architecture Layer | Responsibility |
|-----------|--------------------------|----------------|
| `model/` | **Entities** | Defines business data structures (domain models) |
| `usecase/` | **Use Cases** | Contains business logic, orchestrates operations |
| `handler/` | **Interface Adapters – Controllers** | Handles HTTP requests, calls use cases |
| `dto/request/` | **Interface Adapters – DTO** | Client input validation |
| `dto/response/` | **Interface Adapters – DTO** | Formats data returned to the client |
| `mapping/` | **Interface Adapters – Mapper** | Converts between DTO, Domain, and Model |
| `repository/` | **Interface Adapters – Gateway** | Interface for database operations |

---

## Clean Architecture Data Flow

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
Domain → Model (mapping/)
    ↓
Infrastructure (pkg/, database, redis)
```
---
## Witches Clean Architecture Overview

![arc](/images/arc.png)