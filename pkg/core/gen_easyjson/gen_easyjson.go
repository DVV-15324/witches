package gen_easyjson

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FileSet đánh dấu như 1 cái sổ cái lưu token
func GeneratorEasyJson(fset *token.FileSet, filePathWalk string) {
	structs := []string{}
	// filepath.Walk: Hàm duyệt đệ quy tất cả file và thư mục con trong path (đường dẫn) ./internal/dto
	err := filepath.Walk(filePathWalk, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		// Kiểm tra chỉ chọn xem có file .go trước
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Parse (Phân tích) file .go sang cây AST
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil
		}

		// Kiểm tra struct có implement interface có method Gen() không?
		// Để đơn giản, ta tìm các struct có method receiver là Gen()
		// (có thể parse thêm interface nếu cần)
		hasGenMethod := false
		ast.Inspect(f, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok {
				if fn.Name.Name == "Gen" && fn.Recv != nil {
					hasGenMethod = true
				}
			}
			return true
		})

		if !hasGenMethod {
			return nil
		}

		// Tìm các struct trong file
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				if _, ok := x.Type.(*ast.StructType); ok {
					structName := x.Name.Name
					structs = append(structs, structName)
					// Gọi easyjson CLI
					generateEasyJSON(path)
				}
			}
			return true
		})
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Generated easyjson for structs in files:", structs)
}

func generateEasyJSON(filePath string) {
	// Chạy easyjson -all <file>
	cmd := exec.Command("easyjson", "-all", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: easyjson failed for %s: %v\n", filePath, err)
	}
}
