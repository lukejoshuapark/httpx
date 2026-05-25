package httpx

import (
	"log/slog"
	"net/http"
)

func RequestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		iw := NewInstrumentedResponseWriter(w)
		next.ServeHTTP(iw, r)

		level := levelForStatusCode(iw.StatusCode())
		slog.Log(r.Context(), level, "Request", "status", iw.StatusCode(), "elapsed", iw.Elapsed(), "method", r.Method, "url", r.URL.String())
	})
}

func levelForStatusCode(statusCode int) slog.Level {
	switch {
	case statusCode >= 500:
		return slog.LevelError
	case statusCode >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
