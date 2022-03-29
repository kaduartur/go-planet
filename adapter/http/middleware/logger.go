package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Logger(l *zap.Logger, utc bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			query := r.URL.RawQuery
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			end := time.Now()
			latency := end.Sub(start)
			if utc {
				end = end.UTC()
			}
			defer func() {
				l.Info(path,
					zap.Int("status", ww.Status()),
					zap.String("method", r.Method),
					zap.String("path", path),
					zap.String("query", query),
					zap.String("user-agent", r.UserAgent()),
					zap.String("time", end.Format(time.RFC3339)),
					zap.Duration("latency", latency),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
