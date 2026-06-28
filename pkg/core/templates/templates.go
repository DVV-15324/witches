package templates

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*.tmpl
//go:embed templates/cmd/*.tmpl
//go:embed templates/cmd/server/config/*.tmpl
//go:embed templates/cmd/server/routers/*.tmpl
//go:embed templates/internal/dto/auth/request/*.tmpl
//go:embed templates/internal/dto/auth/response/*.tmpl
//go:embed templates/internal/dto/user/request/*.tmpl
//go:embed templates/internal/dto/user/response/*.tmpl
//go:embed templates/internal/entity/auth/*.tmpl
//go:embed templates/internal/entity/user/*.tmpl
//go:embed templates/internal/handler/auth/*.tmpl
//go:embed templates/internal/handler/user/*.tmpl
//go:embed templates/internal/mapping/*.tmpl
//go:embed templates/internal/middleware/*.tmpl
//go:embed templates/internal/repository/auth/*.tmpl
//go:embed templates/internal/repository/user/*.tmpl
//go:embed templates/internal/usecase/auth/*.tmpl
//go:embed templates/internal/usecase/user/*.tmpl
//go:embed templates/internal/utils/*.tmpl
//go:embed templates/logs/*.tmpl
//go:embed templates/migrate/migrations/*.tmpl
var templates embed.FS

type ProjectConfig struct {
	ProjectName string
	ModuleName  string
}

func TemplatesCreate(projectName string) {
	config := ProjectConfig{
		ProjectName: projectName,
		ModuleName:  projectName,
	}

	fmt.Printf("Creating project: %s\n", projectName)
	fmt.Println("Generating structure...")

	if err := createProjectStructure(config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Project created successfully!")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  witches install\n")
	fmt.Printf("  witches run\n")
}

