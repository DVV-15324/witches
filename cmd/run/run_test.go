package cmd

import (
	"errors"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWitchesCreate_Success(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	err := WitchesCreate("test-project")
	assert.NoError(t, err)

	_, err = os.Stat("test-project/witches.env")
	assert.NoError(t, err)
}

func TestWitchesCreate_ProjectExists(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.Mkdir("existing", 0755)
	err := WitchesCreate("existing")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}
func TestWitchesCreate_OpenFileFail(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	assert.NoError(t, err)

	defer os.Chdir(origWd)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	originalOpenFile := openFileWitchesCreate
	defer func() {
		openFileWitchesCreate = originalOpenFile
	}()

	openFileWitchesCreate = func(
		name string,
		flag int,
		perm os.FileMode,
	) (*os.File, error) {
		return nil, errors.New("mock open file error")
	}

	err = WitchesCreate("test-project")

	assert.Error(t, err)
	assert.EqualError(
		t,
		err,
		"create witches.env: mock open file error",
	)
}
func TestWitchesCreate_MkdirAllFail(t *testing.T) {
	originalMkdirAll := mkdirAllWitchesCreate
	defer func() {
		mkdirAllWitchesCreate = originalMkdirAll
	}()

	mkdirAllWitchesCreate = func(path string, perm os.FileMode) error {
		return errors.New("mock mkdir error")
	}

	err := WitchesCreate("test-project")

	assert.Error(t, err)
	assert.EqualError(
		t,
		err,
		"create project directory: mock mkdir error",
	)
}

func TestWitchesInit_Success(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)
	os.WriteFile("witches.env", []byte("DB_DRIVER=postgres"), 0644)

	err := WitchesInit("postgres")
	assert.NoError(t, err)

	_, err = os.Stat("main.go")
	assert.NoError(t, err)
}

func TestWitchesInit_NoEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)

	err := WitchesInit("postgres")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "witches.env not found")
}

func TestWitchesInit_EnvExists(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(origWd)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	err = os.WriteFile("witches.env", []byte("test"), 0644)
	assert.NoError(t, err)

	err = WitchesInit("postgres")

	assert.NoError(t, err)
}

func TestWitchesInit_GetWorkingDirectoryFail(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(origWd)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	err = os.WriteFile(
		"witches.env",
		[]byte("DB_DRIVER=postgres"),
		0644,
	)
	assert.NoError(t, err)

	originalGetwd := getwdWitchesInit
	defer func() {
		getwdWitchesInit = originalGetwd
	}()

	getwdWitchesInit = func() (string, error) {
		return "", errors.New("mock getwd error")
	}

	err = WitchesInit("postgres")

	assert.Error(t, err)
	assert.EqualError(
		t,
		err,
		"get working directory: mock getwd error",
	)
}

func TestWitchesInit_CreateTemplateFail(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(origWd)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	err = os.WriteFile(
		"witches.env",
		[]byte("DB_DRIVER=postgres"),
		0644,
	)
	assert.NoError(t, err)

	originalInit := initWitchesInit
	defer func() {
		initWitchesInit = originalInit
	}()

	initWitchesInit = func(
		moduleName string,
		DBdriver string,
	) error {
		return errors.New("mock create template error")
	}

	err = WitchesInit("postgres")

	assert.Error(t, err)
	assert.EqualError(
		t,
		err,
		"create template: mock create template error",
	)
}

func TestWitchesAdd_Success(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)
	os.MkdirAll("cmd/server/routers", 0755)
	os.MkdirAll("internal/shared/domain", 0755)
	os.MkdirAll("internal/shared/utils", 0755)

	err := WitchesAdd("user", "postgres")
	assert.NoError(t, err)

	_, err = os.Stat("internal/user/module.go")
	assert.NoError(t, err)
}

func TestWitchesAdd_NoGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	err := WitchesAdd("user", "postgres")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "go.mod not found")
}

func TestWitchesAdd_DomainExists(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)
	os.MkdirAll("internal/user", 0755)

	err := WitchesAdd("user", "postgres")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}
func TestWitchesAdd_GetWorkingDirectoryFail(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	assert.NoError(t, err)

	defer os.Chdir(origWd)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	err = os.WriteFile(
		filepath.Join(tmpDir, "go.mod"),
		[]byte("module test"),
		0644,
	)
	assert.NoError(t, err)

	// Lưu function thật
	originalGetwd := getwdWitchesAdd

	// Khôi phục sau test
	defer func() {
		getwdWitchesAdd = originalGetwd
	}()

	// Mock function
	getwdWitchesAdd = func() (string, error) {
		return "", errors.New("mock getwd error")
	}

	err = WitchesAdd("user", "postgres")

	assert.Error(t, err)
	assert.EqualError(
		t,
		err,
		"get working directory: mock getwd error",
	)
}

