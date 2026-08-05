// pkg/core/easyjson/generator.go
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

// GeneratorEasyJson - Tạo file easyjson cho các struct có GenEasyJson()
//   - input: đường dẫn đến file hoặc thư mục cần gen
//   - output: thư mục đích để lưu file gen (rỗng = cùng thư mục với input)
//   - fset: token.FileSet (có thể nil)
func GeneratorEasyJson(fset *token.FileSet, input string, output string) error {
	//  1. KIỂM TRA ĐẦU VÀO
	if input == "" {
		return fmt.Errorf("input path is empty")
	}

	info, err := os.Stat(input)
	if err != nil {
		return fmt.Errorf("input path does not exist: %s", input)
	}

	// Xác định thư mục làm việc
	var workDir string
	var isFile bool

	if info.IsDir() {
		workDir = input
		isFile = false
	} else {
		workDir = filepath.Dir(input)
		isFile = true
	}

	// 2. XÁC ĐỊNH ĐẦU RA
	outputDir := output
	if outputDir == "" {
		// Mặc định: cùng thư mục với input
		outputDir = workDir
	}

	// Tạo thư mục output nếu chưa có
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	//  3. CHUẨN HÓA ĐƯỜNG DẪN
	absInput, err := filepath.Abs(input)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}
	absOutput, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute output path: %v", err)
	}

	fmt.Printf("Input: %s\n", absInput)
	fmt.Printf("Output: %s\n", absOutput)

	//  4. TÌM CÁC FILE CẦN GEN
	var targetFiles []string

	if isFile {
		// Nếu input là file, chỉ xử lý file đó
		targetFiles = append(targetFiles, absInput)
	} else {
		// Nếu input là thư mục, walk để tìm file
		err = filepath.Walk(absInput, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			// Bỏ qua file test
			if strings.HasSuffix(path, "_test.go") {
				return nil
			}

			if !strings.HasSuffix(path, ".go") {
				return nil
			}

			// Parse file để kiểm tra marker
			localFset := token.NewFileSet()
			if fset != nil {
				localFset = fset
			}

			f, err := parser.ParseFile(localFset, path, nil, parser.ParseComments)
			if err != nil {
				fmt.Printf("  Warning: cannot parse %s: %v\n", path, err)
				return nil
			}

			hasMarker := false
			ast.Inspect(f, func(n ast.Node) bool {
				if fn, ok := n.(*ast.FuncDecl); ok {
					if fn.Name.Name == "GenEasyJson" && fn.Recv != nil {
						hasMarker = true
						return false
					}
				}
				return true
			})

			if hasMarker {
				targetFiles = append(targetFiles, path)
				fmt.Printf("  Found marker in: %s\n", path)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("error walking: %v", err)
		}
	}

	if len(targetFiles) == 0 {
		return fmt.Errorf("no files with GenEasyJson() method found")
	}

	// 5. GỌI EASYJSON CHO TỪNG FILE
	fmt.Printf("Generating easyjson for %d file(s)...\n", len(targetFiles))
	var generatedFiles []string

	for _, path := range targetFiles {
		// Xác định output path
		relPath, err := filepath.Rel(absInput, path)
		if err != nil {
			relPath = filepath.Base(path)
		}
		outputPath := filepath.Join(absOutput, relPath)

		// Đảm bảo thư mục output tồn tại
		outputDirPath := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDirPath, 0755); err != nil {
			fmt.Printf("  Warning: failed to create output dir: %v\n", err)
			continue
		}

		// Tạo file tạm để gen
		tmpFile := outputPath
		if tmpFile == path {
			// Nếu output trùng input, tạo file tạm riêng
			tmpFile = strings.TrimSuffix(path, ".go") + "_tmp.go"
			defer os.Remove(tmpFile)
		}

		// Copy nội dung từ input sang tmp
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("  Warning: failed to read %s: %v\n", path, err)
			continue
		}
		if err := os.WriteFile(tmpFile, content, 0644); err != nil {
			fmt.Printf("  Warning: failed to write tmp file: %v\n", err)
			continue
		}

		// Chạy easyjson
		if err := generateEasyJSON(tmpFile); err != nil {
			fmt.Printf("  Warning: easyjson failed for %s: %v\n", path, err)
			continue
		}

		// Di chuyển file gen từ tmp sang output
		genTmpFile := strings.TrimSuffix(tmpFile, ".go") + "_easyjson.go"
		genOutputFile := strings.TrimSuffix(outputPath, ".go") + "_easyjson.go"

		if _, err := os.Stat(genTmpFile); err == nil {
			if err := os.Rename(genTmpFile, genOutputFile); err != nil {
				fmt.Printf("  Warning: failed to move generated file: %v\n", err)
				continue
			}
			generatedFiles = append(generatedFiles, genOutputFile)
			fmt.Printf("  Generated: %s\n", genOutputFile)
		}
	}

	// 6. KẾT QUẢ
	if len(generatedFiles) == 0 {
		return fmt.Errorf("no files were generated successfully")
	}

	fmt.Printf("Successfully generated %d file(s)\n", len(generatedFiles))
	return nil
}

func generateEasyJSON(filePath string) error {
	// Kiểm tra file tồn tại
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	cmd := exec.Command("easyjson", "-all", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
