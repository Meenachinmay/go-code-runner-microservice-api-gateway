package middleware

import (
	"context"
	"time"

	"go-code-runner-microservice/api-gateway/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryClientLoggingInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()

		log := logger.WithContext(ctx)

		correlationID := logger.GetCorrelationID(ctx)
		if correlationID != "" {
			md := metadata.New(map[string]string{
				"correlation-id": correlationID,
			})
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		log.Info("grpc_client_started",
			zap.String(logger.FieldGRPCMethod, method),
			zap.String("target", cc.Target()),
		)

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start)

		fields := logger.NewFields().
			With(
				zap.String(logger.FieldGRPCMethod, method),
				zap.Duration(logger.FieldLatency, duration),
			)

		if err != nil {
			st, _ := status.FromError(err)
			fields = fields.With(
				zap.String(logger.FieldGRPCCode, st.Code().String()),
				zap.Error(err),
			)

			if st.Code() == codes.Unavailable || st.Code() == codes.DeadlineExceeded || st.Code() == codes.Canceled || st.Code() == codes.Unknown || st.Code() == codes.Internal || st.Code() == codes.ResourceExhausted {
				log.Error("grpc_client_failed", fields...)
			} else {
				log.Warn("grpc_client_error", fields...)
			}
		} else {
			fields = fields.With(zap.String(logger.FieldGRPCCode, codes.OK.String()))
			log.Info("grpc_client_completed", fields...)
		}

		if duration > 500*time.Millisecond {
			log.Warn("slow_grpc_call_detected",
				zap.String("method", method),
				zap.Duration("duration", duration),
			)
		}

		return err
	}
}
