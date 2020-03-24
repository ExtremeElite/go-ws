package impl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)
var HttpChan chan []byte
func HttpHandle(w http.ResponseWriter, r *http.Request)  {
	var (
		body []byte
		err error
	)
	if body,err=ioutil.ReadAll(r.Body);err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	select {
	case HttpChan<-body:
	default:
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	fmt.Println(time.Now().UnixNano())
	w.WriteHeader(http.StatusOK)
}