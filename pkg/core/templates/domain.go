package template

import (
	"embed"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strconv"

	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/DVV-15324/witches/pkg/core/templates/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed domain/dto/request/*.tmpl
//go:embed domain/dto/response/*.tmpl
//go:embed domain/model/*.tmpl
//go:embed domain/handler/*.tmpl
//go:embed domain/mapping/*.tmpl
//go:embed domain/repository/*.tmpl
//go:embed domain/usecase/*.tmpl
//go:embed domain/shared/domain/domain.go.tmpl
var templateSvFS embed.FS

type DomainConfig struct {
	NameCap    string //
	Name       string //
	FolderName string //
	ModuleName string //
}

func (p DomainConfig) GetMuduleName() string {
	return p.ModuleName
}
func AddGoDomain(project string, moduleName string, domainName string) {
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

	if err := generateDomain(project, config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("domain '%s' generated successfully!\n", config.FolderName)
}

func generateDomain(project string, config DomainConfig) error {
	// Vị trí project/internal/folder_name(new domain)
	baseDir := filepath.Join(project, "internal", config.FolderName)

	// Các thư mục cần tạo
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
		// Nối baseDir/folder_name(new domain)/folder cần tạo ở dirs
		path := filepath.Join(baseDir, dir)
		// Tạo Folder
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", path, err)
		}
	}

	// Map template files -> destination files
	files := map[string]string{
		"domain/dto/request/request.go.tmpl":   "dto/request/request.go",
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
	}

	for tmpl, dest := range files {
		utils.RenderTemplate(templateSvFS, baseDir, dest, tmpl, config)
	}

	// gen shared domain
	if err := generateSharedDomain(project, config); err != nil {
		fmt.Printf("Warning: failed to generate shared domain: %v\n", err)
	}

	// Cập nhật key_object.go
	if err := updateKeyObject(project, config); err != nil {
		fmt.Printf("Warning: failed to update key_object.go: %v\n", err)
	}

	return nil
}

// Tạo file shared/domain/Name.go
func generateSharedDomain(projectRoot string, config DomainConfig) error {
	sharedDomainDir := filepath.Join(projectRoot, "internal", "shared", "domain")
	if err := os.MkdirAll(sharedDomainDir, 0755); err != nil {
		return err
	}

	destFile := filepath.Join(sharedDomainDir, config.Name+".go")
	tmplFile := "domain/shared/domain/domain.go.tmpl"

	tmplContent, err := templateSvFS.ReadFile(tmplFile)
	if err != nil {
		return fmt.Errorf("template %s not found", tmplFile)
	}

	tmpl, err := template.New("domain.go.tmpl").Parse(string(tmplContent))
	if err != nil {
		return err
	}

	file, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, config)
}

// Cập nhật key_object
func updateKeyObject(projectRoot string, config DomainConfig) error {
	keyFile := filepath.Join(projectRoot, "internal", "shared", "utils", "key_object.go")

	src, err := os.ReadFile(keyFile)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file error: %w", err)
	}

	// Tìm object có ID lớn nhất và block var chứa nó
	var maxID int
	var targetDecl *ast.GenDecl

	ast.Inspect(node, func(n ast.Node) bool {
		// ast.GenDecl (quản lý các thành phần ở cấp package )
		decl, ok := n.(*ast.GenDecl)
		if !ok || decl.Tok != token.VAR {
			return true
		}

		for _, spec := range decl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || !strings.HasPrefix(vs.Names[0].Name, "Object") {
				continue
			}

			if lit, ok := vs.Values[0].(*ast.BasicLit); ok && lit.Kind == token.INT {
				if id, _ := strconv.Atoi(lit.Value); id > maxID {
					maxID = id
					targetDecl = decl
				}
			}
		}
		return true
	})

	if targetDecl == nil {
		return fmt.Errorf("no Object constant found")
	}

	// Tạo và thêm constant mới
	newSpec := &ast.ValueSpec{
		Names:  []*ast.Ident{ast.NewIdent("Object" + config.NameCap)},
		Type:   &ast.Ident{Name: "uint"},
		Values: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", maxID+1)}},
	}
	targetDecl.Specs = append(targetDecl.Specs, newSpec)

	var buf strings.Builder
	if err := format.Node(&buf, fset, node); err != nil {
		return fmt.Errorf("format code error: %w", err)
	}

	return os.WriteFile(keyFile, []byte(buf.String()), 0644)
}
