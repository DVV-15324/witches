package template

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

//go:embed service/dto/request/*.tmpl
//go:embed service/dto/response/*.tmpl
//go:embed service/entity/*.tmpl
//go:embed service/handler/*.tmpl
//go:embed service/mapping/*.tmpl
//go:embed service/repository/*.tmpl
//go:embed service/usecase/*.tmpl
//go:embed service/shared/model/model.go.tmpl
var templateSvFS embed.FS

type ServiceConfig struct {
	NameCap    string // Book
	Name       string // book
	FolderName string // book-service
	ModuleName string // example_3
}

func AddGoService(projectRoot string, moduleName string, serviceN string) {

	serviceName := strings.ToLower(serviceN)
	serviceNameCap := strings.Title(serviceName)

	config := ServiceConfig{
		NameCap:    serviceNameCap,
		Name:       serviceName,
		FolderName: serviceName + "-service",
		ModuleName: moduleName,
	}

	fmt.Printf("Generating service '%s' ...\n", config.FolderName)

	if err := generateService(projectRoot, config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Service '%s' generated successfully!\n", config.FolderName)
}

func generateService(projectRoot string, config ServiceConfig) error {
	baseDir := filepath.Join(projectRoot, "internal", config.FolderName)

	// Tạo thư mục
	dirs := []string{
		"dto/request",
		"dto/response",
		"entity",
		"handler",
		"mapping",
		"repository",
		"usecase",
	}

	for _, dir := range dirs {
		path := filepath.Join(baseDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", path, err)
		}
	}

	// Map template files -> destination files
	files := map[string]string{
		"service/dto/request/request.go.tmpl":   "dto/request/request.go",
		"service/dto/response/response.go.tmpl": "dto/response/response.go",
		"service/entity/entity.go.tmpl":         "entity/entity.go",
		"service/handler/handler.go.tmpl":       "handler/handler.go",
		"service/handler/create.go.tmpl":        "handler/create.go",
		"service/handler/get.go.tmpl":           "handler/get.go",
		"service/handler/update.go.tmpl":        "handler/update.go",
		"service/handler/delete.go.tmpl":        "handler/delete.go",
		"service/mapping/mapping.go.tmpl":       "mapping/mapping.go",
		"service/repository/repository.go.tmpl": "repository/repository.go",
		"service/repository/create.go.tmpl":     "repository/create.go",
		"service/repository/get.go.tmpl":        "repository/get.go",
		"service/repository/update.go.tmpl":     "repository/update.go",
		"service/repository/delete.go.tmpl":     "repository/delete.go",
		"service/usecase/usecase.go.tmpl":       "usecase/usecase.go",
		"service/usecase/create.go.tmpl":        "usecase/create.go",
		"service/usecase/get.go.tmpl":           "usecase/get.go",
		"service/usecase/update.go.tmpl":        "usecase/update.go",
		"service/usecase/delete.go.tmpl":        "usecase/delete.go",
	}

	for tmpl, dest := range files {
		if err := renderTemplateSv(baseDir, dest, tmpl, config); err != nil {
			return fmt.Errorf("failed to render %s: %v", dest, err)
		}
	}

	// Sinh shared model
	if err := generateSharedModel(projectRoot, config); err != nil {
		fmt.Printf("Warning: failed to generate shared model: %v\n", err)
	}

	// Cập nhật key_object.go
	if err := updateKeyObject(projectRoot, config); err != nil {
		fmt.Printf("Warning: failed to update key_object.go: %v\n", err)
	}

	return nil
}
func generateSharedModel(projectRoot string, config ServiceConfig) error {
	sharedModelDir := filepath.Join(projectRoot, "internal", "shared", "model")
	if err := os.MkdirAll(sharedModelDir, 0755); err != nil {
		return err
	}

	// Tên file: {name}_model.go (VD: book_model.go)
	destFile := filepath.Join(sharedModelDir, config.Name+".go")
	tmplFile := "service/shared/model/model.go.tmpl"

	tmplContent, err := templateSvFS.ReadFile(tmplFile)
	if err != nil {
		return fmt.Errorf("template %s not found", tmplFile)
	}

	tmpl, err := template.New("model.go.tmpl").Parse(string(tmplContent))
	if err != nil {
		return err
	}

	file, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, config)
}

func updateKeyObject(projectRoot string, config ServiceConfig) error {
	keyFile := filepath.Join(projectRoot, "internal", "shared", "utils", "key_object.go")

	// Đọc file hiện tại
	content, err := os.ReadFile(keyFile)
	if err != nil {
		// Nếu file chưa tồn tại, tạo mới
		return createKeyObjectFile(keyFile, config)
	}

	lines := strings.Split(string(content), "\n")

	// Kiểm tra xem constant đã tồn tại chưa
	for _, line := range lines {
		if strings.Contains(line, "Object"+config.NameCap) {
			fmt.Printf("Object%s already exists in key_object.go\n", config.NameCap)
			return nil
		}
	}

	// Tìm max ID hiện tại
	maxID := 0
	re := regexp.MustCompile(`Object\w+\s+uint\s*=\s*(\d+)`)
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			id := 0
			fmt.Sscanf(matches[1], "%d", &id)
			if id > maxID {
				maxID = id
			}
		}
	}

	// Tạo constant mới với ID = maxID + 1
	newLine := fmt.Sprintf("\tObject%s uint = %d", config.NameCap, maxID+1)

	// Tìm vị trí để chèn (sau dòng cuối cùng trong block var)
	insertIndex := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == ")" {
			insertIndex = i
			break
		}
	}

	if insertIndex == -1 {
		// Không tìm thấy dấu đóng ngoặc, thêm vào cuối file
		lines = append(lines, "")
		lines = append(lines, "var (")
		lines = append(lines, newLine)
		lines = append(lines, ")")
	} else {
		// Chèn trước dấu đóng ngoặc
		newLines := make([]string, 0, len(lines)+1)
		newLines = append(newLines, lines[:insertIndex]...)
		newLines = append(newLines, newLine)
		newLines = append(newLines, lines[insertIndex:]...)
		lines = newLines
	}

	// Ghi lại file
	newContent := strings.Join(lines, "\n")
	return os.WriteFile(keyFile, []byte(newContent), 0644)
}

func createKeyObjectFile(keyFile string, config ServiceConfig) error {
	dir := filepath.Dir(keyFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	content := `package utils

var (
	ObjectUser uint = 1
	Object%s   uint = 2
)
`
	return os.WriteFile(keyFile, []byte(fmt.Sprintf(content, config.NameCap)), 0644)
}

func renderTemplateSv(baseDir, destFile, tmplFile string, config ServiceConfig) error {
	tmplContent, err := templateSvFS.ReadFile(tmplFile)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %v", tmplFile, err)
	}

	tmpl, err := template.New(filepath.Base(tmplFile)).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %v", tmplFile, err)
	}

	fullPath := filepath.Join(baseDir, destFile)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %v", fullPath, err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", fullPath, err)
	}
	defer file.Close()

	return tmpl.Execute(file, config)
}
