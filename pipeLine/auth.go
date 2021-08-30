/*
date:2020/8/19 16:46
email:gorouting@qq.com
author:gorouting
description:认证中间件
*/
package pipeLine

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	InvalidParam      = "请传入正确的参数"
	AuthToken         = []string{"token", "sn"}
)

//ws登录
func wsAuth() (name string, err error) {
	response := util.Response{}
	name = MiddlewareRequest["token"]
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
	_err := string(response.Json(InvalidParam, http.StatusUnauthorized, ""))
	query := r.URL.Query()
	if len(query) == 0 {
		err = errors.New(_err)
		return
	}
	if name, ok := hasToken(query); ok {
		return name, err
	}
	if len(name) == 0 {
		err = errors.New(_err)
	}
	MiddlewareRequest["token"] = name
	return
}

func validateToken(token string) (ok bool) {
	defer func() {
		if common.Debug {
			return
		}
		if err := recover(); err != nil {
			log.Println("validate token sql query: ", err)
		}
	}()
	MiddlewareRequest["token"] = token
	var sql = `select count(*) from doorplate where sn = ?`
	var total int
	db.DB.Raw(sql, token).Scan(&total)
	if total >= 1 {
		MiddlewareRequest["token"] = token
		return true
	}
	return false
}

func WsAuthMiddle() Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_, err := wsAuth()
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(TokenUnauthorized))
				log.Println("认证失败")
				return
			}
			println("ws")
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

func hasToken(query url.Values) (token string, ok bool) {
	for _, _token := range AuthToken {
		token = query.Get(_token)
		if len(token) > 0 {
			return token, true
		}
	}
	return "", false
}
