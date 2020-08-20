/**
 * @date:2020/8/19 16:46
 * @email:gorouting@qq.com
 * @author:gorouting
 * @description:认证中间件
**/
package pipeLine

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"ws/core"
)


type Login struct {
	User string
	Pwd string
}
var(
	body       []byte
	dataString string
	login      Login
)
//ws登录
func WsAuth(r *http.Request) (name string,err error)  {
	query:=r.URL.Query()
	if len(query)==0 {
		err=errors.New(`{"C":"Login","M":"验证失败"}`)
	}
	if token,ok:=query["token"];ok{
		if !validateToken(token[0]) {
			err=errors.New(`{"C":"Login","M":"token失效"}`)
		}else {
			name=token[0]
		}
		return
	}else {
		err=errors.New(`{"C":"Login","M":"登陆失败"}`)
	}
	return
}
//http登录
func HttpAuth(r *http.Request)(data string,err error){
	if r.Method=="GET" {
		data=core.HELLO
	}else {
		if body,err=ioutil.ReadAll(r.Body);err!=nil{
			return
		}
		if err=json.Unmarshal(body,&login);err!=nil {
			return
		}
	}
	return
}

func validateToken(token string) (ok bool)  {

	return true
}

func WsAuthMiddle() Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fn(w, r)
		}
	}
}

func HttpAuthMiddle() Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fn(w, r)
		}
	}
}