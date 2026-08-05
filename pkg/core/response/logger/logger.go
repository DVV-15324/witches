package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type EntityLogger struct {
	Log *zap.Logger
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

// Ghi file .log dạng JSON
// maxSize: megabytes
// maxAge: days
// compress: Nén bằng gzip hay không, mặc định không
func NewFileLogger(filePath string, maxSize int, maxBackUps int, maxAge int) (*EntityLogger, error) {
	// Cho Zap biết file lưu log
	// AddSysnc giúp file lưu chuẩn hóa theo Zap để có thể ghi log an toàn
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize, // megabytes
		MaxBackups: maxBackUps,
		MaxAge:     maxAge, //days
		Compress:   true,   //gzip
		LocalTime:  true,
	})

	//Quyết định format log sang json chuẩn hóa theo Zap
	encoder := zapcore.NewJSONEncoder(
		zap.NewProductionEncoderConfig(),
	)

	// Bắt đầu khởi tạo với tạo zapcore(engine log) zap với cấu hình
	// writeSyncer: ghi ra đâu
	// encoder: format log sang json
	// level: lọc log
	// DEBUG  bị bỏ ở đây không sài debug :))
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
		zap.AddCaller(),      //Có file + line
		zap.AddCallerSkip(1), //Bỏ qua 1 layer wrapper(nó chính là gorm wrap), trỏ đúng nơi gọi log
	)

	return &EntityLogger{
		Log: logger,
	}, nil
}
