Đầy đủ và chi tiết về cấu trúc Swagger/OpenAPI 2.0, dễ hiểu và dễ tra cứu.
## Nguồn: https://swagger.io/specification/v2/
---

#  BẢNG TRA CỨU SWAGGER 2.0 (TIẾNG VIỆT)

## 1. ROOT OBJECT (ĐỐI TƯỢNG GỐC)

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `swagger` | string |  Có | Phiên bản Swagger đang dùng, phải là `"2.0"` |
| `info` | Info Object |  Có | Thông tin metadata của API |
| `host` | string |  Không | Tên miền hoặc IP của server (có thể kèm port) |
| `basePath` | string |  Không | Đường dẫn gốc của API, bắt đầu bằng `/` |
| `schemes` | [string] |  Không | Giao thức: `"http"`, `"https"`, `"ws"`, `"wss"` |
| `consumes` | [string] |  Không | Định dạng request nhận vào (global) |
| `produces` | [string] |  Không | Định dạng response trả ra (global) |
| `paths` | Paths Object |  Có | Định nghĩa các endpoint và operations |
| `definitions` | Definitions Object |  Không | Định nghĩa models (schemas) dùng chung |
| `parameters` | Parameters Object |  Không | Định nghĩa parameters dùng chung |
| `responses` | Responses Object |  Không | Định nghĩa responses dùng chung |
| `securityDefinitions` | Security Object |  Không | Định nghĩa các phương thức xác thực |
| `security` | [Security Object] |  Không | Áp dụng security cho toàn bộ API |
| `tags` | [Tag Object] |  Không | Danh sách tags để nhóm các operation |
| `externalDocs` | ExternalDocs Object |  Không | Tài liệu tham khảo bên ngoài |

---

## 2. INFO OBJECT

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `title` | string |  Có | Tên của ứng dụng/API |
| `description` | string |  Không | Mô tả ngắn về API (hỗ trợ Markdown) |
| `termsOfService` | string |  Không | Điều khoản sử dụng (URL) |
| `contact` | Contact Object |  Không | Thông tin liên hệ |
| `license` | License Object |  Không | Thông tin giấy phép |
| `version` | string |  Có | Phiên bản của API (không phải của spec) |

---

## 3. CONTACT OBJECT

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `name` | string |  Không | Tên người/tổ chức liên hệ |
| `url` | string |  Không | URL liên hệ (phải đúng format URL) |
| `email` | string |  Không | Email liên hệ (phải đúng format email) |

---

## 4. LICENSE OBJECT

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `name` | string |  Có | Tên giấy phép |
| `url` | string |  Không | URL giấy phép (phải đúng format URL) |

---

## 5. PATHS OBJECT

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `/{path}` | Path Item |  Có | Đường dẫn tương đối, bắt đầu bằng `/` |
| `x-*` | Any |  Không | Extension tùy chỉnh |

---

## 6. PATH ITEM OBJECT

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `$ref` | string |  Không | Tham chiếu đến định nghĩa bên ngoài |
| `get` | Operation |  Không | Định nghĩa operation GET |
| `put` | Operation |  Không | Định nghĩa operation PUT |
| `post` | Operation |  Không | Định nghĩa operation POST |
| `delete` | Operation |  Không | Định nghĩa operation DELETE |
| `options` | Operation |  Không | Định nghĩa operation OPTIONS |
| `head` | Operation |  Không | Định nghĩa operation HEAD |
| `patch` | Operation |  Không | Định nghĩa operation PATCH |
| `parameters` | [Parameter] |  Không | Parameters áp dụng cho tất cả operations trong path |
| `x-*` | Any |  Không | Extension tùy chỉnh |

---

