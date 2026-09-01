package template

import (
	"embed"
	"fmt"
	"github.com/DVV-15324/witches/pkg/core/templates/utils"
	"os"
	"path/filepath"
)

//go:embed template/pkg/redis/*.tmpl
//go:embed template/*.tmpl
//go:embed template/cmd/*.tmpl
//go:embed template/cmd/server/config/*.tmpl
//go:embed template/cmd/server/routers/*.tmpl
//go:embed template/cmd/server/core/core.go.tmpl
//go:embed template/migrate/mssql/*.tmpl
//go:embed template/migrate/mysql/*.tmpl
//go:embed template/migrate/postgresql/*.tmpl
//go:embed template/internal/shared/utils/*.tmpl
//go:embed template/internal/shared/middleware/*.tmpl
var templateFS embed.FS

type TemplateConfig struct {
	ProjectName string
}

func (p TemplateConfig) GetProjectName() string {
	return p.ProjectName
}

func CreateTemplateGoArc(projectName string, typeDb string) error {
	config := TemplateConfig{
		ProjectName: projectName,
	}
	fmt.Printf("Generating project: %s\n", projectName)
	fmt.Println("Creating structure...")

	if err := createTemplateProjectStructure(config, typeDb); err != nil {
		return fmt.Errorf("create template structure: %w", err)
	}

	fmt.Println("Project created successfully!")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  witches install\n")
	fmt.Printf("  witches run\n")
	return nil
}

func createTemplateProjectStructure(config TemplateConfig, typeDb string) error {
	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	baseDir = filepath.Join(baseDir)

	// Tạo thư mục project
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	var migraUp string
	var migraDown string
	switch typeDb {
	case "mysql":
		migraUp = "template/migrate/mysql/1_create_table.up.sql.tmpl"
		migraDown = "template/migrate/mysql/1_drop_table.down.sql.tmpl"
	case "postgres", "postgresql":
		migraUp = "template/migrate/postgresql/1_create_table.up.sql.tmpl"
		migraDown = "template/migrate/postgresql/1_drop_table.down.sql.tmpl"
	case "sqlserver", "mssql":
		migraUp = "template/migrate/mssql/1_create_table.up.sql.tmpl"
		migraDown = "template/migrate/mssql/1_drop_table.down.sql.tmpl"
	default:
		return fmt.Errorf("error: unsupported database: %s. supported : mysql, postgresql, postgres, mssql, sqlserver", typeDb)
	}

	files := map[string]string{
		// ROOT
		"template/main.go.tmpl":   "main.go",
		"template/README.md.tmpl": "README.md",
		"template/go.mod.tmpl":    "go.mod",

		// CMD
		"template/cmd/root.go.tmpl":                    filepath.Join("cmd", "root.go"),
		"template/cmd/server/config/config.go.tmpl":    filepath.Join("cmd", "server", "config", "config.go"),
		"template/cmd/server/routers/composer.go.tmpl": filepath.Join("cmd", "server", "routers", "composer.go"),
		"template/cmd/server/routers/routers.go.tmpl":  filepath.Join("cmd", "server", "routers", "routers.go"),
		"template/cmd/server/routers/modules.go.tmpl":  filepath.Join("cmd", "server", "routers", "modules.go"),
		"template/cmd/server/core/core.go.tmpl":        filepath.Join("cmd", "server", "core", "core.go"),

		// SHARED - MIDDLEWARE
		"template/internal/shared/middleware/limit.go.tmpl":  filepath.Join("internal", "shared", "middleware", "limit.go"),
		"template/internal/shared/middleware/timing.go.tmpl": filepath.Join("internal", "shared", "middleware", "timing.go"),

		// SHARED - UTILS
		"template/internal/shared/utils/dummy.go.tmpl":      filepath.Join("internal", "shared", "utils", "dummy.go"),
		"template/internal/shared/utils/key_object.go.tmpl": filepath.Join("internal", "shared", "utils", "key_object.go"),
		"template/internal/shared/utils/uid.go.tmpl":        filepath.Join("internal", "shared", "utils", "uid.go"),
		// LOGS
		// MIGRATIONS
		migraUp:   "migrate/migrations/1_create_table.up.sql",
		migraDown: "migrate/migrations/1_drop_table.down.sql",

		// PKG REDIS
		"template/pkg/redis/client.go.tmpl": "pkg/redis/client.go",
	}

	// Tạo thư mục migrate/migrations trước
	migrateDir := filepath.Join(baseDir, "migrate", "migrations")
	if err := os.MkdirAll(migrateDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrate directory: %w", err)
	}

	directories := []string{
		filepath.Join("internal", "shared", "domain"),
	}
	for _, dir := range directories {
		if err := os.MkdirAll(filepath.Join(baseDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	for tmpl, dest := range files {
		// Debug: In ra template đang render
		fmt.Printf("  Rendering: %s -> %s\n", tmpl, dest)
		utils.RenderTemplate(templateFS, baseDir, dest, tmpl, config)
	}

	return nil
}
