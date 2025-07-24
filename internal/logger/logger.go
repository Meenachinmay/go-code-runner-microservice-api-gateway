package logger

import (
	"context"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const (
	loggerKey      contextKey = "logger"
	correlationKey contextKey = "correlation_id"
	requestIDKey   contextKey = "request_id"
)

var (
	globalLogger *zap.Logger
	globalSugar  *zap.SugaredLogger
)

type Config struct {
	Level       string `yaml:"level"`
	Environment string `yaml:"environment"`
	ServiceName string `yaml:"service_name"`
}

func Initialize(cfg Config) error {
	var (
		zapConfig zap.Config
		err       error
	)

	if cfg.ServiceName == "" {
		cfg.ServiceName = "api-gateway"
	}
	if cfg.Environment == "" {
		cfg.Environment = "development"
	}
	if cfg.Level == "" {
		cfg.Level = "info"
	}

	if cfg.Environment == "production" {
		zapConfig = zap.NewProductionConfig()
		zapConfig.DisableStacktrace = true
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.DisableStacktrace = true
	}

	level, err := zapcore.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		return err
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	zapConfig.InitialFields = map[string]interface{}{
		"service": cfg.ServiceName,
		"env":     cfg.Environment,
		"pid":     os.Getpid(),
	}

	globalLogger, err = zapConfig.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		return err
	}

	globalSugar = globalLogger.Sugar()
	return nil
}

func Get() *zap.Logger {
	if globalLogger == nil {
		globalLogger, _ = zap.NewDevelopment()
	}
	return globalLogger
}

func GetSugar() *zap.SugaredLogger {
	if globalSugar == nil {
		globalSugar = Get().Sugar()
	}
	return globalSugar
}

func WithContext(ctx context.Context) *zap.Logger {
	logger := Get()

	if corrID := GetCorrelationID(ctx); corrID != "" {
		logger = logger.With(zap.String("correlation_id", corrID))
	}

	if reqID := GetRequestID(ctx); reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return ctxLogger
	}

	return logger
}

func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func SetCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationKey, id)
}

func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(correlationKey).(string); ok {
		return id
	}
	return ""
}

func SetRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}
