package template

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/DVV-15324/witches/pkg/core/templates/utils"
)

//go:embed template/*.tmpl
//go:embed template/cmd/*.tmpl
//go:embed template/cmd/server/config/*.tmpl
//go:embed template/cmd/server/routers/*.tmpl
//go:embed template/internal/auth-service/dto/request/*.tmpl
//go:embed template/internal/auth-service/dto/response/*.tmpl
//go:embed template/internal/auth-service/entity/*.tmpl
//go:embed template/internal/auth-service/handler/*.tmpl
//go:embed template/internal/auth-service/mapping/*.tmpl
//go:embed template/internal/auth-service/repository/*.tmpl
//go:embed template/internal/auth-service/usecase/*.tmpl
//go:embed template/internal/refresh-service/dto/request/*.tmpl
//go:embed template/internal/refresh-service/dto/response/*.tmpl
//go:embed template/internal/refresh-service/entity/*.tmpl
//go:embed template/internal/refresh-service/handler/*.tmpl
//go:embed template/internal/refresh-service/mapping/*.tmpl
//go:embed template/internal/refresh-service/repository/*.tmpl
//go:embed template/internal/refresh-service/usecase/*.tmpl

//go:embed template/internal/user-service/dto/request/*.tmpl
//go:embed template/internal/user-service/dto/response/*.tmpl
//go:embed template/internal/user-service/entity/*.tmpl
//go:embed template/internal/user-service/handler/*.tmpl
//go:embed template/internal/user-service/mapping/*.tmpl
//go:embed template/internal/user-service/repository/*.tmpl
//go:embed template/internal/user-service/usecase/*.tmpl
//go:embed template/internal/shared/middleware/*.tmpl
//go:embed template/internal/shared/model/*.tmpl
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
	fmt.Printf("  cd %s\n", projectName)
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
		migraUp = "template/migrate/mysql/1_create_table.up.sql"
		migraDown = "template/migrate/mysql/1_drop_table.down.sql"
	case "postgres", "postgresql":
		migraUp = "template/migrate/postgresql/1_create_table.up.sql"
		migraDown = "template/migrate/postgresql/1_drop_table.down.sql"
	case "sqlserver", "mssql":
		migraUp = "template/migrate/mssql/1_create_table.up.sql"
		migraDown = "template/migrate/mssql/1_drop_table.down.sql"
	default:
		log.Fatalf("Error: unsupported database: %s. supported : mysql, postgresql, postgres, mssql, sqlserver", typeDb)
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
		"template/internal/auth-service/dto/request/login.go.tmpl":    filepath.Join("internal", "auth-service", "dto", "request", "login.go"),
		"template/internal/auth-service/dto/request/register.go.tmpl": filepath.Join("internal", "auth-service", "dto", "request", "register.go"),
		"template/internal/auth-service/dto/request/errors.go.tmpl":   filepath.Join("internal", "auth-service", "dto", "request", "errors.go"),
		"template/internal/auth-service/dto/request/validate.go.tmpl": filepath.Join("internal", "auth-service", "dto", "request", "validate.go"),
		"template/internal/auth-service/dto/request/gg.go.tmpl":       filepath.Join("internal", "auth-service", "dto", "request", "gg.go"),
		// AUTH SERVICE - DTO RESPONSE
		"template/internal/auth-service/dto/response/auth.go.tmpl": filepath.Join("internal", "auth-service", "dto", "response", "auth.go"),
		// AUTH SERVICE - ENTITY
		"template/internal/auth-service/entity/entity.go.tmpl": filepath.Join("internal", "auth-service", "entity", "entity.go"),
		// AUTH SERVICE - HANDLER
		"template/internal/auth-service/handler/handler.go.tmpl":  filepath.Join("internal", "auth-service", "handler", "handler.go"),
		"template/internal/auth-service/handler/login.go.tmpl":    filepath.Join("internal", "auth-service", "handler", "login.go"),
		"template/internal/auth-service/handler/registry.go.tmpl": filepath.Join("internal", "auth-service", "handler", "registry.go"),
		"template/internal/auth-service/handler/google.go.tmpl":   filepath.Join("internal", "auth-service", "handler", "google.go"),
		"template/internal/auth-service/handler/logout.go.tmpl":   filepath.Join("internal", "auth-service", "handler", "logout.go"),
		// AUTH SERVICE - MAPPING
		"template/internal/auth-service/mapping/mapping.go.tmpl": filepath.Join("internal", "auth-service", "mapping", "mapping.go"),
		// AUTH SERVICE - REPOSITORY
		"template/internal/auth-service/repository/repository.go.tmpl": filepath.Join("internal", "auth-service", "repository", "repository.go"),
		"template/internal/auth-service/repository/create.go.tmpl":     filepath.Join("internal", "auth-service", "repository", "create.go"),
		"template/internal/auth-service/repository/get.go.tmpl":        filepath.Join("internal", "auth-service", "repository", "get.go"),
		"template/internal/auth-service/repository/update.go.tmpl":     filepath.Join("internal", "auth-service", "repository", "update.go"),
		"template/internal/auth-service/repository/cache.go.tmpl":      filepath.Join("internal", "auth-service", "repository", "cache.go"),
		"template/internal/auth-service/repository/key.go.tmpl":        filepath.Join("internal", "auth-service", "repository", "key.go"),
		// AUTH SERVICE - USECASE
		"template/internal/auth-service/usecase/usecase.go.tmpl":  filepath.Join("internal", "auth-service", "usecase", "usecase.go"),
		"template/internal/auth-service/usecase/login.go.tmpl":    filepath.Join("internal", "auth-service", "usecase", "login.go"),
		"template/internal/auth-service/usecase/register.go.tmpl": filepath.Join("internal", "auth-service", "usecase", "register.go"),
		"template/internal/auth-service/usecase/get.go.tmpl":      filepath.Join("internal", "auth-service", "usecase", "get.go"),
		"template/internal/auth-service/usecase/google.go.tmpl":   filepath.Join("internal", "auth-service", "usecase", "google.go"),
		"template/internal/auth-service/usecase/logout.go.tmpl":   filepath.Join("internal", "auth-service", "usecase", "logout.go"),
		"template/internal/auth-service/usecase/errors.go.tmpl":   filepath.Join("internal", "auth-service", "usecase", "errors.go"),

		// REFRESH SERVICE - DTO
		"template/internal/refresh-service/dto/request/request.go.tmpl":   filepath.Join("internal", "refresh-service", "dto", "request", "request.go"),
		"template/internal/refresh-service/dto/response/response.go.tmpl": filepath.Join("internal", "refresh-service", "dto", "response", "response.go"),
		// REFRESH SERVICE - ENTITY
		"template/internal/refresh-service/entity/entity.go.tmpl": filepath.Join("internal", "refresh-service", "entity", "entity.go"),
		// REFRESH SERVICE - HANDLER
		"template/internal/refresh-service/handler/handler.go.tmpl":  filepath.Join("internal", "refresh-service", "handler", "handler.go"),
		"template/internal/refresh-service/handler/re_token.go.tmpl": filepath.Join("internal", "refresh-service", "handler", "re_token.go"),
		// REFRESH SERVICE - MAPPING
		"template/internal/refresh-service/mapping/mapping.go.tmpl": filepath.Join("internal", "refresh-service", "mapping", "mapping.go"),
		// REFRESH SERVICE - REPOSITORY
		"template/internal/refresh-service/repository/repository.go.tmpl": filepath.Join("internal", "refresh-service", "repository", "repository.go"),
		"template/internal/refresh-service/repository/create.go.tmpl":     filepath.Join("internal", "refresh-service", "repository", "create.go"),
		"template/internal/refresh-service/repository/get.go.tmpl":        filepath.Join("internal", "refresh-service", "repository", "get.go"),
		"template/internal/refresh-service/repository/revoke.go.tmpl":     filepath.Join("internal", "refresh-service", "repository", "revoke.go"),
		"template/internal/refresh-service/repository/cache.go.tmpl":      filepath.Join("internal", "refresh-service", "repository", "cache.go"),
		"template/internal/refresh-service/repository/key.go.tmpl":        filepath.Join("internal", "refresh-service", "repository", "key.go"),
		// REFRESH SERVICE - USECASE
		"template/internal/refresh-service/usecase/usecase.go.tmpl": filepath.Join("internal", "refresh-service", "usecase", "usecase.go"),
		"template/internal/refresh-service/usecase/create.go.tmpl":  filepath.Join("internal", "refresh-service", "usecase", "create.go"),
		"template/internal/refresh-service/usecase/get.go.tmpl":     filepath.Join("internal", "refresh-service", "usecase", "get.go"),
		"template/internal/refresh-service/usecase/refresh.go.tmpl": filepath.Join("internal", "refresh-service", "usecase", "refresh.go"),
		"template/internal/refresh-service/usecase/revoke.go.tmpl":  filepath.Join("internal", "refresh-service", "usecase", "revoke.go"),
		"template/internal/refresh-service/usecase/token.go.tmpl":   filepath.Join("internal", "refresh-service", "usecase", "token.go"),

		// USER SERVICE - DTO
		"template/internal/user-service/dto/request/errors.go.tmpl":    filepath.Join("internal", "user-service", "dto", "request", "errors.go"),
		"template/internal/user-service/dto/request/validate.go.tmpl":  filepath.Join("internal", "user-service", "dto", "request", "validate.go"),
		"template/internal/user-service/dto/response/response.go.tmpl": filepath.Join("internal", "user-service", "dto", "response", "response.go"),
		// USER SERVICE - ENTITY
		"template/internal/user-service/entity/entity.go.tmpl": filepath.Join("internal", "user-service", "entity", "entity.go"),
		// USER SERVICE - HANDLER
		"template/internal/user-service/handler/handler.go.tmpl": filepath.Join("internal", "user-service", "handler", "handler.go"),
		"template/internal/user-service/handler/get.go.tmpl":     filepath.Join("internal", "user-service", "handler", "get.go"),
		// USER SERVICE - MAPPING
		"template/internal/user-service/mapping/mapping.go.tmpl": filepath.Join("internal", "user-service", "mapping", "mapping.go"),
		// USER SERVICE - REPOSITORY
		"template/internal/user-service/repository/repository.go.tmpl": filepath.Join("internal", "user-service", "repository", "repository.go"),
		"template/internal/user-service/repository/create.go.tmpl":     filepath.Join("internal", "user-service", "repository", "create.go"),
		"template/internal/user-service/repository/get.go.tmpl":        filepath.Join("internal", "user-service", "repository", "get.go"),
		// USER SERVICE - USECASE
		"template/internal/user-service/usecase/usecase.go.tmpl": filepath.Join("internal", "user-service", "usecase", "usecase.go"),
		"template/internal/user-service/usecase/create.go.tmpl":  filepath.Join("internal", "user-service", "usecase", "create.go"),
		"template/internal/user-service/usecase/get.go.tmpl":     filepath.Join("internal", "user-service", "usecase", "get.go"),

		// SHARED - MIDDLEWARE
		"template/internal/shared/middleware/cors.go.tmpl":       filepath.Join("internal", "shared", "middleware", "cors.go"),
		"template/internal/shared/middleware/limit.go.tmpl":      filepath.Join("internal", "shared", "middleware", "limit.go"),
		"template/internal/shared/middleware/middleware.go.tmpl": filepath.Join("internal", "shared", "middleware", "middleware.go"),
		"template/internal/shared/middleware/timing.go.tmpl":     filepath.Join("internal", "shared", "middleware", "timing.go"),
		// SHARED - MODEL
		"template/internal/shared/model/auth.go.tmpl":    filepath.Join("internal", "shared", "model", "auth.go"),
		"template/internal/shared/model/refresh.go.tmpl": filepath.Join("internal", "shared", "model", "refresh.go"),
		"template/internal/shared/model/user.go.tmpl":    filepath.Join("internal", "shared", "model", "user.go"),
		// SHARED - UTILS
		"template/internal/shared/utils/helper.go.tmpl":     filepath.Join("internal", "shared", "utils", "helper.go"),
		"template/internal/shared/utils/key_object.go.tmpl": filepath.Join("internal", "shared", "utils", "key_object.go"),
		"template/internal/shared/utils/key_req.go.tmpl":    filepath.Join("internal", "shared", "utils", "key_req.go"),
		"template/internal/shared/utils/mapping.go.tmpl":    filepath.Join("internal", "shared", "utils", "mapping.go"),
		"template/internal/shared/utils/uid.go.tmpl":        filepath.Join("internal", "shared", "utils", "uid.go"),

		// LOGS

		// MIGRATIONS
		migraUp:   "migrate/migrations/1_create_table.up.sql.tmpl",
		migraDown: "migrate/migrations/1_drop_table.down.sql.tmpl",

		// PKG REDIS
		"template/pkg/redis/client.go.tmpl": "pkg/redis/client.go",
	}

	for tmpl, dest := range files {
		utils.RenderTemplate(templateFS, baseDir, dest, tmpl, config)
	}

	return nil
}
