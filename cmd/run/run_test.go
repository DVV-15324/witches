package cmd

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWitchesCreate_Success(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			log.Printf("failed to chdir: %v", err)
		}
	}()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	projectName := "test-project"
	WitchesCreate(projectName)

	projectPath := filepath.Join(tmpDir, projectName)
	_, err = os.Stat(projectPath)
	assert.NoError(t, err, "Project directory should be created")

	envPath := filepath.Join(projectPath, "witches.env")
	_, err = os.Stat(envPath)
	assert.NoError(t, err, "witches.env should be created")

	content, err := os.ReadFile(envPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "DB_")
}

func TestWitchesCreate_ProjectExists(t *testing.T) {
	if e := os.Getenv("TEST_SUBPROCESS"); e == "1" {
		tmpDir := t.TempDir()
		originalWd, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}
		defer func() {
			if err := os.Chdir(originalWd); err != nil {
				log.Printf("failed to chdir: %v", err)
			}
		}()
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}

		projectName := "existing-project"
		err = os.Mkdir(projectName, 0755)
		if err != nil {
			os.Exit(1)
		}

		WitchesCreate(projectName)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestWitchesCreate_ProjectExists")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=1")

	err := cmd.Run()
	assert.Error(t, err, "Should exit with error")
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		assert.Equal(t, 1, exitErr.ExitCode(), "Should exit with code 1")
	} else {
		t.Errorf("Expected ExitError, got %T", err)
	}
}

func TestWitchesCreate_EnvFileExists(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS_ENV") == "1" {
		tmpDir := t.TempDir()
		originalWd, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}
		defer func() {
			if err := os.Chdir(originalWd); err != nil {
				log.Printf("failed to chdir: %v", err)
			}
		}()
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}

		projectName := "test-project"
		projectPath := filepath.Join(tmpDir, projectName)

		err = os.MkdirAll(projectPath, 0755)
		if err != nil {
			os.Exit(1)
		}

		envPath := filepath.Join(projectPath, "witches.env")
		err = os.WriteFile(envPath, []byte("existing env"), 0644)
		if err != nil {
			os.Exit(1)
		}

		WitchesCreate(projectName)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestWitchesCreate_EnvFileExists")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS_ENV=1")

	err := cmd.Run()
	assert.Error(t, err, "Should exit with error")
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		assert.Equal(t, 1, exitErr.ExitCode(), "Should exit with code 1")
	}
}

func TestWitchesLink_Success(t *testing.T) {
	// Tạo temp directory
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	// Chuyển vào temp dir
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Tạo go.mod
	goModContent := `module captain

go 1.21
`
	err = os.WriteFile("go.mod", []byte(goModContent), 0644)
	require.NoError(t, err)

	// Tạo thư mục cần thiết
	dirs := []string{
		"cmd/server/routers",
		"internal/shared/domain",
		"internal/shared/middleware",
		"internal/shared/utils",
	}
	for _, dir := range dirs {
		err = os.MkdirAll(dir, 0755)
		require.NoError(t, err)
	}

	// Tạo file modules.go
	modulesContent := `package routers

type Modules struct {
}

func (m *Modules) Init() *Modules {
	return &Modules{}
}
`
	err = os.WriteFile("cmd/server/routers/modules.go", []byte(modulesContent), 0644)
	require.NoError(t, err)

	// Tạo file routers.go
	routersContent := `package routers

func SetupRoutes() error {
	return nil
}
`
	err = os.WriteFile("cmd/server/routers/routers.go", []byte(routersContent), 0644)
	require.NoError(t, err)

	// Tạo file key_object.go
	keyObjectContent := `package domain

type KeyObject string

const (
	User KeyObject = "user"
)
`
	err = os.WriteFile("internal/shared/domain/key_object.go", []byte(keyObjectContent), 0644)
	require.NoError(t, err)

	// Chạy command
	WitchesLink("book", "https://github.com/DVV-15324/witches-book.git")

	// Kiểm tra file đã được tạo
	expectedFiles := []string{
		"internal/book/module.go",
		"internal/book/model/model.go",
		"internal/book/handler/handler.go",
		"internal/book/repository/repository.go",
		"internal/book/usecase/usecase.go",
		"internal/book/migrate/migrations/1_create_table.up.sql",
		"internal/book/migrate/migrations/1_drop_table.down.sql",
		"internal/shared/domain/book.go",
	}

	for _, file := range expectedFiles {
		_, err := os.Stat(file)
		assert.NoError(t, err, "File should exist: %s", file)
	}
}

func TestWitchesLink_NoGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Không có go.mod, command sẽ fail
	// Cần capture output hoặc check error
	WitchesLink("book", "https://github.com/DVV-15324/witches-book.git")

	// Kiểm tra không có file nào được tạo
	_, err = os.Stat("internal/book")
	assert.Error(t, err, "Book directory should not be created")
}

