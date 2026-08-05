package test

import (
	"testing"

	logger "github.com/DVV-15324/witches/pkg/core/response/logger"
	"go.uber.org/zap"
)

func TestLogger_File(t *testing.T) {
	path := "./test.log"
	logger, err := logger.NewFileLogger(path, 1, 20, 30)
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()

	logger.Info("Hello, I'am Info")
	logger.Warn("Hello, I'am Warn")
	logger.Error("Hello, I'am Error")

	logger.InfoWithFields("Info", zap.String("Field", "hello"))
	logger.WarnWithFields("Info", zap.String("Field", "hello"))
	logger.ErrorWithFields("Info", zap.String("Field", "hello"))
}
