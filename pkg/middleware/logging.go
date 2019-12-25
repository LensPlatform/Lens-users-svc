package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	logger *zap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Debug(
			"request started",
			zap.String("proto", r.Proto),
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.String("remote", r.RemoteAddr),
			zap.String("user-agent", r.UserAgent()),
		)

		defer func(begin time.Time) {
			m.logger.Info("Transport Log",
				zap.String("[METHOD] : ", r.Method),
				zap.String("[URL] : ", r.URL.String()),
				zap.Any("[DURATION] : ", time.Since(begin)))
		}(time.Now())

		next.ServeHTTP(w, r)
	})
}
