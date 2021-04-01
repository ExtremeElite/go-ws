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
	"ws/common"
	"ws/db"
	"ws/util"
)

type Login struct {
	User string
	Pwd  string
}

var (
	body              []byte
	login             Login
	TokenUnauthorized = "token失效"
	NoParam           = "未获取到参数"
	InvalidParam      = "请传入正确的参数"
)

//ws登录
func wsAuth(r *http.Request) (name string, err error) {
	response := util.Response{}
	name, err = GetName(r)
	if err != nil {
		return
	}
	if !validateToken(name) {
		err = errors.New(string(response.Json(TokenUnauthorized, http.StatusUnauthorized, "")))
	}
	return
}

//http登录
func HttpAuth(r *http.Request) (data string, err error) {
	if r.Method == http.MethodGet {
		data = common.HelloWorld
	} else {
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			return
		}
		if err = json.Unmarshal(body, &login); err != nil {
			return
		}
	}
	return
}

func GetName(r *http.Request) (name string, err error) {
	response := util.Response{}
	query := r.URL.Query()
	if len(query) == 0 {
		err = errors.New(string(response.Json(NoParam, http.StatusUnauthorized, "")))
		return
	}
	if token, ok := query["token"]; ok {
		name = token[0]
		return
	}
	if token, ok := query["sn"]; ok {
		name = token[0]
		return
	}
	err = errors.New(string(response.Json(InvalidParam, http.StatusUnauthorized, "")))
	return
}

func validateToken(token string) (ok bool) {
	return true
	var sql = `select count(*) from doorplate where sn = ?`
	var total int
	db.DB.Raw(sql, token).Scan(&total)
	if total >= 1 {
		return true
	}
	return false
}

func WsAuthMiddle() Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_, err := wsAuth(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
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
