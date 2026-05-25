package httpx

import (
	"net/http"
	"strconv"
)

func RespondWithBytes(w http.ResponseWriter, statusCode int, contentType string, bytes []byte) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	w.WriteHeader(statusCode)
	w.Write(bytes)
}
