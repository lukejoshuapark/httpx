package httpx

import "net/http"

func MaximumRequestBodySizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, config.MaximumRequestBodySize())
		next.ServeHTTP(w, r)
	})
}
