package middleware

import (
	"net/http"
)

type MiddleWare func(http.HandlerFunc) http.HandlerFunc

func Loging() MiddleWare  {
	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			
		}
	}
}