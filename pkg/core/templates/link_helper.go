package template

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func ValidateModuleStructure(templateFiles map[string]string, config DomainConfig) error {
	// Các file bắt buộc cho domain
	required := []string{
		fmt.Sprintf("internal/%s/model/model.go", config.Name),
		fmt.Sprintf("internal/%s/handler/handler.go", config.Name),
		fmt.Sprintf("internal/%s/repository/repository.go", config.Name),
		fmt.Sprintf("internal/%s/usecase/usecase.go", config.Name),
		fmt.Sprintf("internal/%s/module.go", config.Name),
	}

	domainFiles := make(map[string][]string)
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

		domain := parts[0]
		filename := strings.Join(parts[1:], "/")
		domainFiles[domain] = append(domainFiles[domain], filename)
	}

	// Validate internal files
	for domain, files := range domainFiles {
		hasRequiredFiles := true
		for _, requiredFile := range required {
			expected := strings.TrimPrefix(requiredFile, "internal/"+domain+"/")
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
			fmt.Printf("Domain '%s' has all required files\n", domain)
		} else {
			fmt.Printf("Domain '%s' missing some files\n", domain)
		}
	}

	// Check migration files
	if len(migrationFiles) > 0 {
		fmt.Printf("Found %d migration files in internal/%s/migrate/migrations/\n", len(migrationFiles), config.Name)
		for _, mf := range migrationFiles {
			fmt.Printf("     - %s\n", mf)
		}
	} else {
		fmt.Printf("No migration files found in internal/%s/migrate/migrations/\n", config.Name)
	}

	// Check shared domain
	if hasSharedDomain {
		fmt.Printf("Found internal/shared/domain/%s.go\n", config.Name)
	}

	if len(domainFiles) == 0 && len(migrationFiles) == 0 && !hasSharedDomain {
		return fmt.Errorf("no domain files, migration files, or shared domain found")
	}

	return nil
}

func AddSharedDomainImport(filePath, domainName, moduleName string) error {
	// Kiểm tra file tồn tại
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File chưa tồn tại, bỏ qua
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return err
	}

	// Kiểm tra xem đã có import shared domain chưa
	importPath := fmt.Sprintf("%s/internal/shared/domain", moduleName)
	hasImport := false

	for _, imp := range node.Imports {
		if strings.Trim(imp.Path.Value, `"`) == importPath {
			hasImport = true
			break
		}
	}

	if !hasImport {
		// Tạo import spec mới
		newImport := &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s"`, importPath),
			},
		}

		// Thêm vào danh sách imports
		if node.Imports == nil {
			node.Imports = []*ast.ImportSpec{newImport}
		} else {
			node.Imports = append(node.Imports, newImport)
		}

		// Ghi lại file
		var buf bytes.Buffer
		if err := format.Node(&buf, fset, node); err != nil {
			return err
		}

		return os.WriteFile(filePath, buf.Bytes(), 0644)
	}

	return nil
}
func AddSharedDomainImportAdvanced(filePath, domainName, moduleName string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return err
	}
	importPath := fmt.Sprintf("%s/shared/domain", moduleName)
	hasImport := false
	for _, imp := range node.Imports {
		if strings.Trim(imp.Path.Value, `"`) == importPath {
			hasImport = true
			break
		}
	}
	if hasImport {
		return nil
	}
	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf(`"%s"`, importPath),
		},
	}
	if len(node.Imports) == 0 {
		node.Imports = []*ast.ImportSpec{newImport}
	} else {

		node.Imports = append(node.Imports, newImport)
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return err
	}

	return os.WriteFile(filePath, buf.Bytes(), 0644)
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
		var newPath string

		if strings.HasPrefix(oldPath, sourceModule) {
			newPath = strings.Replace(oldPath, sourceModule, targetModule, 1)
		}

		if sourceModule != "" && strings.Contains(oldPath, sourceModule) {
			newPath = strings.Replace(oldPath, sourceModule, targetModule, 1)
		}

		if strings.Contains(oldPath, "/internal/") && !strings.HasPrefix(oldPath, targetModule) && !strings.HasPrefix(oldPath, "new_example") {
			parts := strings.Split(oldPath, "/internal/")
			if len(parts) == 2 && !strings.HasPrefix(parts[0], "github.com") && !strings.Contains(parts[0], ".") {

				newPath = targetModule + "/internal/" + parts[1]
			}
		}

		if strings.HasPrefix(oldPath, "internal/") && !strings.HasPrefix(oldPath, targetModule) {
			newPath = targetModule + "/" + oldPath
		}

		if strings.HasPrefix(oldPath, "cmd/") && !strings.HasPrefix(oldPath, targetModule) {
			newPath = targetModule + "/" + oldPath
		}

		if strings.HasPrefix(oldPath, "pkg/") && !strings.HasPrefix(oldPath, targetModule) {
			newPath = targetModule + "/" + oldPath
		}
		if newPath != "" && newPath != oldPath {
			imp.Path.Value = fmt.Sprintf(`"%s"`, newPath)
			changed = true
			fmt.Printf("%s → %s\n", oldPath, newPath)
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

func RewriteAllImportsInDomain(project, domainName, sourceModule, targetModule string) error {
	domainPath := filepath.Join(project, "internal", domainName)

	fmt.Printf("Rewriting imports in domain '%s' from '%s' to '%s'...\n", domainName, sourceModule, targetModule)

	var filesProcessed int
	err := filepath.Walk(domainPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		filesProcessed++
		fmt.Printf("Processing: %s\n", filepath.Base(path))

		if err := RewriteModuleImports(path, sourceModule, targetModule); err != nil {
			return fmt.Errorf("failed to rewrite imports in %s: %w", path, err)
		}
		return nil
	})

	fmt.Printf("Processed %d files\n", filesProcessed)
	return err
}

