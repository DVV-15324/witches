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

func GeneratorEasyJson(fset *token.FileSet, input string, output string) error {
	if input == "" {
		return fmt.Errorf("input path is empty")
	}

	info, err := os.Stat(input)
	if err != nil {
		return fmt.Errorf("input path does not exist: %s", input)
	}

	var workDir string
	var isFile bool

	if info.IsDir() {
		workDir = input
		isFile = false
	} else {
		workDir = filepath.Dir(input)
		isFile = true
	}

	outputDir := output
	if outputDir == "" {
		outputDir = workDir
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	absInput, err := filepath.Abs(input)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}
	absOutput, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute output path: %v", err)
	}

	// fmt.Printf("Input: %s\n", absInput)
	// fmt.Printf("Output: %s\n", absOutput)

	var targetFiles []string

	if isFile {
		targetFiles = append(targetFiles, absInput)
	} else {
		err = filepath.Walk(absInput, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if strings.HasSuffix(path, "_test.go") || !strings.HasSuffix(path, ".go") {
				return nil
			}

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
				//fmt.Printf("  Found marker in: %s\n", path)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("error walking: %v", err)
		}
	}

	// if len(targetFiles) == 0 {
	// 	return fmt.Errorf("no files with GenEasyJson() method found")
	// }

	fmt.Printf("Generating easyjson for %d file(s)...\n", len(targetFiles))
	var generatedFiles []string

	for _, path := range targetFiles {
		relPath, err := filepath.Rel(absInput, path)
		if err != nil {
			relPath = filepath.Base(path)
		}
		outputPath := filepath.Join(absOutput, relPath)

		outputDirPath := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDirPath, 0755); err != nil {
			fmt.Printf("  Warning: failed to create output dir: %v\n", err)
			continue
		}

		// Nếu output khác input, copy file sang output
		if outputPath != path {
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("  Warning: failed to read %s: %v\n", path, err)
				continue
			}
			if err := os.WriteFile(outputPath, content, 0644); err != nil {
				fmt.Printf("  Warning: failed to write output file: %v\n", err)
				continue
			}
		}

		// Xóa file gen cũ
		genOutputFile := strings.TrimSuffix(outputPath, ".go") + "_easyjson.go"
		os.Remove(genOutputFile)

		// Xóa file tmp
		tmpFile := strings.TrimSuffix(outputPath, ".go") + "_easyjson.go.tmp"
		os.Remove(tmpFile)

		// Gen TRỰC TIẾP trên file outputPath
		if err := generateEasyJSON(outputPath); err != nil {
			fmt.Printf("  Warning: easyjson failed for %s: %v\n", path, err)
			continue
		}

		if _, err := os.Stat(genOutputFile); err == nil {
			generatedFiles = append(generatedFiles, genOutputFile)
			fmt.Printf("  Generated: %s\n", genOutputFile)
		}

		os.Remove(tmpFile)
	}

	// if len(generatedFiles) == 0 {
	// 	return fmt.Errorf("no files were generated successfully")
	// }

	fmt.Printf("Successfully generated %d file(s)\n", len(generatedFiles))
	return nil
}

func generateEasyJSON(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	cmd := exec.Command("easyjson", "-all", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
