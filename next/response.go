package next

import (
	"net/http"
)

// ResponseWriter is a response wrapper
type ResponseWriter struct {
	http.ResponseWriter
	writeCodeListener func(code int)
}

// NewResponseWriter create a new ResposneWriter wrapper
func NewResponseWriter(res http.ResponseWriter, writeCodeListener func(code int)) *ResponseWriter {
	return &ResponseWriter{res, writeCodeListener}
}

// Header atisfy the http.ResponseWriter interface
func (w ResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write send content to response
func (w ResponseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}

// WriteHeader send a response code to client
func (w ResponseWriter) WriteHeader(statusCode int) {
	w.writeCodeListener(statusCode)
	w.ResponseWriter.WriteHeader(statusCode)
}
