package template

import (
	"embed"
	"fmt"
	"github.com/DVV-15324/witches/pkg/core/templates/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strings"
)

//go:embed module/dto/request/*.tmpl
//go:embed module/dto/response/*.tmpl
//go:embed module/model/*.tmpl
//go:embed module/handler/*.tmpl
//go:embed module/mapping/*.tmpl
//go:embed module/repository/*.tmpl
//go:embed module/usecase/*.tmpl
//go:embed module/shared/domain/domain.go.tmpl
//go:embed module/module.go.tmpl
//go:embed module/module.env.tmpl
//go:embed module/migrate/mssql/*.tmpl
//go:embed module/migrate/mysql/*.tmpl
//go:embed module/migrate/postgresql/*.tmpl
var templateModuleFS embed.FS

type ModuleConfig struct {
	NameCap     string
	Name        string
	ProjectName string
}

func (p ModuleConfig) GetProjectName() string {
	return p.ProjectName
}

func AddModule(projectPath string, projectName string, moduleName string, typeDb string) {
	moduleName = strings.TrimSpace(moduleName)
	moduleName = strings.ReplaceAll(moduleName, " ", "")
	moduleName = strings.ToLower(moduleName)
	moduleNameCap := cases.Title(language.English).String(moduleName)
	moduleNameCap = strings.ReplaceAll(moduleNameCap, " ", "")
	config := ModuleConfig{
		NameCap:     moduleNameCap,
		Name:        moduleName,
		ProjectName: projectName,
	}

	fmt.Printf("Generating module '%s' ...\n", config.Name)

	if err := generateModule(projectPath, config, typeDb); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	modulesPath := filepath.Join(projectPath, "cmd", "server", "routers", "modules.go")
	if err := AddModuleField(modulesPath, config.Name, config.NameCap, projectName); err != nil {
		fmt.Printf("Warning: failed to add module field: %v\n", err)
	} else {
		fmt.Println("Updated modules.go: added import and field")
	}

	if err := AddModuleInit(modulesPath, config.Name, config.NameCap); err != nil {
		fmt.Printf("Warning: failed to add module init: %v\n", err)
	} else {
		fmt.Println("Updated modules.go: added initialization")
	}
	routersPath := filepath.Join(projectPath, "cmd", "server", "routers", "routers.go")
	if err := AddRouteRegistration(routersPath, config.Name, config.NameCap); err != nil {
		fmt.Printf("Warning: failed to add route registration: %v\n", err)
	} else {
		fmt.Println("Updated routers.go: added route registration for", config.Name)
	}
	fmt.Printf("module '%s' generated successfully!\n", config.Name)
}

func generateModule(projectPath string, config ModuleConfig, typeDb string) error {
	baseDir := filepath.Join(projectPath, "internal", config.Name)

	dirs := []string{
		"dto/request",
		"dto/response",
		"model",
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
	var migraUp string
	var migraDown string
	switch typeDb {
	case "mysql":
		migraUp = "module/migrate/mysql/1_create_table.up.sql.tmpl"
		migraDown = "module/migrate/mysql/1_drop_table.down.sql.tmpl"
	case "postgres", "postgresql":
		migraUp = "module/migrate/postgresql/1_create_table.up.sql.tmpl"
		migraDown = "module/migrate/postgresql/1_drop_table.down.sql.tmpl"
	case "sqlserver", "mssql":
		migraUp = "module/migrate/mssql/1_create_table.up.sql.tmpl"
		migraDown = "module/migrate/mssql/1_drop_table.down.sql.tmpl"
	default:
		return fmt.Errorf("error: unsupported database: %s. supported : mysql, postgresql, postgres, mssql, sqlserver", typeDb)
	}
	files := map[string]string{
		"module/dto/request/request.go.tmpl":   "dto/request/request.go",
		"module/dto/request/validate.go.tmpl":  "dto/request/validate.go",
		"module/dto/request/errors.go.tmpl":    "dto/request/errors.go",
		"module/dto/response/response.go.tmpl": "dto/response/response.go",
		"module/model/model.go.tmpl":           "model/model.go",
		"module/handler/handler.go.tmpl":       "handler/handler.go",
		"module/handler/create.go.tmpl":        "handler/create.go",
		"module/handler/get.go.tmpl":           "handler/get.go",
		"module/handler/update.go.tmpl":        "handler/update.go",
		"module/handler/delete.go.tmpl":        "handler/delete.go",
		"module/mapping/mapping.go.tmpl":       "mapping/mapping.go",
		"module/repository/repository.go.tmpl": "repository/repository.go",
		"module/repository/create.go.tmpl":     "repository/create.go",
		"module/repository/get.go.tmpl":        "repository/get.go",
		"module/repository/update.go.tmpl":     "repository/update.go",
		"module/repository/delete.go.tmpl":     "repository/delete.go",
		"module/usecase/usecase.go.tmpl":       "usecase/usecase.go",
		"module/usecase/create.go.tmpl":        "usecase/create.go",
		"module/usecase/get.go.tmpl":           "usecase/get.go",
		"module/usecase/update.go.tmpl":        "usecase/update.go",
		"module/usecase/delete.go.tmpl":        "usecase/delete.go",
		"module/module.go.tmpl":                "module.go",
		"module/module.env.tmpl":               "module.env",
		// MIGRATIONS
		migraUp:   "migrate/migrations/1_create_table.up.sql",
		migraDown: "migrate/migrations/1_drop_table.down.sql",
	}

	for tmpl, dest := range files {
		utils.RenderTemplate(templateModuleFS, baseDir, dest, tmpl, config)
	}

	if err := generateSharedDomain(projectPath, config, templateModuleFS); err != nil {
		fmt.Printf("Warning: failed to generate shared domain: %v\n", err)
	}

	if err := updateKeyObject(projectPath, config); err != nil {
		fmt.Printf("Warning: failed to update key_object.go: %v\n", err)
	}

	return nil
}
