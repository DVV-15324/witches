---
title: "Cấu hình witches.env"
weight: 3
bookCollapseSection: true
---

Dưới đây là bảng mô tả chi tiết từng biến môi trường trong file `witches.env`:

## Bảng mô tả biến môi trường

| Biến | Mô tả | Loại | Giá trị mặc định |
|------|-------|------|------------------|
| **SERVER CONFIGURATION** |
| `APP_PORT` | Cổng chạy ứng dụng HTTP | `int` | `8080` |
| `APP_HOST` | Hostname binding | `string` | `localhost` |
| `METRIC_PORT` | Cổng cho Prometheus metrics | `int` | `8088` |
| `UID_BITS` | Số bit sử dụng cho ID (dùng trong Snowflake/ULID) | `int` | `26` |
| `REQUEST_KEY` | Key để lưu request context | `string` | `request_context` |
| **LOGGING** |
| `LOG_PATH` | Thư mục lưu log file | `string` | `./logs` |
| **CORS** |
| `CORS_ALLOW_ORIGINS` | Các origin được phép truy cập (* = tất cả) | `string` | `*` |
| `CORS_ALLOW_METHODS` | HTTP methods được phép | `string` | `GET,POST,PUT,DELETE,OPTIONS` |
| `CORS_ALLOW_HEADERS` | Headers được phép gửi lên | `string` | `Content-Type,Authorization` |
| **TIMEZONE** |
| `TIMEZONE` | Múi giờ của ứng dụng | `string` | `UTC` |
| `SUPPORTED_LANGUAGES` | Các ngôn ngữ hỗ trợ (dùng cho i18n) | `string` | `en-US,vi-VN` |
| **SECURITY** |
| `JWT_SECRET` | Khóa bí mật để ký JWT | `string` | `your_secret_key` |
| **DATABASE CONFIGURATION** |
| `DB_DRIVER` | Driver database (postgres/mysql/mssql) | `string` | `your_driver` |
| `DB_USER` | Username kết nối database | `string` | `your_user` |
| `DB_HOST` | Hostname database | `string` | `localhost` |
| `DB_PORT` | Cổng database | `int` | `3306` |
| `DB_NAME` | Tên database | `string` | `your_database` |
| `DB_PASSWORD` | Password kết nối database | `string` | `your_password` |
| **Database Connection Pool** |
| `DB_MAX_OPEN_CONNS` | Số lượng kết nối tối đa mở cùng lúc | `int` | `100` |
| `DB_MAX_IDLE_CONNS` | Số lượng kết nối nhàn rỗi tối đa | `int` | `10` |
| `DB_CONN_MAX_LIFETIME` | Thời gian tồn tại tối đa của 1 kết nối (giây) | `int` | `60` |
| `DB_CONN_MAX_IDLE_TIME` | Thời gian tối đa 1 kết nối ở trạng thái nhàn rỗi (giây) | `int` | `600` |
| `SLOW_THRESHOLD` | Ngưỡng thời gian để coi là truy vấn chậm (giây) | `float` | `5` |
| **REDIS CONFIGURATION** |
| `REDIS_HOST` | Hostname Redis server | `string` | `localhost` |
| `REDIS_PORT` | Cổng Redis server | `int` | `6379` |
| `REDIS_PASSWORD` | Password Redis (để trống nếu không có) | `string` | _empty_ |
| **TOKEN EXPIRATION** |
| `ACCESS_TOKEN_TTL` | Thời gian sống access token (giây) | `int` | `900` (15 phút) |
| `REFRESH_TOKEN_TTL` | Thời gian sống refresh token (giây) | `int` | `604800` (7 ngày) |
| **SESSION SETTINGS** |
| `SESSION_TTL` | Thời gian sống session (giây) | `int` | `604800` (7 ngày) |
| `IDLE_TIMEOUT` | Thời gian session hết hạn khi không hoạt động (giây) | `int` | `1800` (30 phút) |
| **BLACKLIST SETTINGS** |
| `REVOKED_TTL` | Thời gian lưu token đã thu hồi (giây) | `int` | `300` (5 phút) |
| **RATE LIMIT** |
| `RATE_LIMIT_PERIOD` | Chu kỳ giới hạn request (giây) | `int` | `60` |
| `RATE_LIMIT_MAX` | Số request tối đa trong mỗi chu kỳ | `int` | `100` |

---

## Ghi chú

- **Thời gian**: Tất cả giá trị TTL/time đều được tính bằng **giây**.
- **DB_DRIVER**: Các giá trị hợp lệ: `postgres`, `postgresql`, `mysql`, `mssql`, `sqlserver`.
- **CORS_ALLOW_ORIGINS**: Dùng `*` cho mọi origin, hoặc liệt kê cách nhau bởi dấu phẩy.
- **JWT_SECRET**: Nên thay đổi thành khóa bí mật mạnh trong môi trường production.