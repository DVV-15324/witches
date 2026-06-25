package test

import (
	"errors"
	resp "github.com/DVV-15324/witches/pkg/core/response_logger"
	logger "github.com/DVV-15324/witches/pkg/core/response_logger/logger"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestResponseLoggerSuccess(t *testing.T) {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		path := "./response_test.log"
		logger, err := logger.NewFileLogger(path, 1, 20, 30)
		if err != nil {
			t.Fatal(err)
		}
		defer logger.Sync()
		resp.WriteSuccessWithLog(ctx, logger, map[string]string{"name": "vu"})
	})
	r.Run(":8080")

}

func TestResponseLoggerError(t *testing.T) {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		err_resp := resp.NewErrorResponse(404, errors.New("Not Found"), time.Now())
		path := "./response_test.log"
		logger, err := logger.NewFileLogger(path, 1, 20, 30)
		if err != nil {
			t.Fatal(err)
		}
		defer logger.Sync()
		resp.WriteError(ctx, logger, err_resp)
	})
	r.Run(":8080")
}
