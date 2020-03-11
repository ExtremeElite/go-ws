package main

import (
	"log"
	"net/http"
	"strconv"
	"ws/conf"
	"ws/impl"
)

func main(){
	var wsPort= conf.GetConfig().WsPort
	http.HandleFunc("/", impl.WsHandle)
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(wsPort)), nil);err!=nil{
		log.Fatal(err)
	}
}
