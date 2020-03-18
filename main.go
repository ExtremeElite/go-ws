package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"ws/conf"
	"ws/impl"
)

func main(){
	impl.HttpChan=make(chan []byte,1)
	var wsPort= conf.GetConfig().WsPort
	var httpPort=conf.GetConfig().HttpPort
	go func() {
		httpPush:=http.NewServeMux()
		httpPush.HandleFunc("/",impl.HttpHandle)
		if err:=http.ListenAndServe(":"+strconv.Itoa(int(httpPort)), httpPush);err!=nil{
			log.Fatal(err)
		}
	}()
	go func() {
		for{
			select {
			case data:=<-impl.HttpChan:
				fmt.Println(string(data))
			}
		}
	}()
	wsPush:=http.NewServeMux()
	wsPush.HandleFunc("/", impl.WsHandle)
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(wsPort)), wsPush);err!=nil{
		log.Fatal(err)
	}
}
