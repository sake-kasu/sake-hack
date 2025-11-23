package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// コンテキストキー型
type contextKey string

// コンテキストキー定数
const (
	RequestIDKey contextKey = "request_id"
	TraceIDKey   contextKey = "trace_id"
	UserIDKey    contextKey = "user_id"
)

// Init はロガーを初期化する
func Init(level, format string) error {
	var config zap.Config

	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// ログレベル設定
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	var err error
	globalLogger, err = config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	return nil
}

// Get はグローバルロガーを取得する
func Get() *zap.Logger {
	if globalLogger == nil {
		// フォールバック: 初期化されていない場合
		globalLogger, _ = zap.NewDevelopment()
	}
	return globalLogger
}

// Sync はログバッファをフラッシュする
func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}

// WithContext はコンテキストからロガーを取得する
// request_id, trace_id, user_id などのコンテキスト情報を含むロガーを返す
func WithContext(ctx context.Context) *zap.Logger {
	logger := Get()

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		logger = logger.With(zap.String("request_id", requestID))
	}

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		logger = logger.With(zap.String("trace_id", traceID))
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		logger = logger.With(zap.String("user_id", userID))
	}

	return logger
}

// Info はINFOレベルのログを出力する
func Info(ctx context.Context, message string, fields ...zap.Field) {
	WithContext(ctx).Info(message, fields...)
}

// Debug はDEBUGレベルのログを出力する
func Debug(ctx context.Context, message string, fields ...zap.Field) {
	WithContext(ctx).Debug(message, fields...)
}

// Warn はWARNレベルのログを出力する
func Warn(ctx context.Context, message string, fields ...zap.Field) {
	WithContext(ctx).Warn(message, fields...)
}

// Error はERRORレベルのログを出力する
func Error(ctx context.Context, message string, fields ...zap.Field) {
	WithContext(ctx).Error(message, fields...)
}

// LogDatabaseError はデータベースエラーをログ出力する
func LogDatabaseError(ctx context.Context, operation, table string, err error, details ...map[string]interface{}) {
	fields := []zap.Field{
		zap.String("error_type", "database"),
		zap.String("operation", operation),
		zap.String("table", table),
		zap.Error(err),
	}

	// detailsをzap.Fieldに変換
	if len(details) > 0 {
		for k, v := range details[0] {
			fields = append(fields, zap.Any(k, v))
		}
	}

	WithContext(ctx).Error("Database error occurred", fields...)
}

// LogBusinessError はビジネスロジックエラーをログ出力する
func LogBusinessError(ctx context.Context, rule string, err error, details ...map[string]interface{}) {
	fields := []zap.Field{
		zap.String("error_type", "business"),
		zap.String("rule", rule),
		zap.Error(err),
	}

	// detailsをzap.Fieldに変換
	if len(details) > 0 {
		for k, v := range details[0] {
			fields = append(fields, zap.Any(k, v))
		}
	}

	WithContext(ctx).Warn("Business rule violation", fields...)
}

// LogValidationError はバリデーションエラーをログ出力する
func LogValidationError(ctx context.Context, field string, value interface{}, reason string, details ...map[string]interface{}) {
	fields := []zap.Field{
		zap.String("error_type", "validation"),
		zap.String("field", field),
		zap.Any("value", value),
		zap.String("reason", reason),
	}

	// detailsをzap.Fieldに変換
	if len(details) > 0 {
		for k, v := range details[0] {
			fields = append(fields, zap.Any(k, v))
		}
	}

	WithContext(ctx).Warn("Validation error", fields...)
}

// TraceMethodAuto はメソッドの開始と終了を自動でログ出力する
// 使用例: defer logger.TraceMethodAuto(ctx, params)()
func TraceMethodAuto(ctx context.Context, params interface{}) func() {
	start := time.Now()

	// 呼び出し元の関数名を取得
	pc, _, _, ok := runtime.Caller(1)
	methodName := "unknown"
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			methodName = fn.Name()
		}
	}

	logger := WithContext(ctx)

	// メソッド開始ログ
	logger.Debug("Method started",
		zap.String("method", methodName),
		zap.String("phase", "start"),
	)

	return func() {
		// メソッド終了ログ
		duration := time.Since(start)
		logger.Debug("Method completed",
			zap.String("method", methodName),
			zap.String("phase", "end"),
			zap.Int64("duration_ms", duration.Milliseconds()),
		)
	}
}

// InitDefault は開発環境用のデフォルトロガーを初期化する
func InitDefault() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	var err error
	globalLogger, err = config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
}
