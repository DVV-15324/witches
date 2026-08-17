package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	wcmd_utils "github.com/DVV-15324/witches/cmd/utils"
	logger "github.com/DVV-15324/witches/pkg/core/response/logger"
	utils "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func getTestConfig() *wcmd_utils.Config {
	return &wcmd_utils.Config{RequestKey: "request_context"}
}

func getTestLogger(t *testing.T) *logger.ModelLogger {
	// Dùng thư mục tạm của hệ thống, không dùng t.TempDir() để tránh cleanup tự động
	dir, err := os.MkdirTemp("", "test_logs_")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "test_response.log")
	log, err := logger.NewFileLogger(path, 1, 20, 30)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = log.Sync()
		// Cố gắng đóng file nếu có
		// Nếu không có, vẫn xóa thư mục (lỗi sẽ được ignore)
		_ = os.RemoveAll(dir)
	})
	return log
}
func TestResponseSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/test", func(ctx *gin.Context) {
		WriteSuccess(ctx, map[string]string{"name": "vu"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response AppResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Status)
	assert.Equal(t, "Success", response.Message)
	assert.Equal(t, map[string]interface{}{"name": "vu"}, response.Data)
}

func TestResponseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/test", func(ctx *gin.Context) {
		errResp := NewAppError(404, errors.New("Not Found"), time.Now())
		WriteError(ctx, errResp)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	var response AppResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 404, response.Status)
	assert.Equal(t, "Not Found", response.Message)
}

func TestSuccessWithPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/test", func(ctx *gin.Context) {
		pagination := &utils.PaginationResponse{
			Page:       1,
			Limit:      10,
			Total:      100,
			TotalPages: 10,
			HasNext:    true,
			HasPrev:    false,
		}
		WriteSuccessWithPagination(ctx, map[string]string{"name": "vu"}, pagination)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response PaginationResponseWrapper
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Status)
	assert.Equal(t, "Success", response.Message)
	assert.NotNil(t, response.Pagination)
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 10, response.Pagination.Limit)
	assert.Equal(t, int64(100), response.Pagination.Total)
}

func TestResponseSuccessWithLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	cfg := getTestConfig()
	log := getTestLogger(t)

	r.GET("/test", func(ctx *gin.Context) {
		WriteSuccessWithLog(ctx, log, cfg, map[string]string{"name": "vu"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response AppResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Status)
	assert.Equal(t, "Success", response.Message)
}

func TestResponseErrorWithLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	cfg := getTestConfig()
	log := getTestLogger(t)

	r.GET("/test", func(ctx *gin.Context) {
		errResp := NewAppError(404, errors.New("Not Found"), time.Now())
		WriteErrorWithLog(ctx, log, cfg, errResp)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	var response AppResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 404, response.Status)
	assert.Equal(t, "Not Found", response.Message)
}

func TestWriteSuccessWithPaginationAndLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	cfg := getTestConfig()
	log := getTestLogger(t)

	r.GET("/test", func(ctx *gin.Context) {
		pagination := &utils.PaginationResponse{
			Page:       1,
			Limit:      10,
			Total:      100,
			TotalPages: 10,
			HasNext:    true,
			HasPrev:    false,
		}
		WriteSuccessWithPaginationAndLog(ctx, log, cfg, map[string]string{"name": "vu"}, pagination)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response PaginationResponseWrapper
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Status)
	assert.Equal(t, "Success", response.Message)
	assert.NotNil(t, response.Pagination)
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 10, response.Pagination.Limit)
	assert.Equal(t, int64(100), response.Pagination.Total)
}
