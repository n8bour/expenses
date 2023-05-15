package middleware

import (
	"log"
	"net/http"
	"time"
)

func SimpleLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%v] -----> %s: \"%s%s\"  %s", time.Now().Local().Truncate(time.Second), r.Method, r.Host, r.URL.Path, r.Proto)

		next.ServeHTTP(w, r)
	})
}
