package database

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	utils "github.com/DVV-15324/witches/cmd/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateDBURL(t *testing.T) {
	tests := []struct {
		name     string
		driver   string
		user     string
		pass     string
		host     string
		dbname   string
		port     int64
		expected string
		wantErr  bool
	}{
		{
			name:     "MySQL",
			driver:   "mysql",
			user:     "root",
			pass:     "secret",
			host:     "localhost",
			dbname:   "mydb",
			port:     3306,
			expected: "root:secret@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local",
			wantErr:  false,
		},
		{
			name:     "PostgreSQL (postgres)",
			driver:   "postgres",
			user:     "admin",
			pass:     "123",
			host:     "127.0.0.1",
			dbname:   "testdb",
			port:     5432,
			expected: "admin:123@127.0.0.1:5432/testdb?sslmode=disable",
			wantErr:  false,
		},
		{
			name:     "SQL Server (mssql)",
			driver:   "mssql",
			user:     "sa",
			pass:     "Pass@word",
			host:     "localhost",
			dbname:   "master",
			port:     1433,
			expected: "sa:Pass@word@localhost:1433?database=master&encrypt=disable",
			wantErr:  false,
		},
		{
			name:     "SQL Server (sqlserver)",
			driver:   "sqlserver",
			user:     "sa",
			pass:     "Pass@word",
			host:     "localhost",
			dbname:   "master",
			port:     1433,
			expected: "sa:Pass@word@localhost:1433?database=master&encrypt=disable",
			wantErr:  false,
		},
		{
			name:     "Unsupported driver",
			driver:   "oracle",
			user:     "system",
			pass:     "pass",
			host:     "localhost",
			dbname:   "xe",
			port:     1521,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateDBURL(tt.driver, tt.user, tt.pass, tt.host, tt.dbname, tt.port)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestWitchesDBURL(t *testing.T) {
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(origDir)
	}()

	tmpDir := t.TempDir()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	tests := []struct {
		name     string
		driver   string
		config   *utils.Config
		wantErr  bool
		checkEnv func(t *testing.T, content string)
	}{
		{
			name:   "PostgreSQL (postgresql) - valid",
			driver: "postgresql",
			config: &utils.Config{
				DBUser:     "admin",
				DBPassword: "123",
				DBHost:     "127.0.0.1",
				DBName:     "db",
				DBPort:     5432,
			},
			wantErr: false,
		},
		{
			name:   "SQL Server (mssql) - valid",
			driver: "mssql",
			config: &utils.Config{
				DBUser:     "sa",
				DBPassword: "Pass@word",
				DBHost:     "localhost",
				DBName:     "master",
				DBPort:     1433,
			},
			wantErr: false,
		},
		{
			name:   "SQL Server (sqlserver) - valid",
			driver: "sqlserver",
			config: &utils.Config{
				DBUser:     "sa",
				DBPassword: "Pass@word",
				DBHost:     "localhost",
				DBName:     "master",
				DBPort:     1433,
			},
			wantErr: false,
		},
		{
			name:   "MySQL - valid",
			driver: "mysql",
			config: &utils.Config{
				DBUser:     "user",
				DBPassword: "pass",
				DBHost:     "localhost",
				DBName:     "test",
				DBPort:     3306,
			},
			wantErr: false,
		},
		{
			name:   "PostgreSQL - valid",
			driver: "postgres",
			config: &utils.Config{
				DBUser:     "admin",
				DBPassword: "123",
				DBHost:     "127.0.0.1",
				DBName:     "db",
				DBPort:     5432,
			},
			wantErr: false,
		},
		{
			name:   "Unsupported driver",
			driver: "mongodb",
			config: &utils.Config{
				DBUser:     "root",
				DBPassword: "pass",
				DBHost:     "localhost",
				DBName:     "test",
				DBPort:     27017,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WitchesDBURL(tt.driver, tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			envPath := filepath.Join(tmpDir, "witches.env")
			_, err = os.Stat(envPath)
			assert.NoError(t, err, "witches.env should exist")
			content, err := os.ReadFile(envPath)
			assert.NoError(t, err)
			assert.NotEmpty(t, content)
			assert.NotEmpty(t, tt.config.DBUrl)
			assert.Contains(t, string(content), tt.config.DBUrl)
		})
	}
}

func TestWitchesDBURL_EdgeCases(t *testing.T) {
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(origDir)
	}()

	tmpDir := t.TempDir()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	config := &utils.Config{
		DBUser:     "test",
		DBPassword: "pass",
		DBHost:     "localhost",
		DBName:     "testdb",
		DBPort:     3306,
	}

	t.Run("File không tồn tại - tạo mới", func(t *testing.T) {
		os.Remove("witches.env")

		err := WitchesDBURL("mysql", config)
		assert.NoError(t, err)

		content, err := os.ReadFile("witches.env")
		assert.NoError(t, err)
		assert.Contains(t, string(content), "DB_URL=")
		assert.Contains(t, string(content), "test:pass@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local")
	})

	t.Run("File tồn tại nhưng không có DB_URL", func(t *testing.T) {
		envContent := `APP_NAME=test
APP_PORT=8080
# DB_URL=some_old_url`
		err := os.WriteFile("witches.env", []byte(envContent), 0644)
		require.NoError(t, err)

		err = WitchesDBURL("mysql", config)
		assert.NoError(t, err)

		content, err := os.ReadFile("witches.env")
		assert.NoError(t, err)

		assert.Contains(t, string(content), "DB_URL=test:pass@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local")
		assert.Contains(t, string(content), "APP_NAME=test")
		assert.Contains(t, string(content), "APP_PORT=8080")
		assert.Contains(t, string(content), "# DB_URL=some_old_url")
	})
}

func TestWitchesDBURL_ReadFileError(t *testing.T) {
	// Chỉ chạy trên Unix (Linux, macOS) vì Windows không hỗ trợ tốt quyền đọc/ghi file
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows: file permission handling differs")
	}

	// Lưu thư mục làm việc hiện tại để khôi phục sau
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(origDir)
	}()

	// Tạo thư mục tạm và chuyển vào đó
	tmpDir := t.TempDir()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Tạo file witches.env với nội dung bất kỳ
	initialContent := `APP_NAME=test
DB_URL=mysql://root:secret@tcp(localhost:3306)/mydb`
	err = os.WriteFile("witches.env", []byte(initialContent), 0644)
	require.NoError(t, err)

	// Thu hồi toàn bộ quyền đọc, ghi, thực thi cho file
	// (trên Unix, owner vẫn có thể xoá file, nhưng không đọc/ghi được)
	err = os.Chmod("witches.env", 0000)
	require.NoError(t, err)

	// Cấu hình hợp lệ
	config := &utils.Config{
		DBUser:     "test",
		DBPassword: "pass",
		DBHost:     "localhost",
		DBName:     "testdb",
		DBPort:     3306,
	}

	// Gọi hàm cần test – kỳ vọng trả về lỗi vì không đọc được file
	err = WitchesDBURL("mysql", config)
	assert.Error(t, err, "Expected error when witches.env is not readable")

	// (Tuỳ chọn) Kiểm tra thông báo lỗi có chứa từ khoá liên quan
	assert.Contains(t, err.Error(), "read witches.env", "Error message should indicate read failure")
}

func TestWitchesDBURL_WriteFileError(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, "witches.env")
	err := os.WriteFile(envPath, []byte("DB_URL=old"), 0444) // read-only
	require.NoError(t, err)
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	config := &utils.Config{
		DBUser:     "test",
		DBPassword: "pass",
		DBHost:     "localhost",
		DBName:     "testdb",
		DBPort:     3306,
	}
	err = WitchesDBURL("mysql", config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "write witches.env")
}
func TestWitchesDBURL_ReadFileFail(t *testing.T) {
	originalReadFile := readFileWitchesDBURL
	defer func() {
		readFileWitchesDBURL = originalReadFile
	}()

	readFileWitchesDBURL = func(
		filename string,
	) ([]byte, error) {
		return nil, errors.New("mock read file error")
	}

	config := &utils.Config{
		DBUser:     "root",
		DBPassword: "password",
		DBHost:     "localhost",
		DBName:     "test",
		DBPort:     3306,
	}

	err := WitchesDBURL("mysql", config)

	assert.Error(t, err)
	assert.EqualError(
		t,
		err,
		"read witches.env: mock read file error",
	)
}
