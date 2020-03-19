package main

import (
	"log"
	"net/http"
	"strconv"
	"ws/conf"
	"ws/impl"
)

func main(){
	impl.HttpChan=make(chan []byte,1)
	var wsPort= conf.GetConfig().Common.WsPort
	var httpPort=conf.GetConfig().Common.HttpPort
	go func() {
		httpPush:=http.NewServeMux()
		httpPush.HandleFunc("/",impl.HttpHandle)
		if err:=http.ListenAndServe(":"+strconv.Itoa(int(httpPort)), httpPush);err!=nil{
			log.Fatal(err)
		}
	}()
	wsPush:=http.NewServeMux()
	wsPush.HandleFunc("/", impl.WsHandle)
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(wsPort)), wsPush);err!=nil{
		log.Fatal(err)
	}
}