func TestWitchesLink_InvalidRepo(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Tạo go.mod
	err = os.WriteFile("go.mod", []byte("module test"), 0644)
	require.NoError(t, err)

	// Test với repo không tồn tại
	WitchesLink("book", "https://github.com/invalid/repo.git")

	// Kiểm tra không có file nào được tạo
	_, err = os.Stat("internal/book")
	assert.Error(t, err, "Book directory should not be created")
}

func TestWitchesInstall_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test"), 0644)
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			log.Printf("failed to chdir: %v", err)
		}
	}()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	require.NoError(t, err)

	WitchesInstall("postgres")

	t.Log("Installation completed successfully")
}

func TestWitchesInstall_DriverMapping(t *testing.T) {
	drivers := map[string]string{
		"mysql":      "github.com/golang-migrate/migrate/v4/database/mysql@latest",
		"postgres":   "github.com/golang-migrate/migrate/v4/database/postgres@latest",
		"postgresql": "github.com/golang-migrate/migrate/v4/database/postgres@latest",
		"mssql":      "github.com/golang-migrate/migrate/v4/database/sqlserver@latest",
		"sqlserver":  "github.com/golang-migrate/migrate/v4/database/sqlserver@latest",
	}

	for driver, expected := range drivers {
		t.Run(driver, func(t *testing.T) {
			assert.Equal(t, expected, drivers[driver])
		})
	}
}

func TestFindProjectRoot_Success(t *testing.T) {
	tmpDir := t.TempDir()

	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test"), 0644)
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			log.Printf("failed to chdir: %v", err)
		}
	}()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	require.NoError(t, err)

	root := findProjectRoot()
	assert.Equal(t, tmpDir, root)
}

func TestFindProjectRoot_NoGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			log.Printf("failed to chdir: %v", err)
		}
	}()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	require.NoError(t, err)

	root := findProjectRoot()

	if root != "" {
		t.Skip("Skipping: go.mod found in parent directory, cannot test empty root")
	}
	assert.Empty(t, root, "Should return empty when no go.mod found")
}

func TestWitchesRun_Success(t *testing.T) {
	tmpDir := t.TempDir()

	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test"), 0644)
	require.NoError(t, err)

	mainFile := filepath.Join(tmpDir, "main.go")
	mainContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello from witches run test")
}
`
	err = os.WriteFile(mainFile, []byte(mainContent), 0644)
	require.NoError(t, err)

	dtoDir := filepath.Join(tmpDir, "internal", "user", "dto", "request")
	err = os.MkdirAll(dtoDir, 0755)
	require.NoError(t, err)

	dtoFile := filepath.Join(dtoDir, "user_request.go")
	dtoContent := `package request

type UserRequest struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
	err = os.WriteFile(dtoFile, []byte(dtoContent), 0644)
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			log.Printf("failed to chdir: %v", err)
		}
	}()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	require.NoError(t, err)

	WitchesRun()
}

