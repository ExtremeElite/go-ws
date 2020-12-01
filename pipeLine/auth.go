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
	"ws/db"
)

type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}
func (response Response)Json(msg string,code int) string{
	response.Msg=msg
	response.Code=code
	response.Data=""
	data,err:=json.Marshal(response)
	if err!=nil {
		return `["code":404,"msg":"数据错误","data":""]`
	}
	return string(data)
}
type Login struct {
	User string
	Pwd string
}
var(
	body       []byte
	login      Login
)
//ws登录
func wsAuth(r *http.Request) (name string,err error)  {
	response:=Response{}
	name,err=GetName(r)
	if err!=nil{
		return
	}
	if !validateToken(name) {
		err=errors.New(response.Json("token失效",http.StatusUnauthorized))
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

func GetName(r *http.Request) (name string,err error)  {
	response:=Response{}
	query:=r.URL.Query()
	if len(query)==0 {
		err=errors.New(response.Json("未获取到参数",http.StatusUnauthorized))
		return
	}
	if token,ok:=query["token"];ok{
		name=token[0]
		return
	}
	err=errors.New(response.Json("请传入正确的参数",http.StatusUnauthorized))
	return
}

func validateToken(token string) (ok bool)  {
	return true
	var total int
	db.DB.Raw(`select count(*) from hb_shebei where device_nums = ?`,token).Scan(&total)
	if total>=1 {
		return true
	}
	return false
}

func WsAuthMiddle() Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_,err:=wsAuth(r)
			if err!=nil {
				http.Error(w,err.Error(), http.StatusUnauthorized)
				return
			}
			fn(w, r)
		}
	}
}

func HttpAuthMiddle() Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			//time.Sleep(time.Second*3)
			fn(w, r)
		}
	}
}