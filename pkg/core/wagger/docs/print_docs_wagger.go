package docs

import "fmt"

func PrintDocsSwaggerStructTag() {
	fmt.Print(`
+--------------------+----------------------------------+----------------------------------+
| Tag                | Ý nghĩa                          | Ví dụ                            |
+--------------------+----------------------------------+----------------------------------+
| json               | Tên field trong JSON            | json:"name"                      |
| omitempty          | Bỏ qua field nếu rỗng           | json:"age,omitempty"             |
| description        | Mô tả field trong Swagger       | description:"User Name"          |
| example            | Giá trị mẫu                     | example:"Dinh Viet Vu"           |
| binding:"required" | Bắt buộc khi bind request       | binding:"required"               |
| maxLength          | Độ dài chuỗi tối đa             | maxLength:"10"                   |
| minLength          | Độ dài chuỗi tối thiểu          | minLength:"1"                    |
| maximum            | Giá trị số lớn nhất             | maximum:"100"                    |
| minimum            | Giá trị số nhỏ nhất             | minimum:"0"                      |
| default            | Giá trị mặc định                | default:"guest"                  |
| enums              | Danh sách giá trị hợp lệ        | enums:"admin,user,guest"         |
| format             | Định dạng dữ liệu               | format:"email"                   |
| validate           | Validation validator.v10        | validate:"email"                 |
+--------------------+----------------------------------+----------------------------------+
`)
}

func PrintDocsSwaggerMetaData() {
	fmt.Print(`
// @Summary
// @Description
// @Tags
// @Accept json
// @Produce json

// @Param id path string false "Path Parameter"
// @Param request body object false "Request Body"
// @Param keyword query string false "Query Parameter"
// @Param Authorization header string false "Bearer Token"

// @Success 200 {object} object "Success Response"
// @Failure 400 {object} object "Bad Request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 404 {object} object "Not Found"
// @Failure 500 {object} object "Internal Server Error"

// @Router /resource [get]
`)
}
