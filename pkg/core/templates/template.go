package template

import (
	"embed"
	"fmt"
	"github.com/DVV-15324/witches/pkg/core/templates/utils"
	"os"
	"path/filepath"
)

//go:embed template/*.tmpl
//go:embed template/cmd/*.tmpl
//go:embed template/cmd/server/config/*.tmpl
//go:embed template/cmd/server/routers/*.tmpl
//go:embed template/internal/auth/dto/request/*.tmpl
//go:embed template/internal/auth/dto/response/*.tmpl
//go:embed template/internal/auth/model/*.tmpl
//go:embed template/internal/auth/handler/*.tmpl
//go:embed template/internal/auth/mapping/*.tmpl
//go:embed template/internal/auth/repository/*.tmpl
//go:embed template/internal/auth/usecase/*.tmpl
//go:embed template/internal/refresh/dto/request/*.tmpl
//go:embed template/internal/refresh/dto/response/*.tmpl
//go:embed template/internal/refresh/model/*.tmpl
//go:embed template/internal/refresh/handler/*.tmpl
//go:embed template/internal/refresh/mapping/*.tmpl
//go:embed template/internal/refresh/repository/*.tmpl
//go:embed template/internal/refresh/usecase/*.tmpl

//go:embed template/internal/user/dto/request/*.tmpl
//go:embed template/internal/user/dto/response/*.tmpl
//go:embed template/internal/user/model/*.tmpl
//go:embed template/internal/user/handler/*.tmpl
//go:embed template/internal/user/mapping/*.tmpl
//go:embed template/internal/user/repository/*.tmpl
//go:embed template/internal/user/usecase/*.tmpl
//go:embed template/internal/shared/middleware/*.tmpl
//go:embed template/internal/shared/domain/*.tmpl
//go:embed template/internal/shared/utils/*.tmpl
//go:embed template/migrate/mssql/*.tmpl
//go:embed template/migrate/mysql/*.tmpl
//go:embed template/migrate/postgresql/*.tmpl
//go:embed template/pkg/redis/*.tmpl
var templateFS embed.FS

type ProjectConfig struct {
	ModuleName string
}

func (p ProjectConfig) GetMuduleName() string {
	return p.ModuleName
}

func CreateGoArcRefresh(projectName string, typeDb string) {

	config := ProjectConfig{
		ModuleName: projectName,
	}

	fmt.Printf("Generating project: %s\n", projectName)
	fmt.Println("Creating structure...")

	if err := createProjectStructure(config, typeDb); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Project created successfully!")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  witches install\n")
	fmt.Printf("  witches run\n")
}

