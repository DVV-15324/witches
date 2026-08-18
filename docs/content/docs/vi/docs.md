---
title: "Bắt đầu với Witches"
weight: 1
bookCollapseSection: true
---

## Hướng dẫn nhanh

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

Sửa file `witches.env` với thông tin kết nối cơ sở dữ liệu của bạn.

### 4. Tạo URL kết nối database từ file env

```bash
witches database generate
```

### 5. Chạy migration cơ sở dữ liệu

```bash
witches migrate up
```

### 6. Khởi tạo và cài đặt các thư viện phụ thuộc

```bash
witches init      # Tạo các file template
witches install   # Cài đặt các thư viện cần thiết
```

### 7. Chạy máy chủ

```bash
witches run
```
---

<div align="center">
  <br>
  Tạo với ❤️ bởi <a href="https://github.com/DVV-15324">DVV-15324</a>
  <br><br>
</div>