func TestWitchesAdd_AddModuleFail(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	assert.NoError(t, err)

	defer os.Chdir(origWd)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	err = os.WriteFile(
		filepath.Join(tmpDir, "go.mod"),
		[]byte("module test"),
		0644,
	)
	assert.NoError(t, err)

	originalAddModule := addModule
	defer func() {
		addModule = originalAddModule
	}()

	addModule = func(
		projectPath string,
		projectName string,
		moduleName string,
		DBdriver string,
	) error {
		return errors.New("mock add module error")
	}

	err = WitchesAdd("user", "postgres")

	assert.Error(t, err)
	assert.EqualError(
		t,
		err,
		"add module: mock add module error",
	)
}

func TestWitchesRollback_Success(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)
	os.MkdirAll("internal/book", 0755)
	os.MkdirAll("cmd/server/routers", 0755)

	err := WitchesRollback("book")
	assert.NoError(t, err)

	_, err = os.Stat("internal/book")
	assert.Error(t, err)
}

func TestWitchesRollback_NoGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	err := WitchesRollback("book")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "go.mod not found")
}

func TestWitchesRollback_DomainNotExists(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)

	err := WitchesRollback("book")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestWitchesLink_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module captain"), 0644)
	os.MkdirAll("cmd/server/routers", 0755)
	os.MkdirAll("internal/shared/domain", 0755)
	os.MkdirAll("internal/shared/utils", 0755)

	os.WriteFile("cmd/server/routers/modules.go", []byte(`package routers
type Modules struct{}
func (m *Modules) Init() *Modules { return &Modules{} }`), 0644)

	os.WriteFile("cmd/server/routers/routers.go", []byte(`package routers
func SetupRoutes() error { return nil }`), 0644)

	os.WriteFile("internal/shared/utils/key_object.go", []byte(`package utils
var (ObjectDefault int64 = 0)`), 0644)

	err := WitchesLink("book", "https://github.com/DVV-15324/witches-book.git")
	assert.NoError(t, err)

	_, err = os.Stat("internal/book/module.go")
	assert.NoError(t, err)
}

func TestWitchesLink_NoGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	err := WitchesLink("book", "https://github.com/DVV-15324/witches-book.git")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "go.mod not found")
}

func TestWitchesLink_GetWorkingDirectoryFail(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(origWd)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	err = os.WriteFile(
		"go.mod",
		[]byte("module captain"),
		0644,
	)
	assert.NoError(t, err)

	originalGetwd := getwdWitchesLink
	defer func() {
		getwdWitchesLink = originalGetwd
	}()

	getwdWitchesLink = func() (string, error) {
		return "", errors.New("mock getwd error")
	}

	err = WitchesLink(
		"book",
		"https://github.com/DVV-15324/witches-book.git",
	)

	assert.Error(t, err)
	assert.EqualError(
		t,
		err,
		"get working directory: mock getwd error",
	)
}

func TestWitchesLink_LinkDomainFail(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(origWd)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	err = os.WriteFile(
		"go.mod",
		[]byte("module captain"),
		0644,
	)
	assert.NoError(t, err)

	originalLink := linkWitchesLink
	defer func() {
		linkWitchesLink = originalLink
	}()

	linkWitchesLink = func(
		projectPath string,
		moduleName string,
		domainName string,
		repoURL string,
	) error {
		return errors.New("mock link domain error")
	}

	err = WitchesLink(
		"book",
		"https://github.com/DVV-15324/witches-book.git",
	)

	assert.Error(t, err)
	assert.EqualError(
		t,
		err,
		"link domain: mock link domain error",
	)
}

func TestWitchesLink_StatGoModFail(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(origWd)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	// go.mod không tồn tại ở đây.
	// os.Stat sẽ trả os.ErrNotExist, nên branch
	// "go.mod not found" sẽ chạy.
	err = WitchesLink(
		"book",
		"https://github.com/DVV-15324/witches-book.git",
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "go.mod not found")
}

func TestWitchesInstall(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)

	for _, driver := range []string{"mysql", "postgres", "mssql"} {
		t.Run(driver, func(t *testing.T) {
			WitchesInstall(driver)
		})
	}
}