## 7. OPERATION OBJECT

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `tags` | [string] |  Không | Tags để nhóm operation |
| `summary` | string |  Không | Tóm tắt ngắn (< 120 ký tự) |
| `description` | string |  Không | Mô tả chi tiết (hỗ trợ Markdown) |
| `externalDocs` | ExternalDocs |  Không | Tài liệu tham khảo bên ngoài |
| `operationId` | string |  Không | ID duy nhất cho operation |
| `consumes` | [string] |  Không | Định dạng request (ghi đè global) |
| `produces` | [string] |  Không | Định dạng response (ghi đè global) |
| `parameters` | [Parameter] |  Không | Danh sách parameters cho operation |
| `responses` | Responses Object |  Có | Các response có thể trả về |
| `schemes` | [string] |  Không | Giao thức (ghi đè global) |
| `deprecated` | boolean |  Không | Đánh dấu operation đã deprecated (mặc định: false) |
| `security` | [Security] |  Không | Xác thực cho operation (ghi đè global) |
| `x-*` | Any |  Không | Extension tùy chỉnh |

---

## 8. PARAMETER OBJECT

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `name` | string |  Có | Tên parameter (case-sensitive) |
| `in` | string |  Có | Vị trí: `"path"`, `"query"`, `"header"`, `"body"`, `"formData"` |
| `description` | string |  Không | Mô tả parameter |
| `required` | boolean |  Không | Bắt buộc? (path thì luôn true) |

### "body"`:

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `schema` | Schema Object |  Có | Schema định nghĩa kiểu dữ liệu của body |

### Nếu `in` KHÔNG phải `"body"`:

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `type` | string |  Có | Kiểu: `"string"`, `"number"`, `"integer"`, `"boolean"`, `"array"`, `"file"` |
| `format` | string |  Không | Format mở rộng (xem Data Type Formats) |
| `allowEmptyValue` | boolean |  Không | Cho phép giá trị rỗng (mặc định: false) |
| `items` | Items Object |  Không | Nếu type là `"array"`, định nghĩa kiểu của items |
| `collectionFormat` | string |  Không | Format array: `"csv"`, `"ssv"`, `"tsv"`, `"pipes"`, `"multi"` (mặc định: `"csv"`) |
| `default` | * |  Không | Giá trị mặc định |
| `maximum` | number |  Không | Giá trị tối đa |
| `exclusiveMaximum` | boolean |  Không | Loại trừ giá trị tối đa |
| `minimum` | number |  Không | Giá trị tối thiểu |
| `exclusiveMinimum` | boolean |  Không | Loại trừ giá trị tối thiểu |
| `maxLength` | integer |  Không | Độ dài tối đa (string) |
| `minLength` | integer |  Không | Độ dài tối thiểu (string) |
| `pattern` | string |  Không | Regex pattern (string) |
| `maxItems` | integer |  Không | Số lượng items tối đa (array) |
| `minItems` | integer |  Không | Số lượng items tối thiểu (array) |
| `uniqueItems` | boolean |  Không | Các items có unique không? |
| `enum` | [*] |  Không | Danh sách giá trị cho phép |
| `multipleOf` | number |  Không | Bội số của giá trị |

---

## 9. RESPONSES OBJECT

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `default` | Response Object |  Không | Response mặc định cho các status code không được định nghĩa |
| `{status code}` | Response Object |  Không | Response cho HTTP status code cụ thể (VD: `"200"`, `"404"`) |
| `x-*` | Any |  Không | Extension tùy chỉnh |

---

## 10. RESPONSE OBJECT

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `description` | string |  Có | Mô tả response (hỗ trợ Markdown) |
| `schema` | Schema Object |  Không | Cấu trúc dữ liệu của response |
| `headers` | Headers Object |  Không | Danh sách headers trong response |
| `examples` | Example Object |  Không | Ví dụ response |
| `x-*` | Any |  Không | Extension tùy chỉnh |

---

