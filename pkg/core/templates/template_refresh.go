package template

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed go_arc_refresh/*.tmpl
//go:embed go_arc_refresh/cmd/*.tmpl
//go:embed go_arc_refresh/cmd/server/config/*.tmpl
//go:embed go_arc_refresh/cmd/server/routers/*.tmpl
//go:embed go_arc_refresh/internal/dto/auth/request/*.tmpl
//go:embed go_arc_refresh/internal/dto/auth/response/*.tmpl
//go:embed go_arc_refresh/internal/dto/refresh/request/*.tmpl
//go:embed go_arc_refresh/internal/dto/refresh/response/*.tmpl
//go:embed go_arc_refresh/internal/dto/user/request/*.tmpl
//go:embed go_arc_refresh/internal/dto/user/response/*.tmpl
//go:embed go_arc_refresh/internal/entity/auth/*.tmpl
//go:embed go_arc_refresh/internal/entity/refresh/*.tmpl
//go:embed go_arc_refresh/internal/entity/user/*.tmpl
//go:embed go_arc_refresh/internal/handler/auth/*.tmpl
//go:embed go_arc_refresh/internal/handler/refresh/*.tmpl
//go:embed go_arc_refresh/internal/handler/user/*.tmpl
//go:embed go_arc_refresh/internal/mapping/*.tmpl
//go:embed go_arc_refresh/internal/middleware/*.tmpl
//go:embed go_arc_refresh/internal/repository/auth/*.tmpl
//go:embed go_arc_refresh/internal/repository/refresh/*.tmpl
//go:embed go_arc_refresh/internal/repository/user/*.tmpl
//go:embed go_arc_refresh/internal/usecase/auth/*.tmpl
//go:embed go_arc_refresh/internal/usecase/refresh/*.tmpl
//go:embed go_arc_refresh/internal/usecase/user/*.tmpl
//go:embed go_arc_refresh/internal/utils/*.tmpl
//go:embed go_arc_refresh/internal/utils/ratelimit/*.tmpl
//go:embed go_arc_refresh/logs/*.tmpl
//go:embed go_arc_refresh/migrate/migrations/*.tmpl
//go:embed go_arc_refresh/pkg/redis/*.tmpl
var go_arc_refresh embed.FS

type ProjectConfigRefresh struct {
	ProjectName string
	ModuleName  string
}

func CreateGoArcRefresh(projectName string) {
	config := ProjectConfigRefresh{
		ProjectName: projectName,
		ModuleName:  projectName,
	}

	fmt.Printf("Creating project: %s\n", projectName)
	fmt.Println("Generating structure...")

	if err := createProjectStructureRefresh(config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Project created successfully!")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  go run main.go\n")
}

func createProjectStructureRefresh(config ProjectConfigRefresh) error {
	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Tạo thư mục project
	projectDir := filepath.Join(baseDir, config.ProjectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	// Tạo tất cả thư mục cần thiết
	dirs := []string{
		"cmd/server/config",
		"cmd/server/routers",
		"internal/dto/auth/request",
		"internal/dto/auth/response",
		"internal/dto/refresh/request",
		"internal/dto/refresh/response",
		"internal/dto/user/request",
		"internal/dto/user/response",
		"internal/entity/auth",
		"internal/entity/refresh",
		"internal/entity/user",
		"internal/handler/auth",
		"internal/handler/refresh",
		"internal/handler/user",
		"internal/mapping",
		"internal/middleware",
		"internal/repository/auth",
		"internal/repository/refresh",
		"internal/repository/user",
		"internal/usecase/auth",
		"internal/usecase/refresh",
		"internal/usecase/user",
		"internal/utils/ratelimit",
		"logs",
		"migrate/migrations",
		"pkg/redis",
	}

	for _, dir := range dirs {
		path := filepath.Join(projectDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", path, err)
		}
	}

	// Map template files -> destination files
	files := map[string]string{
		// ROOT FILES
		"go_arc_refresh/main.go.tmpl":   "main.go",
		"go_arc_refresh/README.md.tmpl": "README.md",
		"go_arc_access/go.mod.tmpl":     "go.mod",

		// CMD
		"go_arc_refresh/cmd/root.go.tmpl":                     "cmd/root.go",
		"go_arc_refresh/cmd/server/config/config.go.tmpl":     "cmd/server/config/config.go",
		"go_arc_refresh/cmd/server/routers/composer.go.tmpl":  "cmd/server/routers/composer.go",
		"go_arc_refresh/cmd/server/routers/protected.go.tmpl": "cmd/server/routers/protected.go",
		"go_arc_refresh/cmd/server/routers/public.go.tmpl":    "cmd/server/routers/public.go",
		"go_arc_refresh/cmd/server/routers/router.go.tmpl":    "cmd/server/routers/router.go",

		// DTO AUTH REQUEST
		"go_arc_refresh/internal/dto/auth/request/login.go.tmpl":    "internal/dto/auth/request/login.go",
		"go_arc_refresh/internal/dto/auth/request/register.go.tmpl": "internal/dto/auth/request/register.go",
		"go_arc_refresh/internal/dto/auth/request/errors.go.tmpl":   "internal/dto/auth/request/errors.go",
		"go_arc_refresh/internal/dto/auth/request/validate.go.tmpl": "internal/dto/auth/request/validate.go",
		"go_arc_refresh/internal/dto/auth/request/gg.go.tmpl":       "internal/dto/auth/request/gg.go",

		// DTO AUTH RESPONSE
		"go_arc_refresh/internal/dto/auth/response/auth.go.tmpl": "internal/dto/auth/response/auth.go",

		// DTO REFRESH REQUEST
		"go_arc_refresh/internal/dto/refresh/request/refresh.go.tmpl": "internal/dto/refresh/request/refresh.go",

		// DTO REFRESH RESPONSE
		"go_arc_refresh/internal/dto/refresh/response/refresh.go.tmpl": "internal/dto/refresh/response/refresh.go",

		// DTO USER REQUEST
		"go_arc_refresh/internal/dto/user/request/errors.go.tmpl":   "internal/dto/user/request/errors.go",
		"go_arc_refresh/internal/dto/user/request/validate.go.tmpl": "internal/dto/user/request/validate.go",

		// DTO USER RESPONSE
		"go_arc_refresh/internal/dto/user/response/user.go.tmpl": "internal/dto/user/response/user.go",

		// ENTITY AUTH
		"go_arc_refresh/internal/entity/auth/auth.go.tmpl":        "internal/entity/auth/auth.go",
		"go_arc_refresh/internal/entity/auth/auth_google.go.tmpl": "internal/entity/auth/auth_google.go",

		// ENTITY REFRESH
		"go_arc_refresh/internal/entity/refresh/refesh_token.go.tmpl": "internal/entity/refresh/refesh_token.go",
		"go_arc_refresh/internal/entity/refresh/session.go.tmpl":      "internal/entity/refresh/session.go",

		// ENTITY USER
		"go_arc_refresh/internal/entity/user/user.go.tmpl": "internal/entity/user/user.go",

		// HANDLER AUTH
		"go_arc_refresh/internal/handler/auth/auth.go.tmpl":     "internal/handler/auth/auth.go",
		"go_arc_refresh/internal/handler/auth/login.go.tmpl":    "internal/handler/auth/login.go",
		"go_arc_refresh/internal/handler/auth/registry.go.tmpl": "internal/handler/auth/registry.go",
		"go_arc_refresh/internal/handler/auth/google.go.tmpl":   "internal/handler/auth/google.go",
		"go_arc_refresh/internal/handler/auth/logout.go.tmpl":   "internal/handler/auth/logout.go",

		// HANDLER REFRESH
		"go_arc_refresh/internal/handler/refresh/refresh.go.tmpl":  "internal/handler/refresh/refresh.go",
		"go_arc_refresh/internal/handler/refresh/re_token.go.tmpl": "internal/handler/refresh/re_token.go",

		// HANDLER USER
		"go_arc_refresh/internal/handler/user/get.go.tmpl":  "internal/handler/user/get.go",
		"go_arc_refresh/internal/handler/user/user.go.tmpl": "internal/handler/user/user.go",

		// MAPPING
		"go_arc_refresh/internal/mapping/auth.go.tmpl": "internal/mapping/auth.go",
		"go_arc_refresh/internal/mapping/key.go.tmpl":  "internal/mapping/key.go",
		"go_arc_refresh/internal/mapping/user.go.tmpl": "internal/mapping/user.go",

		// MIDDLEWARE
		"go_arc_refresh/internal/middleware/cors.go.tmpl":       "internal/middleware/cors.go",
		"go_arc_refresh/internal/middleware/middleware.go.tmpl": "internal/middleware/middleware.go",
		"go_arc_refresh/internal/middleware/rate_limit.go.tmpl": "internal/middleware/rate_limit.go",

		// REPOSITORY AUTH
		"go_arc_refresh/internal/repository/auth/auth_repo.go.tmpl":   "internal/repository/auth/auth_repo.go",
		"go_arc_refresh/internal/repository/auth/db_create.go.tmpl":   "internal/repository/auth/db_create.go",
		"go_arc_refresh/internal/repository/auth/db_get.go.tmpl":      "internal/repository/auth/db_get.go",
		"go_arc_refresh/internal/repository/auth/db_update.go.tmpl":   "internal/repository/auth/db_update.go",
		"go_arc_refresh/internal/repository/auth/redis_cache.go.tmpl": "internal/repository/auth/redis_cache.go",
		"go_arc_refresh/internal/repository/auth/redis_key.go.tmpl":   "internal/repository/auth/redis_key.go",

		// REPOSITORY REFRESH
		"go_arc_refresh/internal/repository/refresh/db_create.go.tmpl":    "internal/repository/refresh/db_create.go",
		"go_arc_refresh/internal/repository/refresh/db_delete.go.tmpl":    "internal/repository/refresh/db_delete.go",
		"go_arc_refresh/internal/repository/refresh/db_get.go.tmpl":       "internal/repository/refresh/db_get.go",
		"go_arc_refresh/internal/repository/refresh/db_revoke.go.tmpl":    "internal/repository/refresh/db_revoke.go",
		"go_arc_refresh/internal/repository/refresh/redis_cache.go.tmpl":  "internal/repository/refresh/redis_cache.go",
		"go_arc_refresh/internal/repository/refresh/redis_key.go.tmpl":    "internal/repository/refresh/redis_key.go",
		"go_arc_refresh/internal/repository/refresh/refresh_repo.go.tmpl": "internal/repository/refresh/refresh_repo.go",

		// REPOSITORY USER
		"go_arc_refresh/internal/repository/user/db_create.go.tmpl": "internal/repository/user/db_create.go",
		"go_arc_refresh/internal/repository/user/db_get.go.tmpl":    "internal/repository/user/db_get.go",
		"go_arc_refresh/internal/repository/user/user_repo.go.tmpl": "internal/repository/user/user_repo.go",

		// USECASE AUTH
		"go_arc_refresh/internal/usecase/auth/auth.go.tmpl":     "internal/usecase/auth/auth.go",
		"go_arc_refresh/internal/usecase/auth/errors.go.tmpl":   "internal/usecase/auth/errors.go",
		"go_arc_refresh/internal/usecase/auth/login.go.tmpl":    "internal/usecase/auth/login.go",
		"go_arc_refresh/internal/usecase/auth/register.go.tmpl": "internal/usecase/auth/register.go",
		"go_arc_refresh/internal/usecase/auth/get.go.tmpl":      "internal/usecase/auth/get.go",
		"go_arc_refresh/internal/usecase/auth/google.go.tmpl":   "internal/usecase/auth/google.go",
		"go_arc_refresh/internal/usecase/auth/logout.go.tmpl":   "internal/usecase/auth/logout.go",

		// USECASE REFRESH
		"go_arc_refresh/internal/usecase/refresh/create.go.tmpl":        "internal/usecase/refresh/create.go",
		"go_arc_refresh/internal/usecase/refresh/delete.go.tmpl":        "internal/usecase/refresh/delete.go",
		"go_arc_refresh/internal/usecase/refresh/get.go.tmpl":           "internal/usecase/refresh/get.go",
		"go_arc_refresh/internal/usecase/refresh/refresh.go.tmpl":       "internal/usecase/refresh/refresh.go",
		"go_arc_refresh/internal/usecase/refresh/refresh_token.go.tmpl": "internal/usecase/refresh/refresh_token.go",
		"go_arc_refresh/internal/usecase/refresh/revoke.go.tmpl":        "internal/usecase/refresh/revoke.go",
		"go_arc_refresh/internal/usecase/refresh/token.go.tmpl":         "internal/usecase/refresh/token.go",

		// USECASE USER
		"go_arc_refresh/internal/usecase/user/create.go.tmpl": "internal/usecase/user/create.go",
		"go_arc_refresh/internal/usecase/user/get.go.tmpl":    "internal/usecase/user/get.go",
		"go_arc_refresh/internal/usecase/user/user.go.tmpl":   "internal/usecase/user/user.go",

		// UTILS
		"go_arc_refresh/internal/utils/blacklist.go.tmpl": "internal/utils/blacklist.go",
		"go_arc_refresh/internal/utils/context.go.tmpl":   "internal/utils/context.go",
		"go_arc_refresh/internal/utils/helper.go.tmpl":    "internal/utils/helper.go",
		"go_arc_refresh/internal/utils/jwt.go.tmpl":       "internal/utils/jwt.go",
		"go_arc_refresh/internal/utils/session.go.tmpl":   "internal/utils/session.go",

		// UTILS RATELIMIT
		"go_arc_refresh/internal/utils/ratelimit/factory.go.tmpl":   "internal/utils/ratelimit/factory.go",
		"go_arc_refresh/internal/utils/ratelimit/interface.go.tmpl": "internal/utils/ratelimit/interface.go",
		"go_arc_refresh/internal/utils/ratelimit/memory.go.tmpl":    "internal/utils/ratelimit/memory.go",
		"go_arc_refresh/internal/utils/ratelimit/redis.go.tmpl":     "internal/utils/ratelimit/redis.go",

		// LOGS
		"go_arc_refresh/logs/logs.log.tmpl": "logs/logs.log",

		// MIGRATIONS
		"go_arc_refresh/migrate/migrations/1_create_table.up.sql.tmpl": "migrate/migrations/1_create_table.up.sql",
		"go_arc_refresh/migrate/migrations/1_drop_table.down.sql.tmpl": "migrate/migrations/1_drop_table.down.sql",

		// PKG REDIS
		"go_arc_refresh/pkg/redis/client.go.tmpl": "pkg/redis/client.go",
	}

	for tmpl, dest := range files {
		if err := renderTemplateRefresh(projectDir, dest, tmpl, config); err != nil {
			return fmt.Errorf("failed to render %s: %v", dest, err)
		}
	}

	return nil
}

func renderTemplateRefresh(baseDir, destFile, tmplFile string, config ProjectConfigRefresh) error {
	// Đọc template từ embed
	tmplContent, err := go_arc_refresh.ReadFile(tmplFile)
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
