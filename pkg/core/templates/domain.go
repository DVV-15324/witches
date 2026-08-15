package template

import (
	"bytes"
	"embed"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
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
//go:embed domain/module.go.tmpl
var templateSvFS embed.FS

type DomainConfig struct {
	NameCap    string
	Name       string
	FolderName string
	ModuleName string
}

func (p DomainConfig) GetModuleName() string {
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

	// Cập nhật modules.go (thêm import, field, init)
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
		fmt.Println("Updated routers.go: added route registration")
	}
	fmt.Printf("domain '%s' generated successfully!\n", config.FolderName)
}

func generateDomain(project string, config DomainConfig) error {
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
		"domain/module.go.tmpl":                "module.go",
	}

	for tmpl, dest := range files {
		utils.RenderTemplate(templateSvFS, baseDir, dest, tmpl, config)
	}

	if err := generateSharedDomain(project, config); err != nil {
		fmt.Printf("Warning: failed to generate shared domain: %v\n", err)
	}

	if err := updateKeyObject(project, config); err != nil {
		fmt.Printf("Warning: failed to update key_object.go: %v\n", err)
	}

	return nil
}

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

	var maxID int
	var targetDecl *ast.GenDecl

	ast.Inspect(node, func(n ast.Node) bool {
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

	newSpec := &ast.ValueSpec{
		Names:  []*ast.Ident{ast.NewIdent("Object" + config.NameCap)},
		Type:   &ast.Ident{Name: "uint"},
		Values: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", maxID+1)}},
	}
	targetDecl.Specs = append(targetDecl.Specs, newSpec)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return fmt.Errorf("format code error: %w", err)
	}

	return os.WriteFile(keyFile, buf.Bytes(), 0644)
}