func TestWitchesInstall_UnsupportedDriver(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)
	WitchesInstall("unsupported")
}

func TestFindProjectRoot(t *testing.T) {
	// Tạo thư mục mới hoàn toàn, không phải con của temp dir
	tmpDir := t.TempDir()

	// Tạo thư mục con sâu để tránh bị ảnh hưởng bởi go.mod ở cha
	projectRoot := filepath.Join(tmpDir, "project", "root")
	err := os.MkdirAll(projectRoot, 0755)
	require.NoError(t, err)

	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(projectRoot)

	// Test không có go.mod
	root := findProjectRoot()
	assert.Empty(t, root, "Should return empty when no go.mod found")

	// Test có go.mod
	err = os.WriteFile("go.mod", []byte("module test"), 0644)
	require.NoError(t, err)
	root = findProjectRoot()
	assert.Equal(t, projectRoot, root, "Should find root with go.mod")
}
func TestGenerateEasyJSONForAllDTOs_NoRoot(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	// Không có go.mod, generateEasyJSONForAllDTOs sẽ trả về error
	err := generateEasyJSONForAllDTOs()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot find project root")
}
func TestFindProjectRoot_WithNestedDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Tạo go.mod ở root
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644)
	require.NoError(t, err)

	nestedDir := filepath.Join(tmpDir, "cmd", "server")
	err = os.MkdirAll(nestedDir, 0755)
	require.NoError(t, err)

	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)

	// Chuyển vào nested dir
	err = os.Chdir(nestedDir)
	require.NoError(t, err)

	// findProjectRoot chỉ tìm trong thư mục hiện tại, nên không tìm thấy go.mod ở root
	root := findProjectRoot()
	assert.Empty(t, root, "Should return empty when go.mod not in current directory")

	// Tạo go.mod trong nested dir
	err = os.WriteFile(filepath.Join(nestedDir, "go.mod"), []byte("module test"), 0644)
	require.NoError(t, err)

	root = findProjectRoot()
	assert.Equal(t, nestedDir, root, "Should find root with go.mod in current directory")
}

func TestCreateContent(t *testing.T) {
	content := utils.CreateContent()
	assert.Contains(t, content, "APP_PORT=8080")
	assert.Contains(t, content, "DB_DRIVER")
}

func TestRunCmd(t *testing.T) {
	runCmd("go", "version")
	runCmd("nonexistent", "arg")
}

func TestRemoveEasyJSONFiles(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "test_easyjson.go"), []byte("// easyjson"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "normal.go"), []byte("package test"), 0644)

	removeEasyJSONFiles(tmpDir)

	_, err := os.Stat(filepath.Join(tmpDir, "test_easyjson.go"))
	assert.Error(t, err, "EasyJSON file should be removed")

	_, err = os.Stat(filepath.Join(tmpDir, "normal.go"))
	assert.NoError(t, err, "Normal file should remain")
}

func TestGenerateEasyJSONForAllDTOs_Success(t *testing.T) {
	if _, err := exec.LookPath("easyjson"); err != nil {
		t.Skip("easyjson not installed")
	}

	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)

	dtoDir := filepath.Join("internal", "user", "dto", "request")
	os.MkdirAll(dtoDir, 0755)

	dtoContent := `package request
type UserRequest struct {
    ID   int ` + "`json:\"id\"`" + `
    Name string ` + "`json:\"name\"`" + `
}`
	os.WriteFile(filepath.Join(dtoDir, "user_request.go"), []byte(dtoContent), 0644)

	err := generateEasyJSONForAllDTOs()
	assert.NoError(t, err)
}

func TestGenerateEasyJSONForAllDTOs_NoDTO(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(origWd)

	assert.NoError(t, os.Chdir(tmpDir))
	assert.NoError(t, os.WriteFile(
		"go.mod",
		[]byte("module test"),
		0644,
	))
	assert.NoError(t, os.MkdirAll("internal", 0755))

	err = generateEasyJSONForAllDTOs()

	assert.NoError(t, err)
}

func TestGenerateEasyJSONForDir_Success(t *testing.T) {
	if _, err := exec.LookPath("easyjson"); err != nil {
		t.Skip("easyjson not installed")
	}

	tmpDir := t.TempDir()
	dtoContent := `package test
type TestStruct struct {
    ID   int ` + "`json:\"id\"`" + `
    Name string ` + "`json:\"name\"`" + `
}`
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(dtoContent), 0644)

	err := generateEasyJSONForDir(tmpDir, "test")
	assert.NoError(t, err)
}

