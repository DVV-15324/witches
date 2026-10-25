package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	fetchTemplateFilesFromGitFunc = fetchTemplateFilesFromGit
	getCurrentModuleFunc          = GetCurrentModule
	validateModuleStructureFunc   = ValidateModuleStructure

	mkdirAllFunc  = os.MkdirAll
	writeFileFunc = os.WriteFile
	statFunc      = os.Stat

	addModuleFieldFunc       = AddModuleField
	addModuleInitFunc        = AddModuleInit
	addRouteRegistrationFunc = AddRouteRegistration
	updateKeyObjectFunc      = updateKeyObject
	fixAllImportsFunc        = FixAllImports
)

// Public API - hỗ trợ cấu trúc thư mục hiện tại
func AddGoDomainFromLink(project, moduleName, domainName, repoURL string) error {
	domainName = strings.TrimSpace(domainName)
	domainName = strings.ReplaceAll(domainName, " ", "")
	domainName = strings.ToLower(domainName)
	domainNameCap := cases.Title(language.English).String(domainName)
	domainNameCap = strings.ReplaceAll(domainNameCap, " ", "")
	config := ModuleConfig{
		NameCap:     domainNameCap,
		Name:        domainName,
		ProjectName: moduleName,
	}
	sourceModule := GetModuleNameFromRepo(repoURL)
	targetModule, err := getCurrentModuleFunc(project)
	if err != nil {
		return fmt.Errorf("failed to get target module: %w", err)
	}
	fmt.Printf("  Source module: %s\n", sourceModule)
	fmt.Printf("  Target module: %s\n", targetModule)
	templateFiles, err := fetchTemplateFilesFromGitFunc(repoURL, domainName)
	if err != nil {
		return fmt.Errorf("fetch templates: %w", err)
	}
	if len(templateFiles) == 0 {
		return fmt.Errorf("no files found in internal/%s/ or internal/shared/domain/", domainName)
	}
	fmt.Printf("Found %d template files\n", len(templateFiles))
	if err := validateModuleStructureFunc(templateFiles, config); err != nil {
		return fmt.Errorf("invalid module structure: %w", err)
	}
	fmt.Println("Validation passed")
	fmt.Println("Writing files to disk...")
	baseDir := filepath.Join(project, "internal", config.Name)
	const prefix = "internal/"
	sharedDomainDir := filepath.Join(project, "internal", "shared", "domain")
	if err := mkdirAllFunc(sharedDomainDir, 0755); err != nil {
		return fmt.Errorf("failed to create shared domain directory: %w", err)
	}
	for tmplPath, content := range templateFiles {
		var destPath string
		var relPath string
		var isSharedDomain bool
		var isMigration bool
		fmt.Printf("  Processing: %s\n", tmplPath)
		if strings.HasPrefix(tmplPath, "internal/shared/domain/") {
			isSharedDomain = true
			fileName := filepath.Base(tmplPath)
			relPath = filepath.Join("internal", "shared", "domain", fileName)
			destPath = filepath.Join(project, relPath)
		} else if strings.HasPrefix(tmplPath, prefix) {
			relPath = strings.TrimPrefix(tmplPath, prefix)
			if relPath == tmplPath {
				return fmt.Errorf("invalid template path: %s", tmplPath)
			}
			parts := strings.Split(relPath, "/")
			if len(parts) < 2 {
				return fmt.Errorf(
					"invalid path structure: %s",
					relPath,
				)
			}
			if parts[0] != domainName {
				return fmt.Errorf("template belongs to another domain: %s", tmplPath)
			}
			filePath := strings.Join(parts[1:], "/")
			if strings.HasPrefix(filePath, "migrate/migrations/") {
				isMigration = true
			}
			destPath = filepath.Join(baseDir, filePath)
		} else {
			continue
		}
		var buf strings.Builder
		if strings.HasSuffix(tmplPath, ".tmpl") {
			tmpl, err := template.New(filepath.Base(tmplPath)).Parse(content)
			if err != nil {
				return fmt.Errorf("parse template %s: %w", tmplPath, err)
			}

			if err := tmpl.Execute(&buf, config); err != nil {
				return fmt.Errorf("execute template %s: %w", tmplPath, err)
			}
		} else {
			buf.WriteString(content)
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		if err := writeFileFunc(destPath, []byte(buf.String()), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}

		if isSharedDomain {
			fmt.Printf("  Generated shared domain: %s\n", relPath)

		} else if isMigration {
			fmt.Printf("Generated migration: %s\n", filepath.Join("internal", domainName, strings.TrimPrefix(strings.TrimPrefix(relPath, domainName+"/"), "/")))
		} else {
			fmt.Printf("Generated: %s\n", filepath.Join(domainName, strings.TrimPrefix(relPath, domainName+"/")))
		}
	}

	fmt.Println("\nRewriting imports...")

	fmt.Println("\nAdding shared domain imports...")
	filesToUpdate := []string{
		filepath.Join(project, "internal", domainName, "model", "model.go"),
		filepath.Join(project, "internal", domainName, "repository", "repository.go"),
		filepath.Join(project, "internal", domainName, "usecase", "usecase.go"),
		filepath.Join(project, "internal", domainName, "handler", "handler.go"),
		filepath.Join(project, "internal", domainName, "module.go"),
	}

	for _, file := range filesToUpdate {

		if _, err := statFunc(file); os.IsNotExist(err) {
			fmt.Printf("File not found: %s\n", filepath.Base(file))
			continue
		}
	}

	fmt.Println("\nUpdating modules.go and routers.go...")
	modulesPath := filepath.Join(project, "cmd", "server", "routers", "modules.go")
	if _, err := statFunc(modulesPath); err == nil {
		if err := addModuleFieldFunc(modulesPath, config.Name, config.NameCap, moduleName); err != nil {
			fmt.Printf("Warning: failed to add module field: %v\n", err)
		} else {
			fmt.Println("Updated modules.go: added import and field")
		}
		if err := addModuleInitFunc(modulesPath, config.Name, config.NameCap); err != nil {
			fmt.Printf("Warning: failed to add module init: %v\n", err)
		} else {
			fmt.Println("Updated modules.go: added initialization")
		}
	} else {
		fmt.Printf("modules.go not found at: %s\n", modulesPath)
	}
	routersPath := filepath.Join(project, "cmd", "server", "routers", "routers.go")
	if _, err := os.Stat(routersPath); err == nil {
		if err := addRouteRegistrationFunc(routersPath, config.Name, config.NameCap); err != nil {
			fmt.Printf("Warning: failed to add route registration: %v\n", err)
		} else {
			fmt.Println("Updated routers.go: added route registration for", config.Name)
		}
	} else {
		fmt.Printf("routers.go not found at: %s\n", routersPath)
	}
	if err := updateKeyObjectFunc(project, config); err != nil {
		fmt.Printf("Warning: failed to update key_object.go: %v\n", err)
	} else {
		fmt.Println("Updated key_object.go")
	}
	fmt.Println("\nFixing imports in entire project...")
	if err := fixAllImportsFunc(project, sourceModule, targetModule); err != nil {
		fmt.Printf("Warning: failed to fix all imports: %v\n", err)
	} else {
		fmt.Println("Fixed all imports in project")
	}
	fmt.Printf("\nDomain '%s' generated successfully from external repo!\n", config.Name)
	fmt.Println("\nNext steps:")
	fmt.Printf("1. Check generated files in internal/%s/\n", domainName)
	fmt.Printf("2. Check migration files in internal/%s/migrate/migrations/\n", domainName)
	fmt.Println("3. Review module initialization in cmd/server/routers/modules.go")
	fmt.Println("4. Review routes in cmd/server/routers/routers.go")
	fmt.Println("5. Run 'go mod tidy' to clean up dependencies")
	return nil
}
