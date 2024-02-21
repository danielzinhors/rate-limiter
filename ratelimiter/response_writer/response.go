package response_writer

import "net/http"

type RateLimiterResponseWriter interface {
	WriteResponse(w *http.ResponseWriter) error
	WriteError(w *http.ResponseWriter, err error) error
}

type rateLimiterDefaultResponseWriter struct {
	statusCode int
	message    string
}

func NewRateLimiterDefaultResponseWriter() *rateLimiterDefaultResponseWriter {
	responseWriter := &rateLimiterDefaultResponseWriter{}
	responseWriter.statusCode = 429
	responseWriter.message = "you have reached the maximum number of requests or actions allowed within a certain time frame"
	return responseWriter
}

func (rw *rateLimiterDefaultResponseWriter) WriteResponse(w *http.ResponseWriter) error {
	(*w).WriteHeader(rw.statusCode)
	(*w).Write([]byte(rw.message))
	return nil
}

func (rw *rateLimiterDefaultResponseWriter) WriteError(w *http.ResponseWriter, err error) error {
	(*w).WriteHeader(500)
	(*w).Write([]byte("internal server error"))
	return nil
}
