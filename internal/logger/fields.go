package logger

import (
	"time"

	"go.uber.org/zap"
)

const (
	FieldMethod       = "method"
	FieldPath         = "path"
	FieldStatusCode   = "status_code"
	FieldLatency      = "latency_ms"
	FieldClientIP     = "client_ip"
	FieldUserAgent    = "user_agent"
	FieldRequestSize  = "request_size"
	FieldResponseSize = "response_size"
	FieldError        = "error"
	FieldErrorType    = "error_type"
	FieldGRPCService  = "grpc_service"
	FieldGRPCMethod   = "grpc_method"
	FieldGRPCCode     = "grpc_code"
)

type Fields []zap.Field

func NewFields() Fields {
	return Fields{}
}

func (f Fields) Add(fields ...zap.Field) Fields {
	return append(f, fields...)
}

func (f Fields) With(fields ...zap.Field) Fields {
	return f.Add(fields...)
}

func (f Fields) Error(err error) Fields {
	if err != nil {
		return f.Add(
			zap.Error(err),
			zap.String(FieldErrorType, getErrorType(err)),
		)
	}
	return f
}

func (f Fields) Duration(key string, d time.Duration) Fields {
	return f.Add(zap.Float64(key, float64(d.Milliseconds())))
}

func getErrorType(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

func Duration(key string, val time.Duration) zap.Field {
	return zap.Duration(key, val)
}

func Time(key string, val time.Time) zap.Field {
	return zap.Time(key, val)
}

func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}
