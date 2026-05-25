package httpx

import (
	"net/http"
	"time"
)

type InstrumentedResponseWriter struct {
	w http.ResponseWriter

	created    time.Time
	statusCode *int
}

var _ http.ResponseWriter = (*InstrumentedResponseWriter)(nil)

func NewInstrumentedResponseWriter(w http.ResponseWriter) *InstrumentedResponseWriter {
	return &InstrumentedResponseWriter{
		w: w,

		created: time.Now(),
	}
}

func (w *InstrumentedResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *InstrumentedResponseWriter) Write(b []byte) (int, error) {
	return w.w.Write(b)
}

func (w *InstrumentedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = &statusCode
	w.w.WriteHeader(statusCode)
}

func (w *InstrumentedResponseWriter) Elapsed() time.Duration {
	return time.Since(w.created)
}

func (w *InstrumentedResponseWriter) StatusCode() int {
	if w.statusCode == nil {
		return http.StatusOK
	}

	return *w.statusCode
}
