
---

##  "Vị trí các token" là gì?

**Token** là **đơn vị từ vựng nhỏ nhất** trong mã nguồn. Khi trình biên dịch (hoặc `parser`) đọc code, nó không hiểu code như con người đọc văn bản. Nó tách code thành các **token** riêng lẻ.

**"Vị trí các token"** = **dòng nào, cột nào, file nào** mà token đó xuất hiện.

---

##  Ví dụ cụ thể

### Code Go:
```go
package main

import "fmt"

func main() {
    x := 10
    fmt.Println(x)
}
```

### Quá trình tách thành token:

| Token | Giá trị | Vị trí |
|-------|---------|--------|
| `package` | từ khóa | dòng 1, cột 1 |
| `main` | tên | dòng 1, cột 9 |
| `import` | từ khóa | dòng 3, cột 1 |
| `"fmt"` | chuỗi literal | dòng 3, cột 8 |
| `func` | từ khóa | dòng 5, cột 1 |
| `main` | tên hàm | dòng 5, cột 6 |
| `(` | dấu ngoặc | dòng 5, cột 10 |
| `)` | dấu ngoặc | dòng 5, cột 11 |
| `{` | dấu ngoặc nhọn | dòng 5, cột 13 |
| `x` | tên biến | dòng 6, cột 5 |
| `:=` | toán tử gán | dòng 6, cột 7 |
| `10` | số literal | dòng 6, cột 10 |
| `fmt` | tên package | dòng 7, cột 2 |
| `.` | dấu chấm | dòng 7, cột 5 |
| `Println` | tên hàm | dòng 7, cột 6 |
| `(` | dấu ngoặc | dòng 7, cột 13 |
| `x` | tên biến | dòng 7, cột 14 |
| `)` | dấu ngoặc | dòng 7, cột 15 |
| `}` | dấu ngoặc nhọn | dòng 8, cột 1 |

---

##  FileSet lưu vị trí các token để làm gì?

### **Báo lỗi chính xác**

```go
// user.go:10:5: undefined: User
//       ^^^^
//       dòng 10, cột 5
```

Nếu không có FileSet, lỗi sẽ là: `undefined: User` – không biết ở đâu!

### **Debug và log**

```go
// In ra vị trí của từng node trong AST
fmt.Println("Tìm thấy struct tại:", fset.Position(node.Pos()))
// Output: Tìm thấy struct tại: user.go:15:1
```

### **Linting và tool phân tích**

Các tool như `golangci-lint` dùng FileSet để đưa ra cảnh báo kèm vị trí cụ thể.

---

## Cách FileSet lưu trữ vị trí

### `token.FileSet` hoạt động như một cuốn sổ cái:

```go
fset := token.NewFileSet()

// Khi parse file, nó thêm file vào FileSet
f, err := parser.ParseFile(fset, "user.go", nil, 0)

// Bên trong FileSet có 1 map:
// file "user.go" -> { dòng: 1-100, cột: 1-80, ký tự: 0-5000 }
```

### Mỗi node trong AST có `Pos()`:

```go
type Node interface {
    Pos() token.Pos // Vị trí bắt đầu của node
    End() token.Pos // Vị trí kết thúc của node
}
```

`token.Pos` là một **số nguyên** đại diện cho vị trí tuyệt đối trong FileSet.

### Chuyển đổi `token.Pos` thành vị trí dễ đọc:

```go
pos := node.Pos()
position := fset.Position(pos)
fmt.Printf("File: %s, Line: %d, Column: %d\n",
    position.Filename, position.Line, position.Column)
// Output: File: user.go, Line: 15, Column: 1
```

---

## Ví dụ thực tế trong code của bro

```go
fset := token.NewFileSet()

ast.Inspect(f, func(n ast.Node) bool {
    switch x := n.(type) {
    case *ast.TypeSpec:
        // In vị trí của struct
        pos := fset.Position(x.Pos())
        fmt.Printf(" Struct %s at %s:%d\n", x.Name.Name, pos.Filename, pos.Line)

        // In vị trí kết thúc
        end := fset.Position(x.End())
        fmt.Printf("   Ends at %s:%d\n", end.Filename, end.Line)
    }
    return true
})
```

**Output:**
```
 Struct User at internal/dto/user.go:5
   Ends at internal/dto/user.go:12
 Struct Product at internal/dto/product.go:15
   Ends at internal/dto/product.go:20
```

---

##  Tổng quan cấu trúc FileSet

```
token.FileSet
├── file1: "user.go"
│   ├── Dòng 1: position 1-10
│   ├── Dòng 2: position 11-20
│   └── Dòng 3: position 21-30
├── file2: "product.go"
│   ├── Dòng 1: position 31-40
│   └── Dòng 2: position 41-50
└── file3: "order.go"
    └── ...
```

Mỗi `token.Pos` là một số nguyên đại diện cho vị trí **tuyệt đối** trong toàn bộ FileSet. Ví dụ: `token.Pos(123)` có thể là "dòng 10, cột 5 của file user.go".

---

## Tóm tắt

| Khái niệm | Ý nghĩa |
|-----------|---------|
| **Token** | Đơn vị từ vựng nhỏ nhất trong code (từ khóa, tên, số, dấu, ...) |
| **Vị trí token** | Dòng, cột, file mà token xuất hiện |
| **FileSet** | Quản lý vị trí của tất cả token trong nhiều file |
| **token.Pos** | Số nguyên đại diện cho vị trí tuyệt đối |
| **fset.Position(pos)** | Chuyển token.Pos thành tên file + dòng + cột |
| **Công dụng** | Báo lỗi, debug, tool phân tích, sinh code |

---

##  Mẹo debug AST

```go
// In toàn bộ AST với vị trí
ast.Print(fset, f)
```

Hoặc:

```go
// In vị trí của từng node
ast.Inspect(f, func(n ast.Node) bool {
    if n != nil {
        fmt.Printf("%T at %s\n", n, fset.Position(n.Pos()))
    }
    return true
})
```

---
