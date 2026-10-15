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

var (
	removeAll = os.RemoveAll
	remove    = os.Remove
	stat      = os.Stat
)

func RemoveModule(project string, projectName string, moduleName string) error {
	moduleName = strings.TrimSpace(moduleName)
	moduleName = strings.ReplaceAll(moduleName, " ", "")
	moduleName = strings.ToLower(moduleName)
	moduleNameCap := cases.Title(language.English).String(moduleName)
	moduleNameCap = strings.ReplaceAll(moduleNameCap, " ", "")

	fmt.Printf("Removing module '%s' ...\n", moduleName)

	moduleDir := filepath.Join(project, "internal", moduleName)
	if err := removeAll(moduleDir); err != nil {
		return fmt.Errorf("failed to remove module directory: %v", err)
	}
	fmt.Printf("Removed directory: %s\n", moduleDir)

	sharedFile := filepath.Join(project, "internal", "shared", "domain", moduleName+".go")
	if err := remove(sharedFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove shared domain file: %v", err)
	}
	if _, err := stat(sharedFile); err == nil {
		fmt.Printf("Removed file: %s\n", sharedFile)
	}

	keyFile := filepath.Join(project, "internal", "shared", "utils", "key_object.go")
	if err := removeKeyObject(keyFile, moduleNameCap); err != nil {
		fmt.Printf("Warning: failed to remove key_object.go: %v\n", err)
	} else {
		fmt.Println("Removed key_object.go")
	}

	modulesPath := filepath.Join(project, "cmd", "server", "routers", "modules.go")
	if err := removeModuleField(modulesPath, moduleName, moduleNameCap, projectName); err != nil {
		fmt.Printf("Warning: failed to remove modules.go (field): %v\n", err)
	} else {
		fmt.Println("Removed modules.go: removed import and field")
	}

	if err := removeModuleInit(modulesPath, moduleName, moduleNameCap); err != nil {
		fmt.Printf("Warning: failed to remove modules.go (init): %v\n", err)
	} else {
		fmt.Println("Removed modules.go: removed initialization")
	}

	routersPath := filepath.Join(project, "cmd", "server", "routers", "routers.go")
	if err := removeRouteRegistration(routersPath, moduleName, moduleNameCap); err != nil {
		fmt.Printf("Warning: failed to remove routers.go: %v\n", err)
	} else {
		fmt.Println("Removed routers.go: removed route registration")
	}

	fmt.Printf("Module '%s' removed successfully!\n", moduleName)
	return nil
}

func removeKeyObject(filePath, moduleCap string) error {
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
			if vs.Names[0].Name == "Object"+moduleCap {
				targetDecl = decl
				targetIndex = idx
				return false
			}
		}
		return true
	})

	if targetDecl == nil {
		return fmt.Errorf("Object%s not found", moduleCap)
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

func removeModuleField(filePath, module, moduleCap, projectName string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	importPath := fmt.Sprintf("%s/internal/%s", projectName, module)

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
		if len(field.Names) > 0 && field.Names[0].Name == moduleCap {
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

func removeModuleInit(filePath, module, moduleCap string) error {
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
				if ident, ok := assign.Lhs[0].(*ast.Ident); ok && ident.Name == module+"Module" {
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
				if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == moduleCap {
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

func removeRouteRegistration(filePath, module, moduleCap string) error {
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
							if lit.Value == strconv.Quote(module) {
								shouldRemove = true
							}
						}
					}

					if sel.Sel.Name == "RegisterPublicRoutes" || sel.Sel.Name == "RegisterProtectedRoutes" {
						if innerSel, ok := sel.X.(*ast.SelectorExpr); ok {
							if innerSel.Sel.Name == moduleCap {
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
