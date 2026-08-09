[Tiếng Việt](./README.vi.md) | [English](./README.md)

<div align="center">

<img src="../logo/logo.png" alt="witches Logo" width="250"/>

### Go Backend CLI Tool

<p>
  Công cụ tạo scaffold REST API cho <b>Go</b>, được thiết kế đơn giản,
  tiết kiệm tài nguyên.
</p>

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

## ⚠️ Trạng thái

> **Chưa sẵn sàng cho môi trường production.**  
> Đây là dự án cá nhân. Sử dụng với rủi ro của riêng bạn.

---

## Tính năng

| Tính năng | Mô tả |
|-----------|-------|
| **Xác thực JWT** | Xác thực dựa trên token với cơ chế Refresh Token |
| **Swagger** | Tự động tạo tài liệu API |
| **Migration** | Quản lý phiên bản cơ sở dữ liệu |
| **Hash Utilities** | Mã hóa mật khẩu an toàn |
| **Hỗ trợ Redis** | Tùy chọn caching và lưu trữ session |
| **Giới hạn tốc độ** | Bảo vệ API khỏi bị lạm dụng |
| **Hỗ trợ GORM** | ORM cho các thao tác cơ sở dữ liệu |
| **EasyJSON** | Tự động tạo mã JSON serialization |
| **Metrics** | Đo lường, timing và theo dõi |
| **Logging** | Logging có cấu trúc với Zap |
| **CLI** | Giao diện dòng lệnh dựa trên Cobra |
| **Slow Query** | Theo dõi truy vấn database và log truy vấn chậm |

---

## Bắt đầu nhanh

### 1. Cài đặt

```bash
go install github.com/DVV-15324/witches@latest
```

### 2. Tạo dự án mới

```bash
witches create example
cd example
```

### 3. Cấu hình

Chỉnh sửa file `witches.env`

### 4. Thiết lập Database

```bash
witches database generate
```

### 5. Khởi tạo & Cài đặt Dependencies

```bash
witches init      # Tạo file template
witches install   # Cài đặt dependencies
```

### 6. Chạy Server

```bash
witches run
```

---

## Câu lệnh

| Câu lệnh | Mô tả |
|----------|-------|
| `witches create <name>` | Tạo dự án mới |
| `witches init` | Tạo file template và cấu trúc |
| `witches install` | Cài đặt dependencies |
| `witches run` | Khởi động server |
| `witches migrate up` | Chạy tất cả migrations |
| `witches migrate down` | Rollback migration cuối cùng |
| `witches database generate` | Tạo URL kết nối database |
| `witches add service=your_service` | Thêm 1 service mới  |

---

## Giấy phép

MIT License - xem file [LICENSE](LICENSE) để biết thêm chi tiết.

---

<div align="center">
  <br>
  Được tạo với ❤️ bởi <a href="https://github.com/DVV-15324">DVV-15324</a>
  <br><br>
</div>
