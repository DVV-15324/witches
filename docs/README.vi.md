[Vietnamese](./README.vi.md) | [English](./README.md)

<div align="center">

<img src="../logo/logo.png" alt="witches Logo" width="250"/>

### Công cụ CLI Backend Go

<p>
  Scaffolding REST API cho <b>Go</b>, thiết kế đơn giản,
  dùng ít tài nguyên, thực dụng cho phát triển thực tế.
</p>

<p>
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go">
  <img src="https://img.shields.io/badge/Gin-Web_Framework-008ECF?style=for-the-badge">
  <img src="https://img.shields.io/badge/GORM-ORM-25A162?style=for-the-badge&logo=gorm">
</p>

</div>

---

## ⚠️ Trạng thái

**Chưa sẵn sàng để dùng.**  
Đây là dự án cá nhân. Dùng là tự chịu rủi ro.  
Có thể 30 năm nữa mới xong. Hoặc không bao giờ.

---

## Tính năng

| Tính năng | Mô tả |
|-----------|-------|
| **Xác thực JWT** | Auth bằng token + Refresh Token |
| **Swagger** | Tự động tạo tài liệu API |
| **Migration** | Quản lý schema database |
| **Mã hóa** | Băm mật khẩu an toàn |
| **Redis** | Hỗ trợ cache (không bắt buộc) |
| **Giới hạn tốc độ** | Bảo vệ API |
| **GORM** | ORM cho database |
|**EASYJSON** | Tự động hóa Sinh mã EasyJSON |
|****	| Metrics, Timing, Database Tracing, and Structured Logging|

---

## Bắt đầu nhanh

### 1. Cài đặt

```bash
go install github.com/DVV-15324/witches@latest
```

### 2. Tạo project mới

```bash
witches create example
cd example
```

### 3. Cấu hình

Sửa file `witches.env` theo nhu cầu.

### 4. Chạy

```bash
witches run
```

---

## Cấu hình

Xem file `witches.env` để biết tất cả tùy chọn.

---

## Database

Hỗ trợ: MySQL, PostgreSQL, MSSQL

```bash
witches database generate
```

---

## Migration

```bash
witches migrate up
witches migrate down
```

---

## Cấu trúc project

```
├── cmd/           # Điểm vào
├── internal/       # Code nội bộ
│   ├── shared/     # Dùng chung
│   └── *-service/  # Module nghiệp vụ
├── migrate/        # SQL migration
└── pkg/            # Package tái sử dụng
```

---

## Công nghệ sử dụng

- **HTTP:** Gin
- **ORM:** GORM
- **Migration:** Golang-Migrate
- **Cache:** Redis
- **CLI:** Cobra

---

## Giấy phép

MIT License

---

<div align="center">
  <br>
  Làm bởi <a href="https://github.com/DVV-15324">DVV-15324</a>
  <br><br>
</div>
