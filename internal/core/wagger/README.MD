# Hướng dẫn Swagger Annotations cho người mới

Swagger (OpenAPI) là công cụ tự động sinh tài liệu API từ code. Trong Go, thư viện [swaggo/swag](https://github.com/swaggo/swag) sử dụng annotation dạng `// @tên`.

## Cài đặt nhanh

`bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
`

## Cấu trúc annotation chính

Annotations đặt trên hàm main hoặc handler.

`go
// @title Swagger Example API
// @version 1.0
// @description Đây là API mẫu.
// @host localhost:8080
// @BasePath /api/v1
func main() {}
`

## Bảng giải thích các annotation

| Annotation | Ý nghĩa | Ví dụ |
|------------|---------|-------|
| title | Bắt buộc. Tiêu đề API | `// @title Swagger Example API` |
| version | Bắt buộc. Phiên bản API | `// @version 1.0` |
| description | Mô tả API | `// @description This is a sample server` |
| tag.name | Tên nhóm API | `// @tag.name User` |
| tag.description | Mô tả tag | `// @tag.description User APIs` |
| tag.docs.url | Link tài liệu tag | `// @tag.docs.url https://example.com` |
| tag.docs.description | Mô tả docs tag | `// @tag.docs.description API docs` |
| termsOfService | Điều khoản sử dụng | `// @termsOfService http://swagger.io/terms/` |
| contact.name | Tên liên hệ | `// @contact.name API Support` |
| contact.url | Link hỗ trợ | `// @contact.url http://www.swagger.io/support` |
| contact.email | Email liên hệ | `// @contact.email support@swagger.io` |
| license.name | Tên license | `// @license.name Apache 2.0` |
| license.url | Link license | `// @license.url http://www.apache.org/licenses/` |
| host | Domain/IP chạy API | `// @host localhost:8080` |
| BasePath | Base URL API | `// @BasePath /api/v1` |
| accept | Input format | `// @accept json` |
| produce | Output format | `// @produce json` |
| query.collection.format | Format array query | `// @query.collection.format multi` |
| schemes | Giao thức HTTP/HTTPS | `// @schemes http https` |
| externalDocs.description | Mô tả docs ngoài | `// @externalDocs.description OpenAPI` |
| externalDocs.url | Link docs ngoài | `// @externalDocs.url https://swagger.io` |
| x-* | Custom field mở rộng | `// @x-example-key {"key":"value"}` |

## Sử dụng annotation trong handler (router)

`go
// GetUser godoc
// @Summary Lấy thông tin user
// @Description Trả về chi tiết user theo ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {}
`
V
## Mẹo cho người mới

1. Luôn ghi `title` và `version` – nếu thiếu, swag sẽ báo lỗi.
2. Sử dụng `@tag.name` và `@Tags` để nhóm các API liên quan trên giao diện Swagger UI.
3. Các annotation như `@accept` / `@produce` không bắt buộc nhưng nên có để UI hiển thị đúng.
4. Chạy `swag init` mỗi khi thay đổi annotation để cập nhật file `docs/swagger.json` và `docs/swagger.yaml`.
5. Truy cập Swagger UI tại `http://localhost:8080/swagger/index.html` (nếu đã tích hợp swaggo/gin-swagger).

## Tham khảo thêm

- [swaggo/swag GitHub](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)