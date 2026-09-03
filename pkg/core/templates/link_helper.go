package template

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func ValidateModuleStructure(templateFiles map[string]string, config ModuleConfig) error {
	required := []string{
		fmt.Sprintf("internal/%s/model/model.go", config.Name),
		fmt.Sprintf("internal/%s/handler/handler.go", config.Name),
		fmt.Sprintf("internal/%s/repository/repository.go", config.Name),
		fmt.Sprintf("internal/%s/usecase/usecase.go", config.Name),
		fmt.Sprintf("internal/%s/module.go", config.Name),
	}
	moduleFiles := make(map[string][]string)
	migrationFiles := []string{}
	hasSharedDomain := false
	for path := range templateFiles {
		// Kiểm tra migration files - nằm trong domain
		pathMigrate := fmt.Sprintf("internal/%s/migrate/migrations", config.Name)
		if strings.HasPrefix(path, pathMigrate) {
			migrationFiles = append(migrationFiles, path)
			continue
		}
		// Kiểm tra shared domain
		sharedDomainPath := fmt.Sprintf("internal/shared/domain/%s.go", config.Name)
		if path == sharedDomainPath {
			hasSharedDomain = true
			continue
		}
		// Kiểm tra internal files
		parts := strings.Split(strings.TrimPrefix(path, "internal/"), "/")
		if len(parts) < 2 {
			continue
		}
		module := parts[0]
		filename := strings.Join(parts[1:], "/")
		moduleFiles[module] = append(moduleFiles[module], filename)
	}
	// Validate internal files
	for module, files := range moduleFiles {
		hasRequiredFiles := true
		for _, requiredFile := range required {
			expected := strings.TrimPrefix(requiredFile, "internal/"+module+"/")
			found := false

			for _, f := range files {
				if f == expected {
					found = true
					break
				}
			}

			if !found {
				hasRequiredFiles = false
				break
			}
		}

		if hasRequiredFiles {
			fmt.Printf("  Module '%s' has all required files\n", module)
		} else {
			fmt.Printf("  Module '%s' missing some files\n", module)
		}
	}
	// Check migration files
	if len(migrationFiles) > 0 {
		fmt.Printf("  Found %d migration files in internal/%s/migrate/migrations/\n", len(migrationFiles), config.Name)
		for _, mf := range migrationFiles {
			fmt.Printf("    - %s\n", mf)
		}
	} else {
		fmt.Printf("  No migration files found in internal/%s/migrate/migrations/\n", config.Name)
	}
	// Check shared domain
	if hasSharedDomain {
		fmt.Printf("  Found internal/shared/domain/%s.go\n", config.Name)
	}
	if len(moduleFiles) == 0 && len(migrationFiles) == 0 && !hasSharedDomain {
		return fmt.Errorf("no module files, migration files, or shared domain found")
	}
	return nil
}

func RewriteModuleImports(filePath, sourceModule, targetModule string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return err
	}
	changed := false
	for _, imp := range node.Imports {
		oldPath := strings.Trim(imp.Path.Value, `"`)
		newPath := NormalizeImportPath(oldPath, sourceModule, targetModule)
		if newPath != oldPath {
			imp.Path.Value = fmt.Sprintf(`"%s"`, newPath)
			changed = true
			fmt.Printf("    %s → %s\n", oldPath, newPath)
		}
	}
	if !changed {
		return nil
	}
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return err
	}
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func FixAllImports(projectDir, sourceModule, targetModule string) error {
	// Nếu targetModule rỗng, lấy từ go.mod
	if targetModule == "" {
		var err error
		targetModule, err = GetCurrentModule(projectDir)
		if err != nil {
			return fmt.Errorf("failed to get target module: %w", err)
		}
		fmt.Printf("  Using target module from go.mod: %s\n", targetModule)
	}
	fmt.Printf("Fixing all imports from '%s' to '%s'...\n", sourceModule, targetModule)
	var filesProcessed int
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip các thư mục không cần thiết
			skipDirs := []string{".git", "logs", "swagger", "vendor", "node_modules", "tmp"}
			for _, skip := range skipDirs {
				if strings.Contains(path, skip) {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		// Kiểm tra nếu file chứa sourceModule cần sửa
		if sourceModule != "" && !strings.Contains(string(content), sourceModule) {
			return nil
		}
		filesProcessed++
		fmt.Printf("  Fixing: %s\n", filepath.Base(path))
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, content, parser.ParseComments)
		if err != nil {
			return err
		}
		changed := false
		for _, imp := range node.Imports {
			oldPath := strings.Trim(imp.Path.Value, `"`)
			newPath := NormalizeImportPath(oldPath, sourceModule, targetModule)
			if newPath != oldPath {
				imp.Path.Value = fmt.Sprintf(`"%s"`, newPath)
				changed = true
				fmt.Printf("    %s → %s\n", oldPath, newPath)
			}
		}
		if changed {
			var buf bytes.Buffer
			if err := format.Node(&buf, fset, node); err != nil {
				return err
			}
			return os.WriteFile(path, buf.Bytes(), 0644)
		}
		return nil
	})
	fmt.Printf("  Fixed %d files\n", filesProcessed)
	return err
}