func createProjectStructure(config ProjectConfig) error {
	baseDir, _ := os.Getwd()

	// Tạo tất cả thư mục cần thiết
	dirs := []string{
		"cmd/server/config",
		"cmd/server/routers",
		"internal/dto/auth/request",
		"internal/dto/auth/response",
		"internal/dto/user/request",
		"internal/dto/user/response",
		"internal/entity/auth",
		"internal/entity/user",
		"internal/handler/auth",
		"internal/handler/user",
		"internal/mapping",
		"internal/middleware",
		"internal/repository/auth",
		"internal/repository/user",
		"internal/usecase/auth",
		"internal/usecase/user",
		"internal/utils",
		"logs",
		"migrate/migrations",
	}

	for _, dir := range dirs {
		path := filepath.Join(baseDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", path, err)
		}
	}

	// Map template files -> destination files
	files := map[string]string{
		// ROOT FILES
		"templates/main.go.tmpl":   "main.go",
		"templates/go.mod.tmpl":    "go.mod",
		"templates/README.MD.tmpl": "README.md",

		// CMD
		"templates/cmd/root.go.tmpl":                    "cmd/root.go",
		"templates/cmd/server/config/config.go.tmpl":    "cmd/server/config/config.go",
		"templates/cmd/server/routers/composer.go.tmpl": "cmd/server/routers/composer.go",
		"templates/cmd/server/routers/router.go.tmpl":   "cmd/server/routers/router.go",

		// DTO AUTH REQUEST
		"templates/internal/dto/auth/request/login.go.tmpl":    "internal/dto/auth/request/login.go",
		"templates/internal/dto/auth/request/register.go.tmpl": "internal/dto/auth/request/register.go",
		"templates/internal/dto/auth/request/errors.go.tmpl":   "internal/dto/auth/request/errors.go",
		"templates/internal/dto/auth/request/validate.go.tmpl": "internal/dto/auth/request/validate.go",

		// DTO AUTH RESPONSE
		"templates/internal/dto/auth/response/auth.go.tmpl": "internal/dto/auth/response/auth.go",

		// DTO USER REQUEST
		"templates/internal/dto/user/request/errors.go.tmpl":   "internal/dto/user/request/errors.go",
		"templates/internal/dto/user/request/validate.go.tmpl": "internal/dto/user/request/validate.go",

		// DTO USER RESPONSE
		"templates/internal/dto/user/response/user.go.tmpl": "internal/dto/user/response/user.go",

		// ENTITY
		"templates/internal/entity/auth/auth.go.tmpl": "internal/entity/auth/auth.go",
		"templates/internal/entity/user/user.go.tmpl": "internal/entity/user/user.go",

		// HANDLER AUTH
		"templates/internal/handler/auth/auth.go.tmpl":     "internal/handler/auth/auth.go",
		"templates/internal/handler/auth/login.go.tmpl":    "internal/handler/auth/login.go",
		"templates/internal/handler/auth/registry.go.tmpl": "internal/handler/auth/registry.go",

		// HANDLER USER
		"templates/internal/handler/user/get.go.tmpl":  "internal/handler/user/get.go",
		"templates/internal/handler/user/user.go.tmpl": "internal/handler/user/user.go",

		// MAPPING
		"templates/internal/mapping/auth.go.tmpl": "internal/mapping/auth.go",
		"templates/internal/mapping/key.go.tmpl":  "internal/mapping/key.go",
		"templates/internal/mapping/user.go.tmpl": "internal/mapping/user.go",

		// MIDDLEWARE
		"templates/internal/middleware/cors.go.tmpl":       "internal/middleware/cors.go",
		"templates/internal/middleware/middleware.go.tmpl": "internal/middleware/middleware.go",

		// REPOSITORY AUTH
		"templates/internal/repository/auth/auth.go.tmpl":   "internal/repository/auth/auth.go",
		"templates/internal/repository/auth/create.go.tmpl": "internal/repository/auth/create.go",
		"templates/internal/repository/auth/get.go.tmpl":    "internal/repository/auth/get.go",

		// REPOSITORY USER
		"templates/internal/repository/user/create.go.tmpl": "internal/repository/user/create.go",
		"templates/internal/repository/user/get.go.tmpl":    "internal/repository/user/get.go",
		"templates/internal/repository/user/user.go.tmpl":   "internal/repository/user/user.go",

		// USECASE AUTH
		"templates/internal/usecase/auth/auth.go.tmpl":     "internal/usecase/auth/auth.go",
		"templates/internal/usecase/auth/errors.go.tmpl":   "internal/usecase/auth/errors.go",
		"templates/internal/usecase/auth/login.go.tmpl":    "internal/usecase/auth/login.go",
		"templates/internal/usecase/auth/register.go.tmpl": "internal/usecase/auth/register.go",
		"templates/internal/usecase/auth/token.go.tmpl":    "internal/usecase/auth/token.go",

		// USECASE USER
		"templates/internal/usecase/user/create.go.tmpl": "internal/usecase/user/create.go",
		"templates/internal/usecase/user/get.go.tmpl":    "internal/usecase/user/get.go",
		"templates/internal/usecase/user/user.go.tmpl":   "internal/usecase/user/user.go",

		// UTILS
		"templates/internal/utils/connect.go.tmpl": "internal/utils/connect.go",

		// LOGS
		"templates/logs/logs.log.tmpl": "logs/logs.log",

		// MIGRATIONS
		"templates/migrate/migrations/1_create_table.up.sql.tmpl": "migrate/migrations/1_create_table.up.sql",
		"templates/migrate/migrations/1_drop_table.down.sql.tmpl": "migrate/migrations/1_drop_table.down.sql",
	}

	for tmpl, dest := range files {
		if err := renderTemplate(baseDir, dest, tmpl, config); err != nil {
			return fmt.Errorf("failed to render %s: %v", dest, err)
		}
	}

	return nil
}

func renderTemplate(baseDir, destFile, tmplFile string, config ProjectConfig) error {
	// Đọc template từ embed
	tmplContent, err := templates.ReadFile(tmplFile)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %v", tmplFile, err)
	}

	tmpl, err := template.New(filepath.Base(tmplFile)).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %v", tmplFile, err)
	}

	// Tạo file đích
	fullPath := filepath.Join(baseDir, destFile)

	// Tạo thư mục cha nếu chưa tồn tại
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %v", fullPath, err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", fullPath, err)
	}
	defer file.Close()

	// Render template
	return tmpl.Execute(file, config)
}
