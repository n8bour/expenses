package middleware

import (
	"log"
	"net/http"
	"time"
)

type WrapResponseWriter interface {
	http.ResponseWriter
	// Status returns the HTTP status of the request, or 0 if one has not
	// yet been sent.
	Status() int
	// BytesWritten returns the total number of bytes sent to the client.
	BytesWritten() int
}

type basicWriter struct {
	http.ResponseWriter
	code  int
	bytes int
}

// SimpleLogger a simpler logger middleware that logs the requests into the server
func SimpleLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			bw := &basicWriter{ResponseWriter: w}
			t1 := time.Now()
			defer func() {
				log.Printf("%s: \"%s%s\"  %s - %d | %dBytes in %s", r.Method, r.Host, r.RequestURI, r.Proto, bw.code, bw.bytes, time.Since(t1))
			}()

			next.ServeHTTP(bw, r)
		})
}

func (b *basicWriter) WriteHeader(code int) {
	b.code = code
	b.ResponseWriter.WriteHeader(code)

}

func (b *basicWriter) Write(buf []byte) (int, error) {
	n, err := b.ResponseWriter.Write(buf)
	b.bytes += n
	return n, err
}
