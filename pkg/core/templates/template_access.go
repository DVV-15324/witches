package template

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed go_arc_access/*.tmpl
//go:embed go_arc_access/cmd/*.tmpl
//go:embed go_arc_access/cmd/server/config/*.tmpl
//go:embed go_arc_access/cmd/server/routers/*.tmpl
//go:embed go_arc_access/internal/dto/auth/request/*.tmpl
//go:embed go_arc_access/internal/dto/auth/response/*.tmpl
//go:embed go_arc_access/internal/dto/user/request/*.tmpl
//go:embed go_arc_access/internal/dto/user/response/*.tmpl
//go:embed go_arc_access/internal/entity/auth/*.tmpl
//go:embed go_arc_access/internal/entity/user/*.tmpl
//go:embed go_arc_access/internal/handler/auth/*.tmpl
//go:embed go_arc_access/internal/handler/user/*.tmpl
//go:embed go_arc_access/internal/mapping/*.tmpl
//go:embed go_arc_access/internal/middleware/*.tmpl
//go:embed go_arc_access/internal/repository/auth/*.tmpl
//go:embed go_arc_access/internal/repository/user/*.tmpl
//go:embed go_arc_access/internal/usecase/auth/*.tmpl
//go:embed go_arc_access/internal/usecase/user/*.tmpl
//go:embed go_arc_access/internal/utils/*.tmpl
//go:embed go_arc_access/logs/*.tmpl
//go:embed go_arc_access/migrate/migrations/*.tmpl
var go_arc_access embed.FS

type ProjectConfigAccess struct {
	ProjectName string
	ModuleName  string
}

func CreateGoArcAccess(projectName string) {
	config := ProjectConfigAccess{
		ProjectName: projectName,
		ModuleName:  projectName,
	}

	fmt.Printf("Creating project: %s\n", projectName)
	fmt.Println("Generating structure...")

	if err := createProjectStructureAccess(config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Project created successfully!")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  witches install\n")
	fmt.Printf("  witches run\n")
}

func createProjectStructureAccess(config ProjectConfigAccess) error {
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
		"go_arc_access/main.go.tmpl":   "main.go",
		"go_arc_access/go.mod.tmpl":    "go.mod",
		"go_arc_access/README.MD.tmpl": "README.md",

		// CMD
		"go_arc_access/cmd/root.go.tmpl":                    "cmd/root.go",
		"go_arc_access/cmd/server/config/config.go.tmpl":    "cmd/server/config/config.go",
		"go_arc_access/cmd/server/routers/composer.go.tmpl": "cmd/server/routers/composer.go",
		"go_arc_access/cmd/server/routers/router.go.tmpl":   "cmd/server/routers/router.go",

		// DTO AUTH REQUEST
		"go_arc_access/internal/dto/auth/request/login.go.tmpl":    "internal/dto/auth/request/login.go",
		"go_arc_access/internal/dto/auth/request/register.go.tmpl": "internal/dto/auth/request/register.go",
		"go_arc_access/internal/dto/auth/request/errors.go.tmpl":   "internal/dto/auth/request/errors.go",
		"go_arc_access/internal/dto/auth/request/validate.go.tmpl": "internal/dto/auth/request/validate.go",

		// DTO AUTH RESPONSE
		"go_arc_access/internal/dto/auth/response/auth.go.tmpl": "internal/dto/auth/response/auth.go",

		// DTO USER REQUEST
		"go_arc_access/internal/dto/user/request/errors.go.tmpl":   "internal/dto/user/request/errors.go",
		"go_arc_access/internal/dto/user/request/validate.go.tmpl": "internal/dto/user/request/validate.go",

		// DTO USER RESPONSE
		"go_arc_access/internal/dto/user/response/user.go.tmpl": "internal/dto/user/response/user.go",

		// ENTITY
		"go_arc_access/internal/entity/auth/auth.go.tmpl": "internal/entity/auth/auth.go",
		"go_arc_access/internal/entity/user/user.go.tmpl": "internal/entity/user/user.go",

		// HANDLER AUTH
		"go_arc_access/internal/handler/auth/auth.go.tmpl":     "internal/handler/auth/auth.go",
		"go_arc_access/internal/handler/auth/login.go.tmpl":    "internal/handler/auth/login.go",
		"go_arc_access/internal/handler/auth/registry.go.tmpl": "internal/handler/auth/registry.go",

		// HANDLER USER
		"go_arc_access/internal/handler/user/get.go.tmpl":  "internal/handler/user/get.go",
		"go_arc_access/internal/handler/user/user.go.tmpl": "internal/handler/user/user.go",

		// MAPPING
		"go_arc_access/internal/mapping/auth.go.tmpl": "internal/mapping/auth.go",
		"go_arc_access/internal/mapping/key.go.tmpl":  "internal/mapping/key.go",
		"go_arc_access/internal/mapping/user.go.tmpl": "internal/mapping/user.go",

		// MIDDLEWARE
		"go_arc_access/internal/middleware/cors.go.tmpl":       "internal/middleware/cors.go",
		"go_arc_access/internal/middleware/middleware.go.tmpl": "internal/middleware/middleware.go",

		// REPOSITORY AUTH
		"go_arc_access/internal/repository/auth/auth.go.tmpl":   "internal/repository/auth/auth.go",
		"go_arc_access/internal/repository/auth/create.go.tmpl": "internal/repository/auth/create.go",
		"go_arc_access/internal/repository/auth/get.go.tmpl":    "internal/repository/auth/get.go",

		// REPOSITORY USER
		"go_arc_access/internal/repository/user/create.go.tmpl": "internal/repository/user/create.go",
		"go_arc_access/internal/repository/user/get.go.tmpl":    "internal/repository/user/get.go",
		"go_arc_access/internal/repository/user/user.go.tmpl":   "internal/repository/user/user.go",

		// USECASE AUTH
		"go_arc_access/internal/usecase/auth/auth.go.tmpl":     "internal/usecase/auth/auth.go",
		"go_arc_access/internal/usecase/auth/errors.go.tmpl":   "internal/usecase/auth/errors.go",
		"go_arc_access/internal/usecase/auth/login.go.tmpl":    "internal/usecase/auth/login.go",
		"go_arc_access/internal/usecase/auth/register.go.tmpl": "internal/usecase/auth/register.go",
		"go_arc_access/internal/usecase/auth/token.go.tmpl":    "internal/usecase/auth/token.go",

		// USECASE USER
		"go_arc_access/internal/usecase/user/create.go.tmpl": "internal/usecase/user/create.go",
		"go_arc_access/internal/usecase/user/get.go.tmpl":    "internal/usecase/user/get.go",
		"go_arc_access/internal/usecase/user/user.go.tmpl":   "internal/usecase/user/user.go",

		// UTILS
		"go_arc_access/internal/utils/connect.go.tmpl": "internal/utils/connect.go",

		// LOGS
		"go_arc_access/logs/logs.log.tmpl": "logs/logs.log",

		// MIGRATIONS
		"go_arc_access/migrate/migrations/1_create_table.up.sql.tmpl": "migrate/migrations/1_create_table.up.sql",
		"go_arc_access/migrate/migrations/1_drop_table.down.sql.tmpl": "migrate/migrations/1_drop_table.down.sql",
	}

	for tmpl, dest := range files {
		if err := renderTemplateAccess(baseDir, dest, tmpl, config); err != nil {
			return fmt.Errorf("failed to render %s: %v", dest, err)
		}
	}

	return nil
}

func renderTemplateAccess(baseDir, destFile, tmplFile string, config ProjectConfigAccess) error {
	// Đọc template từ embed
	tmplContent, err := go_arc_access.ReadFile(tmplFile)
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
