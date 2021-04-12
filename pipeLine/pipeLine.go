package pipeLine

import (
	"log"
	"net/http"
	"strings"
	"ws/common"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc
var MiddlewareRequest map[string]string

func init() {
	MiddlewareRequest=make(map[string]string)
}
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
				_, _ = w.Write([]byte(common.HelloWorld))
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
func HasName(name string) Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_name,err:=GetName(r)
			MiddlewareRequest["token1"]=_name
			if len(CheckMiddleRequest(name))==0 ||err !=nil {
				w.Header().Set("Connection","close")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(InvalidParam))
				log.Println(InvalidParam)
				return
			}
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
func CheckMiddleRequest(key string) string {
	if name,ok:=MiddlewareRequest[key];ok {
		return name
	}
	return ""
}