func TestGenerateEasyJSONForDir_InvalidDir(t *testing.T) {
	if _, err := exec.LookPath("easyjson"); err != nil {
		t.Skip("easyjson not installed")
	}

	tmpDir := t.TempDir()
	invalidDir := filepath.Join(tmpDir, "nonexistent")

	err := generateEasyJSONForDir(invalidDir, "test")
	assert.Error(t, err)
}

func TestWitchesRun_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)

	mainContent := `package main
import "fmt"
func main() {
    fmt.Println("Hello")
}`
	os.WriteFile("main.go", []byte(mainContent), 0644)

	dtoDir := filepath.Join("internal", "user", "dto", "request")
	os.MkdirAll(dtoDir, 0755)

	dtoContent := `package request
type UserRequest struct {
    ID   int ` + "`json:\"id\"`" + `
    Name string ` + "`json:\"name\"`" + `
}`
	os.WriteFile(filepath.Join(dtoDir, "user_request.go"), []byte(dtoContent), 0644)

	err := WitchesRun()
	if err != nil {
		t.Logf("WitchesRun returned error (expected in test): %v", err)
	}
}

func TestWitchesRun_NoMainGo(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)

	err := WitchesRun()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "run project")
}

func TestWitchesRun_NoGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)

	err := os.Chdir(tmpDir)
	require.NoError(t, err)

	err = WitchesRun()
	assert.Error(t, err)
	// Sửa expected message
	assert.Contains(t, err.Error(), "cannot find project root")
}
func TestWitchesRun_GoRunError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	// Tạo go.mod
	os.WriteFile("go.mod", []byte("module test"), 0644)

	// Tạo main.go có lỗi (syntax error)
	mainContent := `package main
import "fmt"
func main() {
    fmt.Println("Hello" // thiếu dấu đóng ngoặc
}`
	os.WriteFile("main.go", []byte(mainContent), 0644)

	err := WitchesRun()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "run project")
}
func TestWitchesLink_WithFullModules(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module captain"), 0644)
	os.MkdirAll("cmd/server/routers", 0755)
	os.MkdirAll("internal/shared/domain", 0755)
	os.MkdirAll("internal/shared/utils", 0755)

	// Tạo modules.go có InitModules
	modulesContent := `package routers

import "captain/internal/book"

type Modules struct {
    Book *book.Module
}

func (m *Modules) InitModules() *Modules {
    m.Book = &book.Module{}
    return m
}
`
	os.WriteFile("cmd/server/routers/modules.go", []byte(modulesContent), 0644)

	// Tạo routers.go có initModule
	routersContent := `package routers

func initModule(module interface{}) error {
    return nil
}

func SetupRoutes() error {
    return nil
}
`
	os.WriteFile("cmd/server/routers/routers.go", []byte(routersContent), 0644)

	// Tạo key_object.go đầy đủ
	keyObjectContent := `package utils

var (
	ObjectDefault int64 = 0
    ObjectUser   int64 = 1
)
`
	os.WriteFile("internal/shared/utils/key_object.go", []byte(keyObjectContent), 0644)

	err := WitchesLink("book", "https://github.com/DVV-15324/witches-book.git")
	assert.NoError(t, err)

	_, err = os.Stat("internal/book/module.go")
	assert.NoError(t, err)
}
func TestWitchesAdd_FullStructure(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)
	os.MkdirAll("cmd/server/routers", 0755)
	os.MkdirAll("internal/shared/domain", 0755)
	os.MkdirAll("internal/shared/utils", 0755)

	// Tạo các file cần thiết
	os.WriteFile("cmd/server/routers/modules.go", []byte(`package routers
type Modules struct{}
func (m *Modules) InitModules() *Modules { return &Modules{} }`), 0644)

	os.WriteFile("cmd/server/routers/routers.go", []byte(`package routers
func initModule(interface{}) error { return nil }
func SetupRoutes() error { return nil }`), 0644)

	os.WriteFile("internal/shared/utils/key_object.go", []byte(`package utils
var (ObjectDefault int64 = 0)`), 0644)

	err := WitchesAdd("user", "postgres")
	assert.NoError(t, err)

	_, err = os.Stat("internal/user/module.go")
	assert.NoError(t, err)
}
func TestWitchesRollback_FullStructure(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	os.WriteFile("go.mod", []byte("module test"), 0644)
	os.MkdirAll("internal/book", 0755)
	os.MkdirAll("cmd/server/routers", 0755)
	os.MkdirAll("internal/shared/utils", 0755)

	// Tạo các file cần thiết
	os.WriteFile("cmd/server/routers/modules.go", []byte(`package routers
type Modules struct {
    Book *book.Module
}
func (m *Modules) InitModules() *Modules { return m }`), 0644)

	os.WriteFile("cmd/server/routers/routers.go", []byte(`package routers
func initModule(interface{}) error { return nil }
func SetupRoutes() error { return nil }`), 0644)

	os.WriteFile("internal/shared/utils/key_object.go", []byte(`package utils
var (ObjectDefault int64 = 0)`), 0644)

	err := WitchesRollback("book")
	assert.NoError(t, err)

	_, err = os.Stat("internal/book")
	assert.Error(t, err)
}
func TestRemoveEasyJSONFiles_WithError(t *testing.T) {
	tmpDir := t.TempDir()

	// Tạo file nhưng không có quyền xóa (chỉ trên Unix)
	// Trên Windows, test này sẽ skip
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping permission test on Windows")
	}

	// Tạo file read-only
	filePath := filepath.Join(tmpDir, "test_easyjson.go")
	os.WriteFile(filePath, []byte("// easyjson"), 0644)
	os.Chmod(filePath, 0444) // read-only

	removeEasyJSONFiles(tmpDir)
	// Không panic, chỉ log error
}

