package pipeLine

import (
	"net/http"
	"strings"
	"ws/util"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Logging() Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fn(w, r)
		}
	}
}
func Method(m string) Middleware {
	m = strings.ToUpper(m)
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(util.HELLO))
				return
			}
			fn(w, r)
		}
	}
}
func Cors() Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
			fn(w, r)
		}
	}
}
func Use(fn http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	var middlewareLen = len(middlewares) - 1
	for key, _ := range middlewares {
		fn = middlewares[middlewareLen-key](fn)
	}
	return fn
}