func createProjectStructure(config ProjectConfig, typeDb string) error {
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
		return fmt.Errorf("Error: unsupported database: %s. supported : mysql, postgresql, postgres, mssql, sqlserver", typeDb)
	}

	// Map template files -> destination files
	files := map[string]string{
		// ROOT
		"template/main.go.tmpl":   "main.go",
		"template/README.md.tmpl": "README.md",
		"template/go.mod.tmpl":    "go.mod",

		// CMD
		"template/cmd/root.go.tmpl":                     filepath.Join("cmd", "root.go"),
		"template/cmd/server/config/config.go.tmpl":     filepath.Join("cmd", "server", "config", "config.go"),
		"template/cmd/server/routers/composer.go.tmpl":  filepath.Join("cmd", "server", "routers", "composer.go"),
		"template/cmd/server/routers/protected.go.tmpl": filepath.Join("cmd", "server", "routers", "protected.go"),
		"template/cmd/server/routers/public.go.tmpl":    filepath.Join("cmd", "server", "routers", "public.go"),
		"template/cmd/server/routers/routers.go.tmpl":   filepath.Join("cmd", "server", "routers", "routers.go"),

		// AUTH SERVICE - DTO REQUEST
		"template/internal/auth/dto/request/login.go.tmpl":    filepath.Join("internal", "auth", "dto", "request", "login.go"),
		"template/internal/auth/dto/request/register.go.tmpl": filepath.Join("internal", "auth", "dto", "request", "register.go"),
		"template/internal/auth/dto/request/errors.go.tmpl":   filepath.Join("internal", "auth", "dto", "request", "errors.go"),
		"template/internal/auth/dto/request/validate.go.tmpl": filepath.Join("internal", "auth", "dto", "request", "validate.go"),
		"template/internal/auth/dto/request/gg.go.tmpl":       filepath.Join("internal", "auth", "dto", "request", "gg.go"),
		// AUTH SERVICE - DTO RESPONSE
		"template/internal/auth/dto/response/auth.go.tmpl": filepath.Join("internal", "auth", "dto", "response", "auth.go"),
		// AUTH SERVICE - model
		"template/internal/auth/model/model.go.tmpl": filepath.Join("internal", "auth", "model", "model.go"),
		// AUTH SERVICE - HANDLER
		"template/internal/auth/handler/handler.go.tmpl":  filepath.Join("internal", "auth", "handler", "handler.go"),
		"template/internal/auth/handler/login.go.tmpl":    filepath.Join("internal", "auth", "handler", "login.go"),
		"template/internal/auth/handler/registry.go.tmpl": filepath.Join("internal", "auth", "handler", "registry.go"),
		"template/internal/auth/handler/google.go.tmpl":   filepath.Join("internal", "auth", "handler", "google.go"),
		"template/internal/auth/handler/logout.go.tmpl":   filepath.Join("internal", "auth", "handler", "logout.go"),
		// AUTH SERVICE - MAPPING
		"template/internal/auth/mapping/mapping.go.tmpl": filepath.Join("internal", "auth", "mapping", "mapping.go"),
		// AUTH SERVICE - REPOSITORY
		"template/internal/auth/repository/repository.go.tmpl": filepath.Join("internal", "auth", "repository", "repository.go"),
		"template/internal/auth/repository/create.go.tmpl":     filepath.Join("internal", "auth", "repository", "create.go"),
		"template/internal/auth/repository/get.go.tmpl":        filepath.Join("internal", "auth", "repository", "get.go"),
		"template/internal/auth/repository/update.go.tmpl":     filepath.Join("internal", "auth", "repository", "update.go"),
		"template/internal/auth/repository/cache.go.tmpl":      filepath.Join("internal", "auth", "repository", "cache.go"),
		"template/internal/auth/repository/key.go.tmpl":        filepath.Join("internal", "auth", "repository", "key.go"),
		// AUTH SERVICE - USECASE
		"template/internal/auth/usecase/usecase.go.tmpl":  filepath.Join("internal", "auth", "usecase", "usecase.go"),
		"template/internal/auth/usecase/login.go.tmpl":    filepath.Join("internal", "auth", "usecase", "login.go"),
		"template/internal/auth/usecase/register.go.tmpl": filepath.Join("internal", "auth", "usecase", "register.go"),
		"template/internal/auth/usecase/get.go.tmpl":      filepath.Join("internal", "auth", "usecase", "get.go"),
		"template/internal/auth/usecase/google.go.tmpl":   filepath.Join("internal", "auth", "usecase", "google.go"),
		"template/internal/auth/usecase/logout.go.tmpl":   filepath.Join("internal", "auth", "usecase", "logout.go"),
		"template/internal/auth/usecase/errors.go.tmpl":   filepath.Join("internal", "auth", "usecase", "errors.go"),

		// REFRESH SERVICE - DTO
		"template/internal/refresh/dto/request/request.go.tmpl":   filepath.Join("internal", "refresh", "dto", "request", "request.go"),
		"template/internal/refresh/dto/response/response.go.tmpl": filepath.Join("internal", "refresh", "dto", "response", "response.go"),
		// REFRESH SERVICE - model
		"template/internal/refresh/model/model.go.tmpl": filepath.Join("internal", "refresh", "model", "model.go"),
		// REFRESH SERVICE - HANDLER
		"template/internal/refresh/handler/handler.go.tmpl":  filepath.Join("internal", "refresh", "handler", "handler.go"),
		"template/internal/refresh/handler/re_token.go.tmpl": filepath.Join("internal", "refresh", "handler", "re_token.go"),
		// REFRESH SERVICE - MAPPING
		"template/internal/refresh/mapping/mapping.go.tmpl": filepath.Join("internal", "refresh", "mapping", "mapping.go"),
		// REFRESH SERVICE - REPOSITORY
		"template/internal/refresh/repository/repository.go.tmpl": filepath.Join("internal", "refresh", "repository", "repository.go"),
		"template/internal/refresh/repository/create.go.tmpl":     filepath.Join("internal", "refresh", "repository", "create.go"),
		"template/internal/refresh/repository/get.go.tmpl":        filepath.Join("internal", "refresh", "repository", "get.go"),
		"template/internal/refresh/repository/revoke.go.tmpl":     filepath.Join("internal", "refresh", "repository", "revoke.go"),
		"template/internal/refresh/repository/cache.go.tmpl":      filepath.Join("internal", "refresh", "repository", "cache.go"),
		"template/internal/refresh/repository/key.go.tmpl":        filepath.Join("internal", "refresh", "repository", "key.go"),
		// REFRESH SERVICE - USECASE
		"template/internal/refresh/usecase/usecase.go.tmpl": filepath.Join("internal", "refresh", "usecase", "usecase.go"),
		"template/internal/refresh/usecase/create.go.tmpl":  filepath.Join("internal", "refresh", "usecase", "create.go"),
		"template/internal/refresh/usecase/get.go.tmpl":     filepath.Join("internal", "refresh", "usecase", "get.go"),
		"template/internal/refresh/usecase/refresh.go.tmpl": filepath.Join("internal", "refresh", "usecase", "refresh.go"),
		"template/internal/refresh/usecase/revoke.go.tmpl":  filepath.Join("internal", "refresh", "usecase", "revoke.go"),
		"template/internal/refresh/usecase/token.go.tmpl":   filepath.Join("internal", "refresh", "usecase", "token.go"),

		// USER SERVICE - DTO
		"template/internal/user/dto/request/errors.go.tmpl":    filepath.Join("internal", "user", "dto", "request", "errors.go"),
		"template/internal/user/dto/request/validate.go.tmpl":  filepath.Join("internal", "user", "dto", "request", "validate.go"),
		"template/internal/user/dto/response/response.go.tmpl": filepath.Join("internal", "user", "dto", "response", "response.go"),
		// USER SERVICE - model
		"template/internal/user/model/model.go.tmpl": filepath.Join("internal", "user", "model", "model.go"),
		// USER SERVICE - HANDLER
		"template/internal/user/handler/handler.go.tmpl": filepath.Join("internal", "user", "handler", "handler.go"),
		"template/internal/user/handler/get.go.tmpl":     filepath.Join("internal", "user", "handler", "get.go"),
		// USER SERVICE - MAPPING
		"template/internal/user/mapping/mapping.go.tmpl": filepath.Join("internal", "user", "mapping", "mapping.go"),
		// USER SERVICE - REPOSITORY
		"template/internal/user/repository/repository.go.tmpl": filepath.Join("internal", "user", "repository", "repository.go"),
		"template/internal/user/repository/create.go.tmpl":     filepath.Join("internal", "user", "repository", "create.go"),
		"template/internal/user/repository/get.go.tmpl":        filepath.Join("internal", "user", "repository", "get.go"),
		// USER SERVICE - USECASE
		"template/internal/user/usecase/usecase.go.tmpl": filepath.Join("internal", "user", "usecase", "usecase.go"),
		"template/internal/user/usecase/create.go.tmpl":  filepath.Join("internal", "user", "usecase", "create.go"),
		"template/internal/user/usecase/get.go.tmpl":     filepath.Join("internal", "user", "usecase", "get.go"),

		// SHARED - MIDDLEWARE
		"template/internal/shared/middleware/limit.go.tmpl":      filepath.Join("internal", "shared", "middleware", "limit.go"),
		"template/internal/shared/middleware/middleware.go.tmpl": filepath.Join("internal", "shared", "middleware", "middleware.go"),
		"template/internal/shared/middleware/timing.go.tmpl":     filepath.Join("internal", "shared", "middleware", "timing.go"),
		// SHARED - DOMAIN
		"template/internal/shared/domain/auth.go.tmpl":    filepath.Join("internal", "shared", "domain", "auth.go"),
		"template/internal/shared/domain/refresh.go.tmpl": filepath.Join("internal", "shared", "domain", "refresh.go"),
		"template/internal/shared/domain/user.go.tmpl":    filepath.Join("internal", "shared", "domain", "user.go"),
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

	for tmpl, dest := range files {
		utils.RenderTemplate(templateFS, baseDir, dest, tmpl, config)
	}

	return nil
}
