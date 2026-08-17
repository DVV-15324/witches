---
title: "Tài liệu Witches"
weight: 1
bookCollapseSection: true
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
| `witches add <service>` | Thêm 1 service mới  |

---

## Giấy phép

MIT License - xem file [LICENSE](LICENSE) để biết thêm chi tiết.

---

<div align="center">
  <br>
  Được tạo với ❤️ bởi <a href="https://github.com/DVV-15324">DVV-15324</a>
  <br><br>
</div>
