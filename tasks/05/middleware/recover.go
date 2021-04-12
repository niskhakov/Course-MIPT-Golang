package middleware

import (
	"log"
	"net/http"
)

func Recover(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Printf("[ERROR] Panic caught: %v\n", err)
					rw.WriteHeader(http.StatusInternalServerError)
					rw.Write([]byte("Internal Server Error\n"))
				}
			}()
			next.ServeHTTP(rw, req)
		}

		return http.HandlerFunc(fn)
	}
}
