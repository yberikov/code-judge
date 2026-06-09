package middleware

import (
	"code-judge/pkg/logger"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func LoggerMiddleware(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = uuid.New().String()
			}
			reqLogger := base.With("request_id", requestID)
			ctx := logger.WithContext(r.Context(), reqLogger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