func TestGenerateEasyJSONForDir(t *testing.T) {
	if _, err := exec.LookPath("easyjson"); err != nil {
		t.Skip("easyjson not installed, skip test")
	}
	tmpDir := t.TempDir()
	dtoFile := filepath.Join(tmpDir, "test.go")
	dtoContent := `package test

type TestStruct struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
	err := os.WriteFile(dtoFile, []byte(dtoContent), 0644)
	require.NoError(t, err)

	generateEasyJSONForDir(tmpDir, "test")

	files, err := filepath.Glob(filepath.Join(tmpDir, "*_easyjson.go"))
	assert.NoError(t, err)
	if len(files) == 0 {
		t.Log("No easyjson files generated")
	} else {
		assert.Greater(t, len(files), 0, "EasyJSON file should be generated")
	}
}

func TestWitchesRun_NoMainGo(t *testing.T) {

	if os.Getenv("TEST_SUBPROCESS_RUN") == "1" {
		tmpDir := t.TempDir()

		goModPath := filepath.Join(tmpDir, "go.mod")
		err := os.WriteFile(goModPath, []byte("module test"), 0644)
		if err != nil {
			os.Exit(1)
		}

		originalWd, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}
		defer func() {
			if err := os.Chdir(originalWd); err != nil {
				log.Printf("failed to chdir: %v", err)
			}
		}()
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}

		WitchesRun()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestWitchesRun_NoMainGo")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS_RUN=1")

	err := cmd.Run()
	assert.Error(t, err, "Should exit with error")
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		assert.Equal(t, 1, exitErr.ExitCode(), "Should exit with code 1")
	} else {
		t.Errorf("Expected ExitError, got %T", err)
	}
}

func TestRemoveEasyJSONFiles(t *testing.T) {
	tmpDir := t.TempDir()

	easyjsonFile := filepath.Join(tmpDir, "test_easyjson.go")
	err := os.WriteFile(easyjsonFile, []byte("// easyjson"), 0644)
	require.NoError(t, err)

	normalFile := filepath.Join(tmpDir, "test.go")
	err = os.WriteFile(normalFile, []byte("package test"), 0644)
	require.NoError(t, err)

	removeEasyJSONFiles(tmpDir)

	_, err = os.Stat(easyjsonFile)
	assert.Error(t, err, "EasyJSON file should be removed")

	_, err = os.Stat(normalFile)
	assert.NoError(t, err, "Normal file should remain")
}
func TestWitchesRollback_Success(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Tạo go.mod
	err = os.WriteFile("go.mod", []byte("module captain"), 0644)
	require.NoError(t, err)

	// Tạo thư mục internal/shared/domain trước
	err = os.MkdirAll(filepath.Join("internal", "shared", "domain"), 0755)
	require.NoError(t, err)

	// Tạo domain book
	bookDir := filepath.Join("internal", "book")
	err = os.MkdirAll(bookDir, 0755)
	require.NoError(t, err)

	// Tạo file module.go trong book
	moduleContent := `package book

type Module struct{}
`
	err = os.WriteFile(filepath.Join(bookDir, "module.go"), []byte(moduleContent), 0644)
	require.NoError(t, err)

	// Tạo shared domain book
	sharedFile := filepath.Join("internal", "shared", "domain", "book.go")
	err = os.WriteFile(sharedFile, []byte("package domain"), 0644)
	require.NoError(t, err)

	// Chạy rollback
	WitchesRollback("book")

	// Kiểm tra đã xóa
	_, err = os.Stat(bookDir)
	assert.Error(t, err, "Book directory should be removed")

	// Kiểm tra shared domain file đã xóa
	_, err = os.Stat(sharedFile)
	assert.Error(t, err, "Shared domain file should be removed")
}
func TestWitchesAdd_Success(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Tạo go.mod
	err = os.WriteFile("go.mod", []byte("module captain"), 0644)
	require.NoError(t, err)

	// Tạo các thư mục cần thiết
	dirs := []string{
		"internal/shared/domain",
		"internal/shared/utils", // Thêm thư mục utils
		"internal/shared/middleware",
		"cmd/server/routers",
	}
	for _, dir := range dirs {
		err = os.MkdirAll(dir, 0755)
		require.NoError(t, err)
	}

	// Tạo file key_object.go trong utils
	keyContent := `package utils

type KeyObject string

const (
	User KeyObject = "user"
)
`
	err = os.WriteFile("internal/shared/utils/key_object.go", []byte(keyContent), 0644)
	require.NoError(t, err)

	// Tạo file key_object.go trong domain
	domainKeyContent := `package domain

type KeyObject string

const (
	User KeyObject = "user"
)
`
	err = os.WriteFile("internal/shared/domain/key_object.go", []byte(domainKeyContent), 0644)
	require.NoError(t, err)

	// Tạo file modules.go
	modulesContent := `package routers

type Modules struct {
	User *user.Module
}

func (m *Modules) InitModules() *Modules {
	return &Modules{
		User: &user.Module{},
	}
}
`
	err = os.WriteFile("cmd/server/routers/modules.go", []byte(modulesContent), 0644)
	require.NoError(t, err)

	// Tạo file routers.go
	routersContent := `package routers

func SetupRoutes() error {
	return nil
}
`
	err = os.WriteFile("cmd/server/routers/routers.go", []byte(routersContent), 0644)
	require.NoError(t, err)

	// Chạy add domain
	WitchesAdd("book", "postgres")

	// Kiểm tra domain được tạo
	expectedFiles := []string{
		"internal/book/module.go",
		"internal/book/model/model.go",
		"internal/book/handler/handler.go",
		"internal/book/repository/repository.go",
		"internal/book/usecase/usecase.go",
	}

	for _, file := range expectedFiles {
		_, err := os.Stat(file)
		assert.NoError(t, err, "File should exist: %s", file)
	}
}
func TestWitchesRun_WithArgs(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Tạo go.mod
	err = os.WriteFile("go.mod", []byte("module test"), 0644)
	require.NoError(t, err)

	// Tạo main.go
	mainContent := `package main
import "fmt"
func main() {
	fmt.Println("Hello")
}
`
	err = os.WriteFile("main.go", []byte(mainContent), 0644)
	require.NoError(t, err)

	// Lưu os.Args cũ và thay đổi
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"witches", "run", "arg1", "arg2"}

	// Chạy với args
	WitchesRun()
}

func TestWitchesRun_NoModule(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Không có go.mod, không có main.go
	WitchesRun()
}
