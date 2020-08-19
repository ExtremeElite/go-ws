package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc
func Logging() Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				log.Println(r.URL.Path, time.Since(start))
			}()
			fn(w, r)
		}
	}
}
func Method(m string) Middleware {
	m=strings.ToUpper(m)
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			fn(w, r)
		}
	}
}
func Use(fn http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		fn = middleware(fn)
	}
	return fn
}