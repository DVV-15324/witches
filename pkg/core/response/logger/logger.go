package logger

import (
	//	"os"
	//"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type EntityLogger struct {
	Log   *zap.Logger
	Sugar *zap.SugaredLogger
}

// Khi chương trình kết thúc log có thể chưa kịp ghi ra
// Sync có nhiệm vụ đẩy(flush) dữ liệu được lưu trong buffer
// Đảm bảo dự liệu chắc chắn được lưu
func (l *EntityLogger) Sync() error {
	return l.Log.Sync()
}

// Hiển thị thông tin msg của log
func (l *EntityLogger) Info(msg string) {
	l.Log.Info(msg)
}

// Hiển thị thông tin msg của log với level Warn
func (l *EntityLogger) Warn(msg string) {
	l.Log.Warn(msg)
}

// Hiển thị thông tin msg của log với level Error
func (l *EntityLogger) Error(msg string) {
	l.Log.Error(msg)
}

// Thêm ngữ cảnh khi log [ví dụ]:
// logger.InfoWithFields("user login", zap.String("email", "vu@gmail.com"), zap.Int("user_id", 1),)
func (l *EntityLogger) InfoWithFields(msg string, fields ...zap.Field) {
	l.Log.Info(msg, fields...)
}

// Thêm ngữ cảnh khi log với level Warn [ví dụ]:
// logger.InfoWithFields("user login", zap.String("email", "vu@gmail.com"), zap.Int("user_id", 1),)
func (l *EntityLogger) WarnWithFields(msg string, fields ...zap.Field) {
	l.Log.Warn(msg, fields...)
}

// Thêm ngữ cảnh khi log với level Error [ví dụ]:
// logger.InfoWithFields("user login", zap.String("email", "vu@gmail.com"), zap.Int("user_id", 1),)
func (l *EntityLogger) ErrorWithFields(msg string, fields ...zap.Field) {
	l.Log.Error(msg, fields...)
}

// Thêm chèn định dạng văn bản(format) vào msg log [ví dụ]:
// logger.SugarInfo("user login %v", "vu@gmail.com")
func (l *EntityLogger) SugarInfof(msg string, args ...interface{}) {
	l.Sugar.Infof(msg, args...)
}

// Thêm chèn định dạng văn bản(format) vào msg log với level Warn [ví dụ]:
// logger.SugarInfo("user login %v", "vu@gmail.com")
func (l *EntityLogger) SugarWarnf(msg string, args ...interface{}) {
	l.Sugar.Warnf(msg, args...)
}

// Thêm chèn định dạng văn bản(format) vào msg log với level Error [ví dụ]:
// logger.SugarInfo("user login %v", "vu@gmail.com")
func (l *EntityLogger) SugarErrorf(msg string, args ...interface{}) {
	l.Sugar.Errorf(msg, args...)
}

// Thêm ngữ cảnh vào msg log theo kiểu KEYS và VALUES [ví dụ]:
// logger.SugarInfo("user login","email", "vu@gmail.com", "user_id", 1")
func (l *EntityLogger) SugarInfoWithFields(msg string, keysAndValues ...interface{}) {
	l.Sugar.Infow(msg, keysAndValues...)
}

// Thêm ngữ cảnh vào msg log theo kiểu KEYS và VALUES với level Warn [ví dụ]:
// logger.SugarInfo("user login","email", "vu@gmail.com", "user_id", 1")
func (l *EntityLogger) SugarWarnWithFields(msg string, keysAndValues ...interface{}) {
	l.Sugar.Warnw(msg, keysAndValues...)
}

// Thêm ngữ cảnh vào msg log theo kiểu KEYS và VALUES với level Error [ví dụ]:
// logger.SugarInfo("user login","email", "vu@gmail.com", "user_id", 1")
func (l *EntityLogger) SugarErrorWithFields(msg string, keysAndValues ...interface{}) {
	l.Sugar.Errorw(msg, keysAndValues...)
}

// Terminal (stdout) dễ debug
func NewDeloperLogger() (*EntityLogger, error) {

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return &EntityLogger{
		Log:   logger,
		Sugar: logger.Sugar(),
	}, nil
}

// Ghi file .log dạng JSON
// os.O_APPEND: ghi thêm vào cuối file (không ghi đè)
// os.O_CREATE: nếu file chưa tồn tại thì tạo mới
// os.O_WRONLY: chỉ cho phép ghi (write only)
// maxSize: megabytes
// maxAge: days
// compress: Nén bằng gzip hay không, mặc định không
func NewFileLogger(filePath string, maxSize int, maxBackUps int, maxAge int) (*EntityLogger, error) {

	// _, err := os.OpenFile(
	// 	filePath,
	// 	os.O_APPEND|os.O_CREATE|os.O_WRONLY,
	// 	0644,
	// )
	// if err != nil {
	// 	return nil, err
	// }
	// Cho Zap biết file lưu log
	// AddSysnc giúp file lưu chuẩn hóa theo Zap để có thể ghi log an toàn
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize, // megabytes
		MaxBackups: maxBackUps,
		MaxAge:     maxAge, //days
		Compress:   false,  // disabled by default)
	})
	//writeSyncer := safeWriteSyncer{ws}

	//Quyết định format log sang json chuẩn hóa theo Zap
	encoder := zapcore.NewJSONEncoder(
		zap.NewProductionEncoderConfig(),
	)

	// Bắt đầu khởi tạo với tạo zapcore(engine log) zap với cấu hình
	// writeSyncer: ghi ra đâu
	// encoder: format log sang json
	// level: lọc log
	// DEBUG  bị bỏ
	// INFO
	// WARN
	// ERROR
	core := zapcore.NewCore(
		encoder,
		writeSyncer,
		zap.InfoLevel,
	)

	// Logger thật sự bắt đầu
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	return &EntityLogger{
		Log:   logger,
		Sugar: logger.Sugar(),
	}, nil
}

// Stdout (Docker/K8s) hoặc gửi qua hệ thống log (ELK/Loki)
func NewProductionLogger() (*EntityLogger, error) {

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return &EntityLogger{
		Log:   logger,
		Sugar: logger.Sugar(),
	}, nil
}
