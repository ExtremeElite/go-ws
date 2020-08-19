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
				http.Error(w,`<h1>欢迎来到及时服务</h1>`, http.StatusBadRequest)
				return
			}
			fn(w, r)
		}
	}
}
func Cors() Middleware  {
	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request){
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
		}
	}
}
func Use(fn http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		fn = middleware(fn)
	}
	return fn
}