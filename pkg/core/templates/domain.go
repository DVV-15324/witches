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

//go:embed domain/dto/request/*.tmpl
//go:embed domain/dto/response/*.tmpl
//go:embed domain/model/*.tmpl
//go:embed domain/handler/*.tmpl
//go:embed domain/mapping/*.tmpl
//go:embed domain/repository/*.tmpl
//go:embed domain/usecase/*.tmpl
//go:embed domain/shared/domain/domain.go.tmpl
//go:embed domain/module.go.tmpl
//go:embed domain/migrate/mssql/*.tmpl
//go:embed domain/migrate/mysql/*.tmpl
//go:embed domain/migrate/postgresql/*.tmpl
var templateDomainFS embed.FS

type DomainConfig struct {
	NameCap    string
	Name       string
	FolderName string
	ModuleName string
}

func (p DomainConfig) GetModuleName() string {
	return p.ModuleName
}

func AddDomain(project string, moduleName string, domainName string, typeDb string) {
	domainName = strings.TrimSpace(domainName)
	domainName = strings.ReplaceAll(domainName, " ", "")
	domainName = strings.ToLower(domainName)
	domainNameCap := cases.Title(language.English).String(domainName)
	domainNameCap = strings.ReplaceAll(domainNameCap, " ", "")
	config := DomainConfig{
		NameCap:    domainNameCap,
		Name:       domainName,
		FolderName: domainName,
		ModuleName: moduleName,
	}

	fmt.Printf("Generating domain '%s' ...\n", config.FolderName)

	if err := generateDomain(project, config, typeDb); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	modulesPath := filepath.Join(project, "cmd", "server", "routers", "modules.go")
	if err := AddModuleField(modulesPath, config.Name, config.NameCap, moduleName); err != nil {
		fmt.Printf("Warning: failed to add module field: %v\n", err)
	} else {
		fmt.Println("Updated modules.go: added import and field")
	}

	if err := AddModuleInit(modulesPath, config.Name, config.NameCap); err != nil {
		fmt.Printf("Warning: failed to add module init: %v\n", err)
	} else {
		fmt.Println("Updated modules.go: added initialization")
	}
	routersPath := filepath.Join(project, "cmd", "server", "routers", "routers.go")
	if err := AddRouteRegistration(routersPath, config.Name, config.NameCap); err != nil {
		fmt.Printf("Warning: failed to add route registration: %v\n", err)
	} else {
		fmt.Println("Updated routers.go: added route registration for", config.Name)
	}
	fmt.Printf("domain '%s' generated successfully!\n", config.FolderName)
}

func generateDomain(project string, config DomainConfig, typeDb string) error {
	baseDir := filepath.Join(project, "internal", config.FolderName)

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
		migraUp = "domain/migrate/mysql/1_create_table.up.sql.tmpl"
		migraDown = "domain/migrate/mysql/1_drop_table.down.sql.tmpl"
	case "postgres", "postgresql":
		migraUp = "domain/migrate/postgresql/1_create_table.up.sql.tmpl"
		migraDown = "domain/migrate/postgresql/1_drop_table.down.sql.tmpl"
	case "sqlserver", "mssql":
		migraUp = "domain/migrate/mssql/1_create_table.up.sql.tmpl"
		migraDown = "domain/migrate/mssql/1_drop_table.down.sql.tmpl"
	default:
		return fmt.Errorf("error: unsupported database: %s. supported : mysql, postgresql, postgres, mssql, sqlserver", typeDb)
	}
	files := map[string]string{
		"domain/dto/request/request.go.tmpl":   "dto/request/request.go",
		"domain/dto/request/validate.go.tmpl":  "dto/request/validate.go",
		"domain/dto/request/errors.go.tmpl":    "dto/request/errors.go",
		"domain/dto/response/response.go.tmpl": "dto/response/response.go",
		"domain/model/model.go.tmpl":           "model/model.go",
		"domain/handler/handler.go.tmpl":       "handler/handler.go",
		"domain/handler/create.go.tmpl":        "handler/create.go",
		"domain/handler/get.go.tmpl":           "handler/get.go",
		"domain/handler/update.go.tmpl":        "handler/update.go",
		"domain/handler/delete.go.tmpl":        "handler/delete.go",
		"domain/mapping/mapping.go.tmpl":       "mapping/mapping.go",
		"domain/repository/repository.go.tmpl": "repository/repository.go",
		"domain/repository/create.go.tmpl":     "repository/create.go",
		"domain/repository/get.go.tmpl":        "repository/get.go",
		"domain/repository/update.go.tmpl":     "repository/update.go",
		"domain/repository/delete.go.tmpl":     "repository/delete.go",
		"domain/usecase/usecase.go.tmpl":       "usecase/usecase.go",
		"domain/usecase/create.go.tmpl":        "usecase/create.go",
		"domain/usecase/get.go.tmpl":           "usecase/get.go",
		"domain/usecase/update.go.tmpl":        "usecase/update.go",
		"domain/usecase/delete.go.tmpl":        "usecase/delete.go",
		"domain/module.go.tmpl":                "module.go",
		// MIGRATIONS
		migraUp:   "migrate/migrations/1_create_table.up.sql",
		migraDown: "migrate/migrations/1_drop_table.down.sql",
	}

	for tmpl, dest := range files {
		utils.RenderTemplate(templateDomainFS, baseDir, dest, tmpl, config)
	}

	if err := generateSharedDomain(project, config, templateDomainFS); err != nil {
		fmt.Printf("Warning: failed to generate shared domain: %v\n", err)
	}

	if err := updateKeyObject(project, config); err != nil {
		fmt.Printf("Warning: failed to update key_object.go: %v\n", err)
	}

	return nil
}
