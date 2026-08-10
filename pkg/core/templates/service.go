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

//go:embed service/dto/request/*.tmpl
//go:embed service/dto/response/*.tmpl
//go:embed service/entity/*.tmpl
//go:embed service/handler/*.tmpl
//go:embed service/mapping/*.tmpl
//go:embed service/repository/*.tmpl
//go:embed service/usecase/*.tmpl
//go:embed service/shared/model/model.go.tmpl
var templateSvFS embed.FS

type ServiceConfig struct {
	NameCap    string //
	Name       string //
	FolderName string //
	ModuleName string //
}

func (p ServiceConfig) GetMuduleName() string {
	return p.ModuleName
}
func AddGoService(project string, moduleName string, serviceName string) {
	// Xóa cách
	serviceName = strings.TrimSpace(serviceName)
	// Chỉnh chữ thường
	serviceName = strings.ToLower(serviceName)
	// Chữ Hoa đầu
	serviceNameCap := cases.Title(language.English).String(serviceName)
	serviceNameCap = strings.ReplaceAll(serviceNameCap, " ", "")
	serviceName = strings.ReplaceAll(serviceName, " ", "")
	config := ServiceConfig{
		NameCap:    serviceNameCap,
		Name:       serviceName,
		FolderName: serviceName + "-service",
		ModuleName: moduleName,
	}

	fmt.Printf("Generating service '%s' ...\n", config.FolderName)

	if err := generateService(project, config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Service '%s' generated successfully!\n", config.FolderName)
}

func generateService(project string, config ServiceConfig) error {
	// Vị trí project/internal/folder_name(new service)
	baseDir := filepath.Join(project, "internal", config.FolderName)

	// Các thư mục cần tạo
	dirs := []string{
		"dto/request",
		"dto/response",
		"entity",
		"handler",
		"mapping",
		"repository",
		"usecase",
	}

	for _, dir := range dirs {
		// Nối baseDir/folder_name(new service)/folder cần tạo ở dirs
		path := filepath.Join(baseDir, dir)
		// Tạo Folder
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", path, err)
		}
	}

	// Map template files -> destination files
	files := map[string]string{
		"service/dto/request/request.go.tmpl":   "dto/request/request.go",
		"service/dto/response/response.go.tmpl": "dto/response/response.go",
		"service/entity/entity.go.tmpl":         "entity/entity.go",
		"service/handler/handler.go.tmpl":       "handler/handler.go",
		"service/handler/create.go.tmpl":        "handler/create.go",
		"service/handler/get.go.tmpl":           "handler/get.go",
		"service/handler/update.go.tmpl":        "handler/update.go",
		"service/handler/delete.go.tmpl":        "handler/delete.go",
		"service/mapping/mapping.go.tmpl":       "mapping/mapping.go",
		"service/repository/repository.go.tmpl": "repository/repository.go",
		"service/repository/create.go.tmpl":     "repository/create.go",
		"service/repository/get.go.tmpl":        "repository/get.go",
		"service/repository/update.go.tmpl":     "repository/update.go",
		"service/repository/delete.go.tmpl":     "repository/delete.go",
		"service/usecase/usecase.go.tmpl":       "usecase/usecase.go",
		"service/usecase/create.go.tmpl":        "usecase/create.go",
		"service/usecase/get.go.tmpl":           "usecase/get.go",
		"service/usecase/update.go.tmpl":        "usecase/update.go",
		"service/usecase/delete.go.tmpl":        "usecase/delete.go",
	}

	for tmpl, dest := range files {
		utils.RenderTemplate(templateSvFS, baseDir, dest, tmpl, config)
	}

	// gen shared model
	if err := generateSharedModel(project, config); err != nil {
		fmt.Printf("Warning: failed to generate shared model: %v\n", err)
	}

	// Cập nhật key_object.go
	if err := updateKeyObject(project, config); err != nil {
		fmt.Printf("Warning: failed to update key_object.go: %v\n", err)
	}

	return nil
}

// Tạo file shared/model/Name.go
func generateSharedModel(projectRoot string, config ServiceConfig) error {
	sharedModelDir := filepath.Join(projectRoot, "internal", "shared", "model")
	if err := os.MkdirAll(sharedModelDir, 0755); err != nil {
		return err
	}

	destFile := filepath.Join(sharedModelDir, config.Name+".go")
	tmplFile := "service/shared/model/model.go.tmpl"

	tmplContent, err := templateSvFS.ReadFile(tmplFile)
	if err != nil {
		return fmt.Errorf("template %s not found", tmplFile)
	}

	tmpl, err := template.New("model.go.tmpl").Parse(string(tmplContent))
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
func updateKeyObject(projectRoot string, config ServiceConfig) error {
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
