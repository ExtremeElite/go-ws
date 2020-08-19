package core

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"ws/conf"
)
var HttpChan chan []byte
func HttpHandle(w http.ResponseWriter, r *http.Request)  {
	var (
		body []byte
		err error
	)
	if r.Method=="GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(HELLO))
		return
	}
	if body,err=ioutil.ReadAll(r.Body);err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if bodyLen:=len(body);bodyLen>conf.CommonSet.MaxBody {
		res:=`请求体大小为`+strconv.Itoa(bodyLen/1024)+`kb,大于`+strconv.Itoa(conf.CommonSet.MaxBody/1024)+`kb`
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte(res))
		log.Println(res)
		return
	}
	select {
	case HttpChan<-body:
	default:
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Println(r.RemoteAddr+"发来的消息:"+string(body))
	w.Write([]byte(`{contentLength:`+strconv.Itoa((len(body)))+`b}`))
}