package template

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed template/*.tmpl
//go:embed template/cmd/*.tmpl
//go:embed template/cmd/server/config/*.tmpl
//go:embed template/cmd/server/routers/*.tmpl
//go:embed template/internal/dto/auth/request/*.tmpl
//go:embed template/internal/dto/auth/response/*.tmpl
//go:embed template/internal/dto/refresh/request/*.tmpl
//go:embed template/internal/dto/refresh/response/*.tmpl
//go:embed template/internal/dto/user/request/*.tmpl
//go:embed template/internal/dto/user/response/*.tmpl
//go:embed template/internal/entity/auth/*.tmpl
//go:embed template/internal/entity/refresh/*.tmpl
//go:embed template/internal/entity/user/*.tmpl
//go:embed template/internal/handler/auth/*.tmpl
//go:embed template/internal/handler/refresh/*.tmpl
//go:embed template/internal/handler/user/*.tmpl
//go:embed template/internal/mapping/*.tmpl
//go:embed template/internal/middleware/*.tmpl
//go:embed template/internal/repository/auth/*.tmpl
//go:embed template/internal/repository/refresh/*.tmpl
//go:embed template/internal/repository/user/*.tmpl
//go:embed template/internal/usecase/auth/*.tmpl
//go:embed template/internal/usecase/refresh/*.tmpl
//go:embed template/internal/usecase/user/*.tmpl
//go:embed template/internal/utils/*.tmpl
//go:embed template/internal/utils/ratelimit/*.tmpl
//go:embed template/logs/*.tmpl
//go:embed template/migrate/migrations/*.tmpl
//go:embed template/pkg/redis/*.tmpl
var templates embed.FS

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
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  witches install\n")
	fmt.Printf("  witches run\n")
}

func createProjectStructureRefresh(config ProjectConfigRefresh) error {
	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Tạo thư mục project
	if err := os.MkdirAll(baseDir, 0755); err != nil {
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
		path := filepath.Join(baseDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", path, err)
		}
	}

	// Map template files -> destination files
	files := map[string]string{
		// ROOT FILES
		"template/main.go.tmpl":   "main.go",
		"template/README.md.tmpl": "README.md",
		"template/go.mod.tmpl":    "go.mod",

		// CMD
		"template/cmd/root.go.tmpl":                     "cmd/root.go",
		"template/cmd/server/config/config.go.tmpl":     "cmd/server/config/config.go",
		"template/cmd/server/routers/composer.go.tmpl":  "cmd/server/routers/composer.go",
		"template/cmd/server/routers/protected.go.tmpl": "cmd/server/routers/protected.go",
		"template/cmd/server/routers/public.go.tmpl":    "cmd/server/routers/public.go",
		"template/cmd/server/routers/router.go.tmpl":    "cmd/server/routers/router.go",

		// DTO AUTH REQUEST
		"template/internal/dto/auth/request/login.go.tmpl":    "internal/dto/auth/request/login.go",
		"template/internal/dto/auth/request/register.go.tmpl": "internal/dto/auth/request/register.go",
		"template/internal/dto/auth/request/errors.go.tmpl":   "internal/dto/auth/request/errors.go",
		"template/internal/dto/auth/request/validate.go.tmpl": "internal/dto/auth/request/validate.go",
		"template/internal/dto/auth/request/gg.go.tmpl":       "internal/dto/auth/request/gg.go",

		// DTO AUTH RESPONSE
		"template/internal/dto/auth/response/auth.go.tmpl": "internal/dto/auth/response/auth.go",

		// DTO REFRESH REQUEST
		"template/internal/dto/refresh/request/refresh.go.tmpl": "internal/dto/refresh/request/refresh.go",

		// DTO REFRESH RESPONSE
		"template/internal/dto/refresh/response/refresh.go.tmpl": "internal/dto/refresh/response/refresh.go",

		// DTO USER REQUEST
		"template/internal/dto/user/request/errors.go.tmpl":   "internal/dto/user/request/errors.go",
		"template/internal/dto/user/request/validate.go.tmpl": "internal/dto/user/request/validate.go",

		// DTO USER RESPONSE
		"template/internal/dto/user/response/user.go.tmpl": "internal/dto/user/response/user.go",

		// ENTITY AUTH
		"template/internal/entity/auth/auth.go.tmpl":        "internal/entity/auth/auth.go",
		"template/internal/entity/auth/auth_google.go.tmpl": "internal/entity/auth/auth_google.go",

		// ENTITY REFRESH
		"template/internal/entity/refresh/refesh_token.go.tmpl": "internal/entity/refresh/refesh_token.go",
		"template/internal/entity/refresh/session.go.tmpl":      "internal/entity/refresh/session.go",

		// ENTITY USER
		"template/internal/entity/user/user.go.tmpl": "internal/entity/user/user.go",

		// HANDLER AUTH
		"template/internal/handler/auth/auth.go.tmpl":     "internal/handler/auth/auth.go",
		"template/internal/handler/auth/login.go.tmpl":    "internal/handler/auth/login.go",
		"template/internal/handler/auth/registry.go.tmpl": "internal/handler/auth/registry.go",
		"template/internal/handler/auth/google.go.tmpl":   "internal/handler/auth/google.go",
		"template/internal/handler/auth/logout.go.tmpl":   "internal/handler/auth/logout.go",

		// HANDLER REFRESH
		"template/internal/handler/refresh/refresh.go.tmpl":  "internal/handler/refresh/refresh.go",
		"template/internal/handler/refresh/re_token.go.tmpl": "internal/handler/refresh/re_token.go",

		// HANDLER USER
		"template/internal/handler/user/get.go.tmpl":  "internal/handler/user/get.go",
		"template/internal/handler/user/user.go.tmpl": "internal/handler/user/user.go",

		// MAPPING
		"template/internal/mapping/auth.go.tmpl": "internal/mapping/auth.go",
		"template/internal/mapping/key.go.tmpl":  "internal/mapping/key.go",
		"template/internal/mapping/user.go.tmpl": "internal/mapping/user.go",

		// MIDDLEWARE
		"template/internal/middleware/cors.go.tmpl":       "internal/middleware/cors.go",
		"template/internal/middleware/middleware.go.tmpl": "internal/middleware/middleware.go",
		"template/internal/middleware/rate_limit.go.tmpl": "internal/middleware/rate_limit.go",

		// REPOSITORY AUTH
		"template/internal/repository/auth/auth_repo.go.tmpl":   "internal/repository/auth/auth_repo.go",
		"template/internal/repository/auth/db_create.go.tmpl":   "internal/repository/auth/db_create.go",
		"template/internal/repository/auth/db_get.go.tmpl":      "internal/repository/auth/db_get.go",
		"template/internal/repository/auth/db_update.go.tmpl":   "internal/repository/auth/db_update.go",
		"template/internal/repository/auth/redis_cache.go.tmpl": "internal/repository/auth/redis_cache.go",
		"template/internal/repository/auth/redis_key.go.tmpl":   "internal/repository/auth/redis_key.go",

		// REPOSITORY REFRESH
		"template/internal/repository/refresh/db_create.go.tmpl":    "internal/repository/refresh/db_create.go",
		"template/internal/repository/refresh/db_delete.go.tmpl":    "internal/repository/refresh/db_delete.go",
		"template/internal/repository/refresh/db_get.go.tmpl":       "internal/repository/refresh/db_get.go",
		"template/internal/repository/refresh/db_revoke.go.tmpl":    "internal/repository/refresh/db_revoke.go",
		"template/internal/repository/refresh/redis_cache.go.tmpl":  "internal/repository/refresh/redis_cache.go",
		"template/internal/repository/refresh/redis_key.go.tmpl":    "internal/repository/refresh/redis_key.go",
		"template/internal/repository/refresh/refresh_repo.go.tmpl": "internal/repository/refresh/refresh_repo.go",

		// REPOSITORY USER
		"template/internal/repository/user/db_create.go.tmpl": "internal/repository/user/db_create.go",
		"template/internal/repository/user/db_get.go.tmpl":    "internal/repository/user/db_get.go",
		"template/internal/repository/user/user_repo.go.tmpl": "internal/repository/user/user_repo.go",

		// USECASE AUTH
		"template/internal/usecase/auth/auth.go.tmpl":     "internal/usecase/auth/auth.go",
		"template/internal/usecase/auth/errors.go.tmpl":   "internal/usecase/auth/errors.go",
		"template/internal/usecase/auth/login.go.tmpl":    "internal/usecase/auth/login.go",
		"template/internal/usecase/auth/register.go.tmpl": "internal/usecase/auth/register.go",
		"template/internal/usecase/auth/get.go.tmpl":      "internal/usecase/auth/get.go",
		"template/internal/usecase/auth/google.go.tmpl":   "internal/usecase/auth/google.go",
		"template/internal/usecase/auth/logout.go.tmpl":   "internal/usecase/auth/logout.go",

		// USECASE REFRESH
		"template/internal/usecase/refresh/create.go.tmpl":        "internal/usecase/refresh/create.go",
		"template/internal/usecase/refresh/delete.go.tmpl":        "internal/usecase/refresh/delete.go",
		"template/internal/usecase/refresh/get.go.tmpl":           "internal/usecase/refresh/get.go",
		"template/internal/usecase/refresh/refresh.go.tmpl":       "internal/usecase/refresh/refresh.go",
		"template/internal/usecase/refresh/refresh_token.go.tmpl": "internal/usecase/refresh/refresh_token.go",
		"template/internal/usecase/refresh/revoke.go.tmpl":        "internal/usecase/refresh/revoke.go",
		"template/internal/usecase/refresh/token.go.tmpl":         "internal/usecase/refresh/token.go",

		// USECASE USER
		"template/internal/usecase/user/create.go.tmpl": "internal/usecase/user/create.go",
		"template/internal/usecase/user/get.go.tmpl":    "internal/usecase/user/get.go",
		"template/internal/usecase/user/user.go.tmpl":   "internal/usecase/user/user.go",

		// UTILS
		"template/internal/utils/blacklist.go.tmpl": "internal/utils/blacklist.go",
		"template/internal/utils/context.go.tmpl":   "internal/utils/context.go",
		"template/internal/utils/helper.go.tmpl":    "internal/utils/helper.go",
		"template/internal/utils/jwt.go.tmpl":       "internal/utils/jwt.go",
		"template/internal/utils/session.go.tmpl":   "internal/utils/session.go",

		// UTILS RATELIMIT
		"template/internal/utils/ratelimit/factory.go.tmpl":   "internal/utils/ratelimit/factory.go",
		"template/internal/utils/ratelimit/interface.go.tmpl": "internal/utils/ratelimit/interface.go",
		"template/internal/utils/ratelimit/memory.go.tmpl":    "internal/utils/ratelimit/memory.go",
		"template/internal/utils/ratelimit/redis.go.tmpl":     "internal/utils/ratelimit/redis.go",

		// LOGS
		"template/logs/logs.log.tmpl": "logs/logs.log",

		// MIGRATIONS
		"template/migrate/migrations/1_create_table.up.sql.tmpl": "migrate/migrations/1_create_table.up.sql",
		"template/migrate/migrations/1_drop_table.down.sql.tmpl": "migrate/migrations/1_drop_table.down.sql",

		// PKG REDIS
		"template/pkg/redis/client.go.tmpl": "pkg/redis/client.go",
	}

	for tmpl, dest := range files {
		if err := renderTemplateRefresh(baseDir, dest, tmpl, config); err != nil {
			return fmt.Errorf("failed to render %s: %v", dest, err)
		}
	}

	return nil
}

func renderTemplateRefresh(baseDir, destFile, tmplFile string, config ProjectConfigRefresh) error {
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
