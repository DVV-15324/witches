package easyjson

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

func GeneratorEasyJson(fset *token.FileSet, filePathWalk string) {
	var targetFiles []string // 1. Tạo danh sách file cần gen

	err := filepath.Walk(filePathWalk, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return nil
		}

		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil
		}

		hasMarker := false
		// 2. Kiểm tra có method TableName không
		ast.Inspect(f, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok {
				if fn.Name.Name == "TableName" && fn.Recv != nil {
					hasMarker = true
					return false // Tìm thấy rồi -> dừng duyệt tiếp
				}
			}
			return true
		})

		// 3. Nếu có dấu hiệu, lưu file vào danh sách
		if hasMarker {
			targetFiles = append(targetFiles, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking:", err)
		return
	}

	// 4. CHỈ GỌI EASYJSON 1 LẦN DUY NHẤT cho mỗi FILE
	for _, path := range targetFiles {
		fmt.Printf("Generating easyjson for file: %s\n", path)
		generateEasyJSON(path)
	}
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
