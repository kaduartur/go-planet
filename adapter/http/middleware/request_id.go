package middleware

import (
	"net/http"

	"github.com/hashicorp/go-uuid"
	"github.com/kaduartur/go-planet/pkg/log"
)

const RequestIDHeader = "X-Request-Id"

func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID, _ = uuid.GenerateUUID()
		}

		ctx := log.NewLogContext(r.Context(), map[string]interface{}{
			"X-Request-ID": requestID,
		})

		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
