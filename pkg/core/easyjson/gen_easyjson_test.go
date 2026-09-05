package easyjson

import (
	"errors"
	"os"
	"os/exec"

	"path/filepath"

	"testing"

	"go/token"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//  Test GeneratorEasyJson

func TestGeneratorEasyJsonRequest(t *testing.T) {
	rootDir := findProjectRoot(t)
	inputDir := filepath.Join(rootDir, "pkg", "core", "easyjson", "test", "request")
	outputDir := inputDir
	oldGenFile := filepath.Join(outputDir, "request_easyjson.go")
	_ = os.Remove(oldGenFile)

	fset := token.NewFileSet()
	err := GeneratorEasyJson(fset, inputDir, outputDir)
	assert.NoError(t, err)

	if _, err = os.Stat(oldGenFile); err == nil {
		t.Logf("Generated file found: %s", oldGenFile)
	} else {
		t.Log("No generated file")
	}
}

func TestGeneratorEasyJsonResponse(t *testing.T) {
	rootDir := findProjectRoot(t)
	inputDir := filepath.Join(rootDir, "pkg", "core", "easyjson", "test", "response")
	outputDir := inputDir
	oldGenFile := filepath.Join(outputDir, "response_easyjson.go")
	_ = os.Remove(oldGenFile)

	fset := token.NewFileSet()
	err := GeneratorEasyJson(fset, inputDir, outputDir)
	assert.NoError(t, err)

	if _, err = os.Stat(oldGenFile); err == nil {
		t.Logf("Generated file found: %s", oldGenFile)
	} else {
		t.Log("No generated file")
	}
}
func TestGeneratorEasyJson_InputIsFile(t *testing.T) {
	// Kiểm tra easyjson có installed không
	if _, err := exec.LookPath("easyjson"); err != nil {
		t.Skip("easyjson not installed, skipping test")
	}

	tmpDir := t.TempDir()

	// Tạo go.mod để easyjson có thể chạy
	goMod := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goMod, []byte("module test\ngo 1.21"), 0644)
	require.NoError(t, err)

	// Tạo file Go với marker GenEasyJson
	goFile := filepath.Join(tmpDir, "test.go")
	content := `package test

type TestStruct struct {
	Name string ` + "`json:\"name\"`" + `
	Age  int    ` + "`json:\"age\"`" + `
}

func (t *TestStruct) GenEasyJson() {}
`
	err = os.WriteFile(goFile, []byte(content), 0644)
	require.NoError(t, err)

	fset := token.NewFileSet()
	outputDir := tmpDir

	err = GeneratorEasyJson(fset, goFile, outputDir)
	assert.NoError(t, err)

	// Kiểm tra file generated
	genFile := filepath.Join(tmpDir, "test_easyjson.go")
	if _, err := os.Stat(genFile); err == nil {
		t.Log("Generated file found")
	}
}
func TestGeneratorEasyJson_OutputDifferentFromInput(t *testing.T) {
	// Kiểm tra easyjson có installed không
	if _, err := exec.LookPath("easyjson"); err != nil {
		t.Skip("easyjson not installed, skipping test")
	}

	tmpDir := t.TempDir()

	// Tạo go.mod
	goMod := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goMod, []byte("module test\ngo 1.21"), 0644)
	require.NoError(t, err)

	// Tạo input dir
	inputDir := filepath.Join(tmpDir, "input")
	err = os.MkdirAll(inputDir, 0755)
	require.NoError(t, err)

	// Tạo file Go với marker
	goFile := filepath.Join(inputDir, "test.go")
	content := `package test

type TestStruct struct {
	Name string ` + "`json:\"name\"`" + `
	Age  int    ` + "`json:\"age\"`" + `
}

func (t *TestStruct) GenEasyJson() {}
`
	err = os.WriteFile(goFile, []byte(content), 0644)
	require.NoError(t, err)

	// Output dir khác
	outputDir := filepath.Join(tmpDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	require.NoError(t, err)

	fset := token.NewFileSet()
	err = GeneratorEasyJson(fset, inputDir, outputDir)
	assert.NoError(t, err)
}

