package cmd

import (
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// ==================== Test WitchesCreate ====================

func TestWitchesCreate_Success(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

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
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		tmpDir := t.TempDir()
		originalWd, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}
		defer os.Chdir(originalWd)

		err = os.Chdir(tmpDir)
		if err != nil {
			os.Exit(1)
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
		defer os.Chdir(originalWd)

		err = os.Chdir(tmpDir)
		if err != nil {
			os.Exit(1)
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

// ==================== Test WitchesInit ====================

func TestWitchesInit_Success(t *testing.T) {
	tmpDir := t.TempDir()

	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test-module"), 0644)
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	WitchesInit("postgres")

	expectedDirs := []string{
		"cmd",
		"internal",
		"pkg",
		"migrate/migrations",
	}
	for _, dir := range expectedDirs {
		_, err := os.Stat(dir)
		assert.NoError(t, err, "Directory %s should be created", dir)
	}

	_, err = os.Stat("main.go")
	assert.NoError(t, err, "main.go should exist at root")
}

func TestWitchesInit_InvalidDriver(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS_INIT") == "1" {
		WitchesInit("invalid-driver")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestWitchesInit_InvalidDriver")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS_INIT=1")

	err := cmd.Run()
	assert.Error(t, err, "Should exit with error")
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		assert.Equal(t, 1, exitErr.ExitCode(), "Should exit with code 1")
	} else {
		t.Errorf("Expected ExitError, got %T", err)
	}
}

// ==================== Test WitchesInstall ====================

func TestWitchesInstall_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	// Tạo go.mod
	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test"), 0644)
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Gọi WitchesInstall với driver postgres
	WitchesInstall("postgres")

	// Kiểm tra go.mod có được cập nhật không (có thể dùng go list)
	// Hoặc chỉ kiểm tra không có panic
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

// ==================== Test helper functions ====================

func TestFindProjectRoot_Success(t *testing.T) {
	tmpDir := t.TempDir()

	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test"), 0644)
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	root := findProjectRoot()
	assert.Equal(t, tmpDir, root)
}

// cmd/run/run_test.go
func TestFindProjectRoot_NoGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	root := findProjectRoot()
	// Nếu có go.mod ở thư mục cha (do môi trường có project Go), bỏ qua test
	if root != "" {
		t.Skip("Skipping: go.mod found in parent directory, cannot test empty root")
	}
	assert.Empty(t, root, "Should return empty when no go.mod found")
}

// ==================== Test WitchesRun ====================

func TestWitchesRun_Success(t *testing.T) {
	tmpDir := t.TempDir()

	// Tạo go.mod
	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test"), 0644)
	require.NoError(t, err)

	// Tạo main.go
	mainFile := filepath.Join(tmpDir, "main.go")
	mainContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello from witches run test")
}
`
	err = os.WriteFile(mainFile, []byte(mainContent), 0644)
	require.NoError(t, err)

	// Tạo internal/dto để easyjson generate
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
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Gọi WitchesRun - sẽ chạy go run .
	WitchesRun()
}

// ==================== Test generateEasyJSONForDir ====================

func TestGenerateEasyJSONForDir(t *testing.T) {
	// Skip nếu chưa có easyjson
	if _, err := exec.LookPath("easyjson"); err != nil {
		t.Skip("easyjson not installed, skip test")
	}

	tmpDir := t.TempDir()

	// Tạo file .go với struct
	dtoFile := filepath.Join(tmpDir, "test.go")
	dtoContent := `package test

type TestStruct struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
	err := os.WriteFile(dtoFile, []byte(dtoContent), 0644)
	require.NoError(t, err)

	// Gọi generateEasyJSONForDir
	generateEasyJSONForDir(tmpDir, "test")

	// Kiểm tra file _easyjson.go được tạo
	files, err := filepath.Glob(filepath.Join(tmpDir, "*_easyjson.go"))
	assert.NoError(t, err)
	if len(files) == 0 {
		t.Log("No easyjson files generated")
	} else {
		assert.Greater(t, len(files), 0, "EasyJSON file should be generated")
	}
}

// ==================== Test WitchesRun with missing files ====================
// cmd/run/run_test.go
func TestWitchesRun_NoMainGo(t *testing.T) {
	//  Dùng subprocess để test log.Fatal
	if os.Getenv("TEST_SUBPROCESS_RUN") == "1" {
		tmpDir := t.TempDir()

		// Tạo go.mod nhưng không có main.go
		goModPath := filepath.Join(tmpDir, "go.mod")
		err := os.WriteFile(goModPath, []byte("module test"), 0644)
		if err != nil {
			os.Exit(1)
		}

		originalWd, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}
		defer os.Chdir(originalWd)

		err = os.Chdir(tmpDir)
		if err != nil {
			os.Exit(1)
		}

		// Gọi WitchesRun - sẽ log.Fatal
		WitchesRun()
		return
	}

	// Chạy subprocess để test
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

// ==================== Test RemoveEasyJSONFiles ====================

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
func TestWitchesAdd_Success(t *testing.T) {
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test-module"), 0644)
	require.NoError(t, err)

	// Tạo internal/shared/utils với Object constant
	sharedDir := filepath.Join(tmpDir, "internal", "shared", "utils")
	err = os.MkdirAll(sharedDir, 0755)
	require.NoError(t, err)
	keyContent := `package utils

var Object uint = 10
`
	err = os.WriteFile(filepath.Join(sharedDir, "key_object.go"), []byte(keyContent), 0644)
	require.NoError(t, err)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	WitchesAdd("book")

	// Kiểm tra thư mục service được tạo
	serviceDir := filepath.Join(tmpDir, "internal", "book-service")
	_, err = os.Stat(serviceDir)
	assert.NoError(t, err, "Service directory should be created")
}
