package pipeLine

import (
	"log"
	"net/http"
	"strings"
	"ws/util"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

var MiddlewareRequest map[string]string

func init() {
	MiddlewareRequest = make(map[string]string)
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
				_, _ = w.Write([]byte(util.MethodNotAllowed))
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
			_name, err := GetName(r)
			MiddlewareRequest["token"] = _name
			if len(CheckMiddleRequest(name)) == 0 || err != nil {
				w.Header().Set("Connection", "close")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(InvalidParam))
				log.Println(InvalidParam)
				return
			}
			fn(w, r)
		}
	}
}

// 		   	  requests
//            	 |
//            	 v
// +----------  one ---------+
// | +--------  two -------+ |
// | | +------ three -----+| |
// | | |                 | | |
// | | |       handler   | | |
// | | |                 | | |
// | | +----  after_one -+ | |
// | +------  after_two ---+ |
// +-------- after_three --+ |
//            |
//            v
//         responses

// handler:=pipLine.Use(handler)
// before := pipeLine.Before(handler,one,two,three)
// after:=pipeLine.After(before,after_one,after_two,after_three)
// return after
func Use(fn http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	var middlewareLen = len(middlewares) - 1
	for key := range middlewares {
		fn = middlewares[middlewareLen-key](fn)
	}
	return fn
}

// 前置中间件
func Before(fn http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	var middlewareLen = len(middlewares) - 1
	for key := range middlewares {
		fn = middlewares[middlewareLen-key](fn)
	}
	return fn
}

// 后置中间件
func After(fn http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	var before = fn
	fn = func(w http.ResponseWriter, r *http.Request) {}
	var middlewareLen = len(middlewares) - 1
	for key := range middlewares {
		fn = middlewares[middlewareLen-key](fn)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		before(w, r)
		fn(w, r)
	}
}

//			 requests
//	          	|
//	         	v
//
// +---------- three---------+
// | +-------- two -------+ |
// | | +------ one  -----+| |
// | | |                 | | |
// | | |      handler    | | |
// | | |                 | | |
// | | +---- after_one  -+ | |
// | +------ after_two  ---+ |
// +-------- after_three --+ |
//
//	   			|
//	   			v
//			responses
//
// handler := pipeLine.Next(handler,one,two,three,after_one,after_two,after_three)
// return handler
func Next(fn http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for key := range middlewares {
		fn = middlewares[key](fn)
	}
	return fn
}
func CheckMiddleRequest(key string) string {
	if name, ok := MiddlewareRequest[key]; ok {
		return name
	}
	return ""
}

// 是否为指定域名请求
func LocalRequest(hosts []string) Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for host := range hosts {
				if r.Host != hosts[host] {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(util.HostNotAllowed))
					return
				}
			}
			fn(w, r)
		}
	}
}
