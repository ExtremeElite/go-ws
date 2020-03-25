package middleware

import (
	"net/http"
	"ws/middleware/httpLog"
)

type MiddleWare func(http.HandlerFunc) http.HandlerFunc

func Loging() MiddleWare  {
	return httpLog.Log()
}