func FixAllImports(project, sourceModule, targetModule string) error {
	fmt.Printf("Fixing all imports from '%s' to '%s'...\n", sourceModule, targetModule)

	var filesProcessed int
	err := filepath.Walk(project, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.Contains(path, ".git") ||
				strings.Contains(path, "logs") ||
				strings.Contains(path, "swagger") {
				return filepath.SkipDir
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
		if !strings.Contains(string(content), "new_example") &&
			!strings.Contains(string(content), sourceModule) {
			return nil
		}

		filesProcessed++
		fmt.Printf("Fixing: %s\n", filepath.Base(path))

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, content, parser.ParseComments)
		if err != nil {
			return err
		}

		changed := false
		for _, imp := range node.Imports {
			oldPath := strings.Trim(imp.Path.Value, `"`)
			var newPath string

			// Thay thế tất cả các trường hợp
			if strings.Contains(oldPath, "new_example") {
				newPath = strings.Replace(oldPath, "new_example", targetModule, -1)
			}
			if strings.Contains(oldPath, sourceModule) {
				newPath = strings.Replace(oldPath, sourceModule, targetModule, -1)
			}

			if newPath != "" && newPath != oldPath {
				imp.Path.Value = fmt.Sprintf(`"%s"`, newPath)
				changed = true
				fmt.Printf("%s → %s\n", oldPath, newPath)
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

	fmt.Printf("Fixed %d files\n", filesProcessed)
	return err
}

func GetModuleNameFromRepo(repoURL string) string {
	repoURL = strings.TrimPrefix(repoURL, "https://")
	repoURL = strings.TrimPrefix(repoURL, "http://")
	repoURL = strings.TrimSuffix(repoURL, ".git")
	return repoURL
}
func GetCurrentModule(project string) (string, error) {
	goModPath := filepath.Join(project, "go.mod")

	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	return "", fmt.Errorf("module directive not found in go.mod")
}
