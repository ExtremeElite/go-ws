package broker

import (
	"log"
	"net/http"
)
var HttpChan chan PushData
func HttpHandle(w http.ResponseWriter, r *http.Request)  {
	err:=httpBroker(w,r)
	if err!=nil{
		log.Println("http error",err.Error())
		return
	}
}