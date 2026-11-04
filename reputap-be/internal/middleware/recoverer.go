package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/kittu533/reputap-be/internal/response"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
				)
				response.InternalError(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
