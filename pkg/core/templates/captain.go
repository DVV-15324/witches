package template

import (
	"embed"
	"fmt"
	"github.com/DVV-15324/witches/pkg/core/templates/utils"
	"os"
	"path/filepath"
)

//go:embed captain/pkg/redis/*.tmpl
//go:embed captain/*.tmpl
//go:embed captain/cmd/*.tmpl
//go:embed captain/cmd/server/config/*.tmpl
//go:embed captain/cmd/server/routers/*.tmpl
//go:embed captain/cmd/server/core/core.go.tmpl
//go:embed captain/migrate/mssql/*.tmpl
//go:embed captain/migrate/mysql/*.tmpl
//go:embed captain/migrate/postgresql/*.tmpl
//go:embed captain/internal/shared/utils/*.tmpl
//go:embed captain/internal/shared/middleware/*.tmpl
var templateCaptainFS embed.FS

type CaptainConfig struct {
	ModuleName string
}

func (p CaptainConfig) GetModuleName() string {
	return p.ModuleName
}

func CreateCaptainGoArc(projectName string, typeDb string) {
	config := CaptainConfig{
		ModuleName: projectName,
	}
	fmt.Printf("Generating project: %s\n", projectName)
	fmt.Println("Creating structure...")

	if err := createCaptainProjectStructure(config, typeDb); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Project created successfully!")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  witches install\n")
	fmt.Printf("  witches run\n")
}

func createCaptainProjectStructure(config CaptainConfig, typeDb string) error {
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
		migraUp = "captain/migrate/mysql/1_create_table.up.sql.tmpl"
		migraDown = "captain/migrate/mysql/1_drop_table.down.sql.tmpl"
	case "postgres", "postgresql":
		migraUp = "captain/migrate/postgresql/1_create_table.up.sql.tmpl"
		migraDown = "captain/migrate/postgresql/1_drop_table.down.sql.tmpl"
	case "sqlserver", "mssql":
		migraUp = "captain/migrate/mssql/1_create_table.up.sql.tmpl"
		migraDown = "captain/migrate/mssql/1_drop_table.down.sql.tmpl"
	default:
		return fmt.Errorf("error: unsupported database: %s. supported : mysql, postgresql, postgres, mssql, sqlserver", typeDb)
	}

	files := map[string]string{
		// ROOT
		"captain/main.go.tmpl":   "main.go",
		"captain/README.md.tmpl": "README.md",
		"captain/go.mod.tmpl":    "go.mod",

		// CMD
		"captain/cmd/root.go.tmpl":                    filepath.Join("cmd", "root.go"),
		"captain/cmd/server/config/config.go.tmpl":    filepath.Join("cmd", "server", "config", "config.go"),
		"captain/cmd/server/routers/composer.go.tmpl": filepath.Join("cmd", "server", "routers", "composer.go"),
		"captain/cmd/server/routers/routers.go.tmpl":  filepath.Join("cmd", "server", "routers", "routers.go"),
		"captain/cmd/server/routers/modules.go.tmpl":  filepath.Join("cmd", "server", "routers", "modules.go"),
		"captain/cmd/server/core/core.go.tmpl":        filepath.Join("cmd", "server", "core", "core.go"),

		// SHARED - MIDDLEWARE
		"captain/internal/shared/middleware/limit.go.tmpl":  filepath.Join("internal", "shared", "middleware", "limit.go"),
		"captain/internal/shared/middleware/timing.go.tmpl": filepath.Join("internal", "shared", "middleware", "timing.go"),

		// SHARED - UTILS
		"captain/internal/shared/utils/dummy.go.tmpl":      filepath.Join("internal", "shared", "utils", "dummy.go"),
		"captain/internal/shared/utils/key_object.go.tmpl": filepath.Join("internal", "shared", "utils", "key_object.go"),
		"captain/internal/shared/utils/uid.go.tmpl":        filepath.Join("internal", "shared", "utils", "uid.go"),
		// LOGS
		// MIGRATIONS
		migraUp:   "migrate/migrations/1_create_table.up.sql",
		migraDown: "migrate/migrations/1_drop_table.down.sql",

		// PKG REDIS
		"captain/pkg/redis/client.go.tmpl": "pkg/redis/client.go",
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
		utils.RenderTemplate(templateCaptainFS, baseDir, dest, tmpl, config)
	}

	return nil
}