// AddModuleField thêm import và field vào struct Modules
func AddModuleField(filePath, domain, domainCamel, moduleName string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// 1. Thêm import nếu chưa có
	importPath := fmt.Sprintf("%s/internal/%s", moduleName, domain)
	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(importPath),
		},
	}
	importExists := false
	ast.Inspect(node, func(n ast.Node) bool {
		if imp, ok := n.(*ast.ImportSpec); ok {
			if imp.Path.Value == strconv.Quote(importPath) {
				importExists = true
				return false
			}
		}
		return true
	})
	if !importExists {
		var importDecl *ast.GenDecl
		ast.Inspect(node, func(n ast.Node) bool {
			if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.IMPORT {
				importDecl = decl
				return false
			}
			return true
		})
		if importDecl != nil {
			importDecl.Specs = append(importDecl.Specs, newImport)
		} else {
			newImportDecl := &ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: []ast.Spec{newImport},
			}
			node.Decls = append([]ast.Decl{newImportDecl}, node.Decls...)
		}
	}

	// 2. Tìm struct Modules và thêm field
	var targetStruct *ast.StructType
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name == "Modules" {
				if st, ok := x.Type.(*ast.StructType); ok {
					targetStruct = st
					return false
				}
			}
		}
		return true
	})
	if targetStruct == nil {
		return fmt.Errorf("struct Modules not found in %s", filePath)
	}
	for _, field := range targetStruct.Fields.List {
		if len(field.Names) > 0 && field.Names[0].Name == domainCamel {
			return fmt.Errorf("field %s already exists", domainCamel)
		}
	}
	newField := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(domainCamel)},
		Type: &ast.StarExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent(domain),
				Sel: ast.NewIdent(domainCamel + "Module"),
			},
		},
	}
	targetStruct.Fields.List = append(targetStruct.Fields.List, newField)

	// Ghi lại file
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return err
	}
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// AddModuleInit thêm khởi tạo module vào InitModules và return statement
func AddModuleInit(filePath, domain, domainCamel string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Tìm hàm InitModules
	var targetFunc *ast.FuncDecl
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name.Name == "InitModules" {
				targetFunc = x
				return false
			}
		}
		return true
	})
	if targetFunc == nil {
		return fmt.Errorf("function InitModules not found in %s", filePath)
	}

	// Tìm return statement
	var returnStmt *ast.ReturnStmt
	ast.Inspect(targetFunc.Body, func(n ast.Node) bool {
		if rs, ok := n.(*ast.ReturnStmt); ok {
			returnStmt = rs
			return false
		}
		return true
	})
	if returnStmt == nil {
		return fmt.Errorf("return statement not found in InitModules")
	}

	// Tìm composite literal của Modules trong return statement (xử lý cả &Modules{...})
	var compLit *ast.CompositeLit
	ast.Inspect(returnStmt, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CompositeLit:
			switch t := x.Type.(type) {
			case *ast.Ident:
				if t.Name == "Modules" {
					compLit = x
					return false
				}
			case *ast.SelectorExpr:
				if t.Sel.Name == "Modules" {
					compLit = x
					return false
				}
			}
		case *ast.UnaryExpr:
			if x.Op == token.AND {
				if cl, ok := x.X.(*ast.CompositeLit); ok {
					switch t := cl.Type.(type) {
					case *ast.Ident:
						if t.Name == "Modules" {
							compLit = cl
							return false
						}
					case *ast.SelectorExpr:
						if t.Sel.Name == "Modules" {
							compLit = cl
							return false
						}
					}
				}
			}
		}
		return true
	})
	if compLit == nil {
		return fmt.Errorf("cannot find composite literal for Modules in return statement")
	}

	// Kiểm tra xem field đã tồn tại trong return chưa
	for _, elt := range compLit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == domainCamel {
				return fmt.Errorf("field %s already exists in return statement", domainCamel)
			}
		}
	}

	// Tạo câu lệnh khởi tạo: bookModule := book.NewBookModule(core)
	assignStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent(domain + "Module")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{&ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent(domain),
				Sel: ast.NewIdent("New" + domainCamel + "Module"),
			},
			Args: []ast.Expr{ast.NewIdent("core")},
		}},
	}

	// Chèn assignStmt vào trước returnStmt
	newBody := make([]ast.Stmt, 0, len(targetFunc.Body.List)+1)
	for _, stmt := range targetFunc.Body.List {
		if stmt == returnStmt {
			newBody = append(newBody, assignStmt)
		}
		newBody = append(newBody, stmt)
	}
	targetFunc.Body.List = newBody

	// Thêm field vào composite literal
	newField := &ast.KeyValueExpr{
		Key:   ast.NewIdent(domainCamel),
		Value: ast.NewIdent(domain + "Module"),
	}
	compLit.Elts = append(compLit.Elts, newField)

	// Ghi lại file
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return fmt.Errorf("format code error: %w", err)
	}
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// AddRouteRegistration thêm route registration vào routers.go
func AddRouteRegistration(filePath, domain, domainCamel string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Tìm hàm RegisterRoutes
	var targetFunc *ast.FuncDecl
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name.Name == "RegisterRoutes" {
				targetFunc = x
				return false
			}
		}
		return true
	})
	if targetFunc == nil {
		return fmt.Errorf("function RegisterRoutes not found in %s", filePath)
	}

	// Tìm dòng "modules.Book.RegisterProtectedRoutes" để chèn sau
	// Cách đơn giản: tìm return hoặc dòng cuối của protected routes
	var lastProtectedCall *ast.ExprStmt
	ast.Inspect(targetFunc.Body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ExprStmt:
			// Kiểm tra xem đây có phải là lời gọi RegisterProtectedRoutes không
			if call, ok := x.X.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if sel.Sel.Name == "RegisterProtectedRoutes" {
						lastProtectedCall = x
					}
				}
			}
		}
		return true
	})

	if lastProtectedCall == nil {
		return fmt.Errorf("no RegisterProtectedRoutes call found")
	}

	// Tạo lời gọi mới: modules.Book.RegisterProtectedRoutes(gen, &rateLimit, authMiddleware)
	newCall := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("modules"),
					Sel: ast.NewIdent(domainCamel),
				},
				Sel: ast.NewIdent("RegisterProtectedRoutes"),
			},
			Args: []ast.Expr{
				ast.NewIdent("gen"),
				&ast.UnaryExpr{
					Op: token.AND,
					X:  ast.NewIdent("rateLimit"),
				},
				ast.NewIdent("authMiddleware"),
			},
		},
	}

	// Chèn newCall sau lastProtectedCall
	newBody := make([]ast.Stmt, 0, len(targetFunc.Body.List)+1)
	for _, stmt := range targetFunc.Body.List {
		newBody = append(newBody, stmt)
		if stmt == lastProtectedCall {
			newBody = append(newBody, newCall)
		}
	}
	targetFunc.Body.List = newBody

	// Ghi lại file
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return err
	}
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}
