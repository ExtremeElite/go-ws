package server

import (
	"log"
	"net/http"
	"strings"
	"time"
	"ws/auth"
)

type Middleware func(http.Handler) http.Handler

func withMiddlewares(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// --- built-in middleware ---

func method(m string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, `{"code":405,"msg":"method not allowed"}`, http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func cors() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
			next.ServeHTTP(w, r)
		})
	}
}

func localOnly() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := strings.Split(r.Host, ":")[0]
			if host != "127.0.0.1" && host != "localhost" {
				http.Error(w, `{"code":406,"msg":"host not allowed"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func maxBody(n int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, int64(n))
			next.ServeHTTP(w, r)
		})
	}
}

func authMold(moldFn func() int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get("token")
			if token == "" {
				token = r.URL.Query().Get("sn")
			}
			mold := moldFn()
			if mold != 0 && !auth.ValidateToken(token, mold) {
				http.Error(w, `{"code":401,"msg":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func logging() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
			next.ServeHTTP(w, r) // ← handler 在这里执行
		})
	}
}

// 前置 + 后置中间件示例：埋点耗时
func timing() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r) // ← handler 执行

			// ── 这里是 handler 之后 ──
			log.Printf("[timing] %s %s took %v", r.Method, r.URL.Path, time.Since(start))
		})
	}
}