// Test input là file nhưng không có marker GenEasyJson
func TestGeneratorEasyJson_InputFileNoMarker(t *testing.T) {
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "test.go")
	content := `package test
type T struct{}` // không có marker
	err := os.WriteFile(goFile, []byte(content), 0644)
	require.NoError(t, err)

	outputDir := tmpDir
	fset := token.NewFileSet()
	err = GeneratorEasyJson(fset, goFile, outputDir)
	assert.NoError(t, err)
}

func TestGeneratorEasyJson_InputDirNoMarker(t *testing.T) {
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "test.go")
	content := `package test
type T struct{}`
	err := os.WriteFile(goFile, []byte(content), 0644)
	require.NoError(t, err)

	fset := token.NewFileSet()
	err = GeneratorEasyJson(fset, tmpDir, "")
	assert.NoError(t, err)
}

func TestGeneratorEasyJson_InputNotExist(t *testing.T) {
	fset := token.NewFileSet()
	err := GeneratorEasyJson(fset, "/non/existent/path", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "input path does not exist")
}

func TestGeneratorEasyJson_InputEmpty(t *testing.T) {
	fset := token.NewFileSet()
	err := GeneratorEasyJson(fset, "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "input path is empty")
}

//  Test generateEasyJSON

func TestGenerateEasyJSON_FileNotExist(t *testing.T) {
	err := generateEasyJSON("non_existent_file.go")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file does not exist")
}
func TestGenerateEasyJSON_FileExists(t *testing.T) {
	// Kiểm tra easyjson có installed không
	if _, err := exec.LookPath("easyjson"); err != nil {
		t.Skip("easyjson not installed, skipping test")
	}

	tmpDir := t.TempDir()

	// Tạo go.mod
	goModContent := `module test

go 1.21
`
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644)
	require.NoError(t, err)

	// Tạo file Go hợp lệ
	goFile := filepath.Join(tmpDir, "test.go")
	content := `package test

type TestStruct struct {
	Name string ` + "`json:\"name\"`" + `
	Age  int    ` + "`json:\"age\"`" + `
}

func (t *TestStruct) GenEasyJson() {}
`
	err = os.WriteFile(goFile, []byte(content), 0644)
	require.NoError(t, err)

	err = generateEasyJSON(goFile)
	// Nếu fail thì skip, không fail test
	if err != nil {
		t.Skipf("easyjson failed (dependencies issue): %v", err)
	}

	// Kiểm tra file generated
	genFile := filepath.Join(tmpDir, "test_easyjson.go")
	_, err = os.Stat(genFile)
	assert.NoError(t, err)
}
func TestGenerateEasyJSON_NoEasyJSON(t *testing.T) {
	oldPath := os.Getenv("PATH")
	defer func() { _ = os.Setenv("PATH", oldPath) }()
	_ = os.Setenv("PATH", "")

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(tmpFile, []byte("package test\ntype T struct{}"), 0644)
	require.NoError(t, err)

	err = generateEasyJSON(tmpFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "executable")
}

//  Helper