// GetModuleNameFromRepo lấy module name từ repo URL
func GetModuleNameFromRepo(repoURL string) string {
	repoURL = strings.TrimPrefix(repoURL, "https://")
	repoURL = strings.TrimPrefix(repoURL, "http://")
	repoURL = strings.TrimSuffix(repoURL, ".git")
	return repoURL
}

// GetCurrentModule lấy module name từ go.mod của project
func GetCurrentModule(projectDir string) (string, error) {
	goModPath := filepath.Join(projectDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			return moduleName, nil
		}
	}
	return "", fmt.Errorf("module directive not found in go.mod")
}

// NormalizeImportPath chuẩn hóa import path dựa trên source và target module
func NormalizeImportPath(oldPath, sourceModule, targetModule string) string {
	if oldPath == "" || targetModule == "" {
		return oldPath
	}

	// GIỮ NGUYÊN EXTERNAL PKG (chứa /pkg/)
	// Ví dụ: github.com/DVV-15324/witches-book/pkg/redis -> giữ nguyên
	if strings.Contains(oldPath, "/pkg/") {
		return oldPath
	}

	// GIỮ NGUYÊN CÁC THƯ VIỆN EXTERNAL KHÁC
	externalPrefixes := []string{
		"github.com/",
		"gitlab.com/",
		"golang.org/",
		"go.uber.org/",
		"google.golang.org/",
		"gopkg.in/",
	}

	// Kiểm tra nếu là external library (không phải sourceModule)
	isExternal := false
	for _, prefix := range externalPrefixes {
		if strings.HasPrefix(oldPath, prefix) {
			// Nếu là sourceModule (witches-book) thì KHÔNG skip
			if sourceModule != "" && strings.Contains(oldPath, sourceModule) {
				// Kiểm tra nếu có /internal/ thì vẫn sửa
				if strings.Contains(oldPath, "/internal/") {
					isExternal = false
					break
				}

				// Kiểm tra nếu có /cmd/ thì vẫn sửa
				if strings.Contains(oldPath, "/cmd/") {
					isExternal = false
					break
				}

				// Nếu là /pkg/ thì skip (đã xử lý ở trên)
				if strings.Contains(oldPath, "/pkg/") {
					isExternal = true
					break
				}
			}
			isExternal = true
			break
		}
	}
	if isExternal {
		return oldPath
	}

	// THAY THẾ SOURCE MODULE -> TARGET MODULE (cho internal code)
	if sourceModule != "" && strings.Contains(oldPath, sourceModule) {
		return strings.Replace(oldPath, sourceModule, targetModule, -1)
	}

	// XỬ LÝ IMPORT TƯƠNG ĐỐI (internal/, cmd/, pkg/)
	if !strings.HasPrefix(oldPath, targetModule) {
		if strings.HasPrefix(oldPath, "internal/") ||
			strings.HasPrefix(oldPath, "cmd/") ||
			strings.HasPrefix(oldPath, "pkg/") {
			return targetModule + "/" + oldPath
		}
	}

	return oldPath
}