func TestWitchesRun_GenerateEasyJSONFail(t *testing.T) {
	originalGenerate := generateAllDTOs
	defer func() {
		generateAllDTOs = originalGenerate
	}()

	generateAllDTOs = func() error {
		return errors.New("mock easyjson error")
	}

	err := WitchesRun()

	assert.Error(t, err)
	assert.EqualError(
		t,
		err,
		"generate easyjson: mock easyjson error",
	)
}

func TestWitchesRun_RunProjectFail(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(origWd)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	err = os.WriteFile(
		"go.mod",
		[]byte("module test"),
		0644,
	)
	assert.NoError(t, err)

	err = os.WriteFile(
		"main.go",
		[]byte(`package main

func main() {
	panic("mock run error")
}`),
		0644,
	)
	assert.NoError(t, err)

	err = WitchesRun()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "run project:")
}
func TestGenerateEasyJSONForAllDTOs_NoProjectRoot(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(origWd)

	assert.NoError(t, os.Chdir(tmpDir))

	err = generateEasyJSONForAllDTOs()

	assert.Error(t, err)
	assert.EqualError(t, err, "cannot find project root")
}
func TestGenerateEasyJSONForAllDTOs_RequestResponse(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(origWd)

	assert.NoError(t, os.Chdir(tmpDir))

	assert.NoError(t, os.WriteFile(
		"go.mod",
		[]byte("module test"),
		0644,
	))

	requestDir := filepath.Join(
		tmpDir,
		"internal", "user", "dto", "request",
	)

	responseDir := filepath.Join(
		tmpDir,
		"internal", "user", "dto", "response",
	)

	assert.NoError(t, os.MkdirAll(requestDir, 0755))
	assert.NoError(t, os.MkdirAll(responseDir, 0755))

	assert.NoError(t, os.WriteFile(
		filepath.Join(requestDir, "create.go"),
		[]byte(`package request

type CreateRequest struct {
	Name string `+"`json:\"name\"`"+`
}
`),
		0644,
	))

	assert.NoError(t, os.WriteFile(
		filepath.Join(responseDir, "user.go"),
		[]byte(`package response

type UserResponse struct {
	ID   int    `+"`json:\"id\"`"+`
	Name string `+"`json:\"name\"`"+`
}
`),
		0644,
	))

	err = generateEasyJSONForAllDTOs()

	assert.NoError(t, err)
}
func TestGenerateEasyJSONForDir_Error(t *testing.T) {
	tmpDir := t.TempDir()

	originalGenerate := generateEasyJSON
	defer func() {
		generateEasyJSON = originalGenerate
	}()

	generateEasyJSON = func(
		fset *token.FileSet,
		inputDir string,
		outputDir string,
	) error {
		return errors.New("mock generator error")
	}

	err := generateEasyJSONForDir(tmpDir, "request")

	assert.Error(t, err)
	assert.EqualError(
		t,
		err,
		"request error: mock generator error",
	)
}
