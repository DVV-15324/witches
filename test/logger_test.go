package test

import (
	logger "core-v/pkg/core/response_logger/logger"
	"go.uber.org/zap"
	"testing"
)

func TestLogger_WithInfoProduct(t *testing.T) {

	logger, err := logger.NewDeloperLogger()
	if err != nil {
		t.Fatal(err)
	}

	defer logger.Sync()

	logger.Info("Hello product")
}

func TestLogger_WithInfoDev(t *testing.T) {

	logger, err := logger.NewProductionLogger()
	if err != nil {
		t.Fatal(err)
	}

	defer logger.Sync()

	logger.Info("Hello")
	logger.InfoWithFields("Hello", zap.String("email", "vu@gmail.com"))

	logger.Warn("Hello")
	logger.WarnWithFields("Hello", zap.String("email", "vu@gmail.com"))

	logger.Error("Hello")
	logger.ErrorWithFields("Hello", zap.String("email", "vu@gmail.com"))
}

func TestLogger_SugarInfo(t *testing.T) {

	logger, err := logger.NewProductionLogger()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()

	logger.SugarInfof("Hello %v", "Sugar")
	logger.SugarInfoWithFields("Hello", "email", "vu@gmail.com")

	logger.SugarWarnf("Hello %v", "Sugar")
	logger.SugarWarnWithFields("Hello", "email", "vu@gmail.com")

	logger.SugarErrorf("Hello %v", "Sugar")
	logger.SugarErrorWithFields("Hello", "email", "vu@gmail.com")
}

func TestLogger_File(t *testing.T) {
	path := "./test.log"
	logger, err := logger.NewFileLogger(path, 1, 20, 30)
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()

	logger.SugarInfof("Hello %v", "Sugar")
	logger.SugarInfoWithFields(
		"Hello",
		"email", "vu@gmail.com",
	)
	logger.SugarWarnf("Hello %v", "Sugar")
	logger.SugarWarnWithFields("Hello", "email", "vu@gmail.com")

	logger.SugarErrorf("Hello %v", "Sugar")
	logger.SugarErrorWithFields("Hello", "email", "vu@gmail.com")

}