## 11. SCHEMA OBJECT

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `$ref` | string |  Không | Tham chiếu JSON Pointer |
| `format` | string |  Không | Format mở rộng |
| `title` | string |  Không | Tiêu đề |
| `description` | string |  Không | Mô tả (hỗ trợ Markdown) |
| `default` | * |  Không | Giá trị mặc định |
| `multipleOf` | number |  Không | Bội số |
| `maximum` | number |  Không | Giá trị tối đa |
| `exclusiveMaximum` | boolean |  Không | Loại trừ max |
| `minimum` | number |  Không | Giá trị tối thiểu |
| `exclusiveMinimum` | boolean |  Không | Loại trừ min |
| `maxLength` | integer |  Không | Độ dài tối đa |
| `minLength` | integer |  Không | Độ dài tối thiểu |
| `pattern` | string |  Không | Regex pattern |
| `maxItems` | integer |  Không | Số items tối đa |
| `minItems` | integer |  Không | Số items tối thiểu |
| `uniqueItems` | boolean |  Không | Các items có unique không? |
| `maxProperties` | integer |  Không | Số properties tối đa |
| `minProperties` | integer |  Không | Số properties tối thiểu |
| `required` | [string] |  Không | Danh sách properties bắt buộc |
| `enum` | [*] |  Không | Danh sách giá trị cho phép |
| `type` | string |  Không | Kiểu dữ liệu: `"object"`, `"array"`, `"string"`, `"number"`, `"integer"`, `"boolean"` |
| `items` | Schema Object |  Không | Nếu type là `"array"`, định nghĩa items |
| `allOf` | [Schema] |  Không | Kế thừa/composition |
| `properties` | Map |  Không | Các properties của object |
| `additionalProperties` | boolean/Schema |  Không | Cho phép properties bổ sung |
| `discriminator` | string |  Không | Hỗ trợ polymorphism |
| `readOnly` | boolean |  Không | Property chỉ đọc (mặc định: false) |
| `xml` | XML Object |  Không | Metadata XML |
| `externalDocs` | ExternalDocs |  Không | Tài liệu tham khảo |
| `example` | Any |  Không | Ví dụ giá trị |

---

## 12. TAG OBJECT

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `name` | string |  Có | Tên tag |
| `description` | string |  Không | Mô tả tag (hỗ trợ Markdown) |
| `externalDocs` | ExternalDocs |  Không | Tài liệu tham khảo |

---

## 13. SECURITY SCHEME OBJECT

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `type` | string |  Có | Loại: `"basic"`, `"apiKey"`, `"oauth2"` |
| `description` | string |  Không | Mô tả |
| `name` | string |  Có (apiKey) | Tên header hoặc query parameter |
| `in` | string |  Có (apiKey) | Vị trí: `"query"` hoặc `"header"` |
| `flow` | string |  Có (oauth2) | Loại flow: `"implicit"`, `"password"`, `"application"`, `"accessCode"` |
| `authorizationUrl` | string |  Có (oauth2) | URL xác thực |
| `tokenUrl` | string |  Có (oauth2) | URL token |
| `scopes` | Scopes Object |  Có (oauth2) | Danh sách scopes |

---

## 14. EXTERNAL DOCUMENTATION OBJECT

| Tên trường | Loại | Bắt buộc | Mô tả |
|-----------|------|----------|-------|
| `description` | string |  Không | Mô tả ngắn |
| `url` | string |  Có | URL tài liệu (phải đúng format URL) |

---

##  CÁC KIỂU DỮ LIỆU (DATA TYPES)

| Kiểu | Mô tả | Format thường dùng |
|------|-------|-------------------|
| `string` | Chuỗi | `"date"`, `"date-time"`, `"email"`, `"uuid"`, `"password"` |
| `number` | Số thực | `"float"`, `"double"` |
| `integer` | Số nguyên | `"int32"`, `"int64"` |
| `boolean` | Boolean | - |
| `array` | Mảng | - |
| `object` | Object | - |
| `file` | File (chỉ dùng trong formData) | - |

---

##  CÁC LOẠI `in` (VỊ TRÍ PARAMETER)

| Giá trị | Mô tả |
|---------|-------|
| `path` | Tham số trong URL (`/users/{id}`) |
| `query` | Tham số sau dấu `?` (`?page=1`) |
| `header` | Tham số trong header (`Authorization`) |
| `body` | Payload trong body (POST, PUT) |
| `formData` | Tham số dạng form (multipart/form-data) |

---

##  CÁC LOẠI `collectionFormat` (FORMAT ARRAY)

| Giá trị | Mô tả | Ví dụ |
|---------|-------|-------|
| `csv` | Comma separated values | `foo,bar` |
| `ssv` | Space separated values | `foo bar` |
| `tsv` | Tab separated values | `foo\tbar` |
| `pipes` | Pipe separated values | `foo\|bar` |
| `multi` | Multiple parameters | `foo=bar&foo=baz` |

---

