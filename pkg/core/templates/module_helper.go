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
)

func generateSharedDomain(projectPath string, config ModuleConfig, templateFs embed.FS) error {
	sharedDomainDir := filepath.Join(projectPath, "internal", "shared", "domain")
	if err := os.MkdirAll(sharedDomainDir, 0755); err != nil {
		return err
	}
	destFile := filepath.Join(sharedDomainDir, config.Name+".go")
	tmplFile := "module/shared/domain/domain.go.tmpl"
	tmplContent, err := templateFs.ReadFile(tmplFile)
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
	defer func() {
		_ = file.Close()
	}()
	return tmpl.Execute(file, config)
}

func updateKeyObject(projectPath string, config ModuleConfig) error {
	keyFile := filepath.Join(projectPath, "internal", "shared", "utils", "key_object.go")
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
		if !ok || decl.Tok != token.CONST {
			return true
		}
		for _, spec := range decl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 {
				continue
			}
			//  Kiểm tra tất cả các constant bắt đầu bằng "Object"
			if strings.HasPrefix(vs.Names[0].Name, "Object") {
				if lit, ok := vs.Values[0].(*ast.BasicLit); ok && lit.Kind == token.INT {
					id, err := strconv.Atoi(lit.Value)
					if err != nil {
						continue
					}
					if id > maxID || targetDecl == nil {
						maxID = id
						targetDecl = decl
					}
				}
			}
		}
		return true
	})
	if targetDecl == nil {
		return fmt.Errorf("no Object constant found")
	}

	//  Tạo constant mới với giá trị tăng thêm 1
	newSpec := &ast.ValueSpec{
		Names:  []*ast.Ident{ast.NewIdent("Object" + config.NameCap)},
		Type:   &ast.Ident{Name: "KeyObject"},
		Values: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", maxID+1)}},
	}
	targetDecl.Specs = append(targetDecl.Specs, newSpec)
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return fmt.Errorf("format code error: %w", err)
	}
	return os.WriteFile(keyFile, buf.Bytes(), 0644)
}

func AddModuleField(filePath, module, moduleCamel, projectName string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	importPath := fmt.Sprintf("%s/internal/%s", projectName, module)
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
		if len(field.Names) > 0 && field.Names[0].Name == moduleCamel {
			return fmt.Errorf("field %s already exists", moduleCamel)
		}
	}
	newField := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(moduleCamel)},
		Type: &ast.StarExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent(module),
				Sel: ast.NewIdent(moduleCamel + "Module"),
			},
		},
	}
	targetStruct.Fields.List = append(targetStruct.Fields.List, newField)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return err
	}
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func AddModuleInit(filePath, module, moduleCamel string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

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

	for _, elt := range compLit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == moduleCamel {
				return fmt.Errorf("field %s already exists in return statement", moduleCamel)
			}
		}
	}

	assignStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent(module + "Module")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{&ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent(module),
				Sel: ast.NewIdent("New" + moduleCamel + "Module"),
			},
			Args: []ast.Expr{ast.NewIdent("core")},
		}},
	}

	newBody := make([]ast.Stmt, 0, len(targetFunc.Body.List)+1)
	for _, stmt := range targetFunc.Body.List {
		if stmt == returnStmt {
			newBody = append(newBody, assignStmt)
		}
		newBody = append(newBody, stmt)
	}
	targetFunc.Body.List = newBody

	newField := &ast.KeyValueExpr{
		Key:   ast.NewIdent(moduleCamel),
		Value: ast.NewIdent(module + "Module"),
	}
	compLit.Elts = append(compLit.Elts, newField)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return fmt.Errorf("format code error: %w", err)
	}
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func AddRouteRegistration(filePath, domain, domainCamel string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var targetFunc *ast.FuncDecl
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name.Name == "initModule" {
				targetFunc = x
				return false
			}
		}
		return true
	})
	if targetFunc == nil {
		return fmt.Errorf("function initModule not found in %s", filePath)
	}

	// Đảm bảo body không nil và không rỗng
	if targetFunc.Body == nil {
		targetFunc.Body = &ast.BlockStmt{
			List: []ast.Stmt{},
		}
	}

	var exists bool
	ast.Inspect(targetFunc.Body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if innerSel, ok := sel.X.(*ast.SelectorExpr); ok {
					if ident, ok := innerSel.X.(*ast.Ident); ok && ident.Name == "modules" {
						if innerSel.Sel.Name == domainCamel {
							exists = true
							return false
						}
					}
				}
			}
		}
		return true
	})
	if exists {
		return fmt.Errorf("domain %s already registered in initModule", domainCamel)
	}

	addTagStmt := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("gen"),
				Sel: ast.NewIdent("AddTag"),
			},
			Args: []ast.Expr{
				&ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(domain)},
				&ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(domainCamel + " endpoints")},
			},
		},
	}

	publicStmt := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("modules"),
					Sel: ast.NewIdent(domainCamel),
				},
				Sel: ast.NewIdent("RegisterPublicRoutes"),
			},
			Args: []ast.Expr{
				ast.NewIdent("gen"),
			},
		},
	}

	protectedStmt := &ast.ExprStmt{
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
			},
		},
	}

	//  Nếu body rỗng, thêm trực tiếp, nếu không thì chèn vào cuối
	if len(targetFunc.Body.List) == 0 {
		targetFunc.Body.List = []ast.Stmt{addTagStmt, publicStmt, protectedStmt}
	} else {
		// Chèn vào cuối body
		lastStmt := targetFunc.Body.List[len(targetFunc.Body.List)-1]
		newBody := make([]ast.Stmt, 0, len(targetFunc.Body.List)+3)
		for _, stmt := range targetFunc.Body.List {
			newBody = append(newBody, stmt)
			if stmt == lastStmt {
				newBody = append(newBody, addTagStmt)
				newBody = append(newBody, publicStmt)
				newBody = append(newBody, protectedStmt)
			}
		}
		targetFunc.Body.List = newBody
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return err
	}
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}
