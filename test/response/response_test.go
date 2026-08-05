package test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	resp "github.com/DVV-15324/witches/pkg/core/response"
	logger "github.com/DVV-15324/witches/pkg/core/response/logger"
	utils "github.com/DVV-15324/witches/pkg/core/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResponseSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/test", func(ctx *gin.Context) {
		resp.WriteSuccess(ctx, map[string]string{"name": "vu"})
	})

	// Test request
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response resp.AppResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Status)
	assert.Equal(t, "Success", response.Message)
	assert.Equal(t, map[string]interface{}{"name": "vu"}, response.Data)
	t.Logf("Response success test passed: %+v", response)
}

func TestResponseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/test", func(ctx *gin.Context) {
		errResp := resp.NewAppError(404, errors.New("Not Found"), time.Now())
		resp.WriteError(ctx, errResp)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	var response resp.AppResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 404, response.Status)
	assert.Equal(t, "Not Found", response.Message)
	t.Logf("Response error test passed: %+v", response)
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
		resp.WriteSuccessWithPagination(ctx, map[string]string{"name": "vu"}, pagination)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response resp.PaginationResponseWrapper
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Status)
	assert.Equal(t, "Success", response.Message)
	assert.NotNil(t, response.Pagination)
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 10, response.Pagination.Limit)
	assert.Equal(t, int64(100), response.Pagination.Total)
	t.Logf("Pagination test passed: %+v", response)
}

func TestResponseSuccessWithLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Tạo logger cho test
	path := "./response_test.log"
	log, err := logger.NewFileLogger(path, 1, 20, 30)
	if err != nil {
		t.Fatal(err)
	}
	defer log.Sync()

	r.GET("/test", func(ctx *gin.Context) {
		resp.WriteSuccessWithLog(ctx, log, map[string]string{"name": "vu"}, "")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response resp.AppResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Status)
	assert.Equal(t, "Success", response.Message)
	t.Logf("Response success with log passed")
}

func TestResponseErrorWithLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Tạo logger cho test
	path := "./response_test.log"
	log, err := logger.NewFileLogger(path, 1, 20, 30)
	if err != nil {
		t.Fatal(err)
	}
	defer log.Sync()

	r.GET("/test", func(ctx *gin.Context) {
		errResp := resp.NewAppError(404, errors.New("Not Found"), time.Now())
		resp.WriteErrorWithLog(ctx, log, errResp, "")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	var response resp.AppResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 404, response.Status)
	assert.Equal(t, "Not Found", response.Message)
	t.Logf("Response error with log passed")
}

func TestWriteSuccessWithPaginationAndLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Tạo logger cho test
	path := "./response_test.log"
	log, err := logger.NewFileLogger(path, 1, 20, 30)
	if err != nil {
		t.Fatal(err)
	}
	defer log.Sync()

	r.GET("/test", func(ctx *gin.Context) {
		pagination := &utils.PaginationResponse{
			Page:       1,
			Limit:      10,
			Total:      100,
			TotalPages: 10,
			HasNext:    true,
			HasPrev:    false,
		}
		resp.WriteSuccessWithPaginationAndLog(ctx, log, map[string]string{"name": "vu"}, pagination, "")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response resp.PaginationResponseWrapper
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Status)
	assert.Equal(t, "Success", response.Message)
	assert.NotNil(t, response.Pagination)
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 10, response.Pagination.Limit)
	assert.Equal(t, int64(100), response.Pagination.Total)
	t.Logf("Pagination with log test passed")
}
