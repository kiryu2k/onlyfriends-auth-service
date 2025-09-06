package grpc

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type validatable interface {
	Validate() error
}

func ValidationMiddleware(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp any, err error) {
	request, ok := req.(validatable)
	if !ok {
		return handler(ctx, req)
	}

	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.WithMessage(err, "validate request body").Error())
	}

	return handler(ctx, req)
}

func LoggingMiddleware(logger *zap.Logger) grpc.UnaryServerInterceptor {
	fn := func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		logger := logger.WithOptions(
			zap.AddCallerSkip(1),
		).With(convertLogFields(fields)...)

		switch level {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		default:
			logger.Error(msg)
		}
	}

	return selector.UnaryServerInterceptor(
		logging.UnaryServerInterceptor(logging.LoggerFunc(fn)),
		selector.MatchFunc(skipHealthcheck),
	)
}

func convertLogFields(fields []any) []zap.Field {
	result := make([]zap.Field, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		key, _ := fields[i].(string)
		value := fields[i+1]

		result = append(result, fieldByType(key, value))
	}

	return result
}

func fieldByType(key string, value any) zap.Field {
	switch v := value.(type) {
	case string:
		return zap.String(key, v)
	case int:
		return zap.Int(key, v)
	case bool:
		return zap.Bool(key, v)
	default:
		return zap.Any(key, v)
	}
}

func skipHealthcheck(_ context.Context, meta interceptors.CallMeta) bool {
	return meta.FullMethod() != healthpb.Health_Check_FullMethodName
}
