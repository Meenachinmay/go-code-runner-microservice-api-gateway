package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-code-runner-microservice/api-gateway/internal/logger"
	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		start := time.Now()

		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = requestID
		}
		c.Set("correlation_id", correlationID)

		ctx := c.Request.Context()
		ctx = logger.SetRequestID(ctx, requestID)
		ctx = logger.SetCorrelationID(ctx, correlationID)
		c.Request = c.Request.WithContext(ctx)

		log := logger.WithContext(ctx).With(
			zap.String(logger.FieldMethod, c.Request.Method),
			zap.String(logger.FieldPath, c.Request.URL.Path),
			zap.String(logger.FieldClientIP, c.ClientIP()),
			zap.String(logger.FieldUserAgent, c.Request.UserAgent()),
		)

		ctx = logger.ToContext(ctx, log)
		c.Request = c.Request.WithContext(ctx)

		var requestBody []byte
		if c.Request.Body != nil && shouldLogBody(c.Request.Method, c.Request.URL.Path) {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		fields := logger.NewFields().
			With(zap.Int64(logger.FieldRequestSize, c.Request.ContentLength))

		if len(requestBody) > 0 && len(requestBody) < 1000 { // Only log small bodies
			fields = fields.With(zap.String("request_body", string(requestBody)))
		}

		log.Info("request_started", fields...)

		c.Next()

		latency := time.Since(start)

		responseFields := logger.NewFields().
			With(
				zap.Int(logger.FieldStatusCode, c.Writer.Status()),
				zap.Duration(logger.FieldLatency, latency),
				zap.Int(logger.FieldResponseSize, blw.body.Len()),
			)

		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Writer.Header().Set("X-Correlation-ID", correlationID)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				responseFields = responseFields.Error(e.Err)
			}
		}

		switch {
		case c.Writer.Status() >= 500:
			log.Error("request_failed", responseFields...)
		case c.Writer.Status() >= 400:
			log.Warn("request_error", responseFields...)
		case c.Writer.Status() >= 300:
			log.Info("request_redirect", responseFields...)
		default:
			log.Info("request_completed", responseFields...)
		}

		if latency > time.Second {
			log.Warn("slow_request_detected",
				zap.Duration("latency", latency),
				zap.String("path", c.Request.URL.Path),
			)
		}
	}
}

func shouldLogBody(method, path string) bool {
	skipPaths := []string{
		"/api/v1/companies/login",
		"/api/v1/companies/register",
	}

	for _, skipPath := range skipPaths {
		if path == skipPath {
			return false
		}
	}

	return method == "POST" || method == "PUT" || method == "PATCH"
}

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log := logger.WithContext(c.Request.Context())
				log.Error("panic_recovered",
					zap.Any("error", err),
					zap.Stack("stack"),
				)

				c.JSON(500, gin.H{
					"success":    false,
					"error":      "Internal server error",
					"request_id": logger.GetRequestID(c.Request.Context()),
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}