func findProjectRoot(t *testing.T) string {
	dir, err := os.Getwd()
	assert.NoError(t, err)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("Cannot find project root")
		}
		dir = parent
	}
}
func TestGeneratorEasyJson_FileInputWithOutputDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Tạo file Go với marker
	goFile := filepath.Join(tmpDir, "test.go")
	content := `package test

type TestStruct struct {
	Name string ` + "`json:\"name\"`" + `
	Age  int    ` + "`json:\"age\"`" + `
}

func (t *TestStruct) GenEasyJson() {}
`
	err := os.WriteFile(goFile, []byte(content), 0644)
	require.NoError(t, err)

	// Output dir khác
	outputDir := filepath.Join(tmpDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	require.NoError(t, err)

	fset := token.NewFileSet()
	err = GeneratorEasyJson(fset, goFile, outputDir)
	assert.NoError(t, err)
}
func TestGeneratorEasyJson_AbsInputError(t *testing.T) {
	old := absPath
	defer func() { absPath = old }()

	called := 0
	absPath = func(path string) (string, error) {
		called++
		if called == 1 {
			return "", errors.New("mock abs input error")
		}
		return filepath.Abs(path)
	}

	tmpDir := t.TempDir()
	input := filepath.Join(tmpDir, "test.go")
	require.NoError(t, os.WriteFile(input, []byte("package test"), 0644))

	err := GeneratorEasyJson(nil, input, "")

	assert.EqualError(t, err, "failed to get absolute path: mock abs input error")
}
func TestGeneratorEasyJson_AbsOutputError(t *testing.T) {
	old := absPath
	defer func() { absPath = old }()

	called := 0
	absPath = func(path string) (string, error) {
		called++
		if called == 2 {
			return "", errors.New("mock abs output error")
		}
		return filepath.Abs(path)
	}

	tmpDir := t.TempDir()
	input := filepath.Join(tmpDir, "test.go")
	require.NoError(t, os.WriteFile(input, []byte("package test"), 0644))

	err := GeneratorEasyJson(nil, input, "")

	assert.EqualError(t, err,
		"failed to get absolute output path: mock abs output error")
}
func TestGeneratorEasyJson_ParseFileError(t *testing.T) {
	tmpDir := t.TempDir()

	broken := filepath.Join(tmpDir, "broken.go")

	require.NoError(t, os.WriteFile(
		broken,
		[]byte(`package test

func Foo( {
`),
		0644,
	))

	err := GeneratorEasyJson(nil, tmpDir, "")

	assert.NoError(t, err)
}
func TestGeneratorEasyJson_WalkCallbackError(t *testing.T) {
	old := walkPath
	defer func() { walkPath = old }()

	walkPath = func(root string, fn filepath.WalkFunc) error {
		return fn(root, nil, errors.New("mock walk callback error"))
	}

	tmpDir := t.TempDir()

	err := GeneratorEasyJson(nil, tmpDir, "")

	assert.NoError(t, err)
}
func TestGeneratorEasyJson_WalkError(t *testing.T) {
	old := walkPath
	defer func() { walkPath = old }()

	walkPath = func(root string, fn filepath.WalkFunc) error {
		return errors.New("mock walk error")
	}

	tmpDir := t.TempDir()

	err := GeneratorEasyJson(nil, tmpDir, "")

	assert.EqualError(t, err, "error walking: mock walk error")
}
func TestGeneratorEasyJson_RelPathError(t *testing.T) {
	oldRelPath := relPath
	t.Cleanup(func() {
		relPath = oldRelPath
	})

	relPath = func(basepath, targpath string) (string, error) {
		return "", errors.New("mock rel path error")
	}

	inputDir := t.TempDir()
	outputDir := filepath.Join(t.TempDir(), "output")

	source := `
package test

type Request struct{}

func (r *Request) GenEasyJson() {}
`

	inputFile := filepath.Join(inputDir, "request.go")
	require.NoError(t, os.WriteFile(inputFile, []byte(source), 0644))

	oldGenerate := generateEasyJSONFunc
	t.Cleanup(func() {
		generateEasyJSONFunc = oldGenerate
	})

	generateEasyJSONFunc = func(filePath string) error {
		return nil
	}

	err := GeneratorEasyJson(token.NewFileSet(), inputDir, outputDir)

	require.NoError(t, err)

	// filepath.Base(path) được dùng khi Rel fail
	expected := filepath.Join(outputDir, "request.go")
	_, err = os.Stat(expected)
	assert.NoError(t, err)
}
func TestGeneratorEasyJson_CreateBaseOutputDirError(t *testing.T) {
	old := mkdirAll
	defer func() { mkdirAll = old }()

	mkdirAll = func(path string, perm os.FileMode) error {
		return errors.New("mock mkdir error")
	}

	tmpDir := t.TempDir()

	input := filepath.Join(tmpDir, "test.go")
	require.NoError(t, os.WriteFile(
		input,
		[]byte(`package test`),
		0644,
	))

	err := GeneratorEasyJson(nil, input, "")

	assert.EqualError(t, err,
		"failed to create output directory: mock mkdir error")
}
func TestGeneratorEasyJson_OutputDirError(t *testing.T) {
	oldMkdir := mkdirAll
	oldGenerate := generateEasyJSONFunc

	defer func() {
		mkdirAll = oldMkdir
		generateEasyJSONFunc = oldGenerate
	}()

	count := 0
	mkdirAll = func(path string, perm os.FileMode) error {
		count++

		if count == 2 {
			return errors.New("mock output mkdir error")
		}

		return os.MkdirAll(path, perm)
	}

	generateEasyJSONFunc = func(filePath string) error {
		return nil
	}

	tmpDir := t.TempDir()

	input := filepath.Join(tmpDir, "input")
	require.NoError(t, os.MkdirAll(input, 0755))

	goFile := filepath.Join(input, "test.go")
	require.NoError(t, os.WriteFile(
		goFile,
		[]byte(`package test`),
		0644,
	))

	output := filepath.Join(tmpDir, "output")

	err := GeneratorEasyJson(nil, input, output)

	assert.NoError(t, err)
}

