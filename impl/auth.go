package impl

import (
	"errors"
	"net/http"
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
func HttpAuth(r *http.Request)(err error){

	return
}

func validateToken(token string) (ok bool)  {

	return
}

