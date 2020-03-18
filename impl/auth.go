package impl

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	HELLO=`<h1>欢迎来到Gorouting即时通讯服务</h1>`
)
var(
	pongNum uint8
	data []byte
	err error
	dataString string
	currentTime int64
)
//ws登录
func WsAuth(r *http.Request) (err error)  {
	query:=r.URL.Query()
	if len(query)==0 {
		err=errors.New(`{"M":"checkinok","ID":"xx1","NAME":"xx2","T":"xx3"}`)
	}
	if token,ok:=query["token"];ok{
		if !validateToken(token[0]) {
			err=errors.New("token过期")
		}
		return
	}else {
		err=errors.New(`{"M":""}`)
	}
	return
}
//http登录
func HttpAuth(r *http.Request)(data string,err error){
	type Login struct {
		User string
		Pwd string
	}
	var(
		login Login
		body []byte
	)
	if r.Method=="GET" {
		data=HELLO
	}else {
		body, err= ioutil.ReadAll(r.Body)
		if err!=nil {
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