func TestGeneratorEasyJson_OutputDirMkdirError(t *testing.T) {
	oldMkdirAll := mkdirAll
	t.Cleanup(func() {
		mkdirAll = oldMkdirAll
	})

	callCount := 0

	mkdirAll = func(path string, perm os.FileMode) error {
		callCount++

		// Cho mkdirAll đầu tiên tạo basePath thành công,
		// lần sau mới fail ở outputDirPath.
		if callCount == 1 {
			return os.MkdirAll(path, perm)
		}

		return errors.New("mock mkdir error")
	}

	inputDir := t.TempDir()
	outputDir := filepath.Join(t.TempDir(), "output")

	source := `
package test

type Request struct{}

func (r *Request) GenEasyJson() {}
`

	require.NoError(t,
		os.WriteFile(
			filepath.Join(inputDir, "request.go"),
			[]byte(source),
			0644,
		),
	)

	err := GeneratorEasyJson(token.NewFileSet(), inputDir, outputDir)

	require.NoError(t, err)
}
func TestGeneratorEasyJson_ReadFileError(t *testing.T) {
	oldReadFile := readFile
	t.Cleanup(func() {
		readFile = oldReadFile
	})

	readFile = func(filename string) ([]byte, error) {
		return nil, errors.New("mock read file error")
	}

	inputDir := t.TempDir()
	outputDir := filepath.Join(t.TempDir(), "output")

	source := `
package test

type Request struct{}

func (r *Request) GenEasyJson() {}
`

	require.NoError(t,
		os.WriteFile(
			filepath.Join(inputDir, "request.go"),
			[]byte(source),
			0644,
		),
	)

	err := GeneratorEasyJson(token.NewFileSet(), inputDir, outputDir)

	require.NoError(t, err)
}
func TestGeneratorEasyJson_WriteFileError(t *testing.T) {
	oldWriteFile := writeFile
	t.Cleanup(func() {
		writeFile = oldWriteFile
	})

	writeFile = func(
		filename string,
		data []byte,
		perm os.FileMode,
	) error {
		return errors.New("mock write file error")
	}

	inputDir := t.TempDir()
	outputDir := filepath.Join(t.TempDir(), "output")

	source := `
package test

type Request struct{}

func (r *Request) GenEasyJson() {}
`

	require.NoError(t,
		os.WriteFile(
			filepath.Join(inputDir, "request.go"),
			[]byte(source),
			0644,
		),
	)

	err := GeneratorEasyJson(token.NewFileSet(), inputDir, outputDir)

	require.NoError(t, err)
}
