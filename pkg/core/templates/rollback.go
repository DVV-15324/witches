package template

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func RollbackDomain(project string, moduleName string, domainName string) error {
	domainName = strings.TrimSpace(domainName)
	domainName = strings.ReplaceAll(domainName, " ", "")
	domainName = strings.ToLower(domainName)
	domainNameCap := cases.Title(language.English).String(domainName)
	domainNameCap = strings.ReplaceAll(domainNameCap, " ", "")

	fmt.Printf("Rolling back domain '%s' ...\n", domainName)

	domainDir := filepath.Join(project, "internal", domainName)
	if err := os.RemoveAll(domainDir); err != nil {
		return fmt.Errorf("failed to remove domain directory: %v", err)
	}
	fmt.Printf("Removed directory: %s\n", domainDir)

	sharedFile := filepath.Join(project, "internal", "shared", "domain", domainName+".go")
	if err := os.Remove(sharedFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove shared domain file: %v", err)
	}
	if _, err := os.Stat(sharedFile); err == nil {
		fmt.Printf("Removed file: %s\n", sharedFile)
	}

	keyFile := filepath.Join(project, "internal", "shared", "utils", "key_object.go")
	if err := rollbackKeyObject(keyFile, domainNameCap); err != nil {
		fmt.Printf("Warning: failed to rollback key_object.go: %v\n", err)
	} else {
		fmt.Println("Rolled back key_object.go")
	}

	modulesPath := filepath.Join(project, "cmd", "server", "routers", "modules.go")
	if err := rollbackModuleField(modulesPath, domainName, domainNameCap, moduleName); err != nil {
		fmt.Printf("Warning: failed to rollback modules.go (field): %v\n", err)
	} else {
		fmt.Println("Rolled back modules.go: removed import and field")
	}

	if err := rollbackModuleInit(modulesPath, domainName, domainNameCap); err != nil {
		fmt.Printf("Warning: failed to rollback modules.go (init): %v\n", err)
	} else {
		fmt.Println("Rolled back modules.go: removed initialization")
	}

	routersPath := filepath.Join(project, "cmd", "server", "routers", "routers.go")
	if err := rollbackRouteRegistration(routersPath, domainName, domainNameCap); err != nil {
		fmt.Printf("Warning: failed to rollback routers.go: %v\n", err)
	} else {
		fmt.Println("Rolled back routers.go: removed route registration")
	}

	fmt.Printf("Domain '%s' rolled back successfully!\n", domainName)
	return nil
}

func rollbackKeyObject(filePath, domainCap string) error {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file error: %w", err)
	}

	var targetDecl *ast.GenDecl
	var targetIndex int

	ast.Inspect(node, func(n ast.Node) bool {
		decl, ok := n.(*ast.GenDecl)
		if !ok || decl.Tok != token.VAR {
			return true
		}

		for idx, spec := range decl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 {
				continue
			}
			if vs.Names[0].Name == "Object"+domainCap {
				targetDecl = decl
				targetIndex = idx
				return false
			}
		}
		return true
	})

	if targetDecl == nil {
		return fmt.Errorf("Object%s not found", domainCap)
	}

	targetDecl.Specs = append(targetDecl.Specs[:targetIndex], targetDecl.Specs[targetIndex+1:]...)

	if len(targetDecl.Specs) == 0 {
		var newDecls []ast.Decl
		for _, decl := range node.Decls {
			if decl != targetDecl {
				newDecls = append(newDecls, decl)
			}
		}
		node.Decls = newDecls
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return fmt.Errorf("format code error: %w", err)
	}

	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func rollbackModuleField(filePath, domain, domainCap, moduleName string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	importPath := fmt.Sprintf("%s/internal/%s", moduleName, domain)

	var importDecl *ast.GenDecl
	ast.Inspect(node, func(n ast.Node) bool {
		if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.IMPORT {
			importDecl = decl
			return false
		}
		return true
	})

	if importDecl != nil {
		var newSpecs []ast.Spec
		for _, spec := range importDecl.Specs {
			if imp, ok := spec.(*ast.ImportSpec); ok {
				if imp.Path.Value == strconv.Quote(importPath) {
					continue
				}
			}
			newSpecs = append(newSpecs, spec)
		}
		importDecl.Specs = newSpecs

		if len(importDecl.Specs) == 0 {
			var newDecls []ast.Decl
			for _, decl := range node.Decls {
				if decl != importDecl {
					newDecls = append(newDecls, decl)
				}
			}
			node.Decls = newDecls
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

	var newFields []*ast.Field
	for _, field := range targetStruct.Fields.List {
		if len(field.Names) > 0 && field.Names[0].Name == domainCap {
			continue
		}
		newFields = append(newFields, field)
	}
	targetStruct.Fields.List = newFields

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return err
	}
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func rollbackModuleInit(filePath, domain, domainCap string) error {
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

	var newBody []ast.Stmt
	for _, stmt := range targetFunc.Body.List {
		if assign, ok := stmt.(*ast.AssignStmt); ok {
			if len(assign.Lhs) > 0 {
				if ident, ok := assign.Lhs[0].(*ast.Ident); ok && ident.Name == domain+"Module" {
					continue
				}
			}
		}
		newBody = append(newBody, stmt)
	}
	targetFunc.Body.List = newBody

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

	if compLit != nil {
		var newElts []ast.Expr
		for _, elt := range compLit.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == domainCap {
					continue
				}
			}
			newElts = append(newElts, elt)
		}
		compLit.Elts = newElts
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return fmt.Errorf("format code error: %w", err)
	}
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func rollbackRouteRegistration(filePath, domain, domainCap string) error {
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

	var newBody []ast.Stmt
	for _, stmt := range targetFunc.Body.List {
		shouldRemove := false

		if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
			if call, ok := exprStmt.X.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {

					if sel.Sel.Name == "AddTag" && len(call.Args) > 0 {
						if lit, ok := call.Args[0].(*ast.BasicLit); ok {
							if lit.Value == strconv.Quote(domain) {
								shouldRemove = true
							}
						}
					}

					if sel.Sel.Name == "RegisterPublicRoutes" || sel.Sel.Name == "RegisterProtectedRoutes" {
						if innerSel, ok := sel.X.(*ast.SelectorExpr); ok {
							if innerSel.Sel.Name == domainCap {
								shouldRemove = true
							}
						}
					}
				}
			}
		}

		if !shouldRemove {
			newBody = append(newBody, stmt)
		}
	}
	targetFunc.Body.List = newBody

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return err
	}
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}
