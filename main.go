package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"ws/broker"
	"ws/conf"
	"ws/db"
	"ws/router"
)

func init()  {
	broker.HttpChan =make(chan broker.PushData,1)
}
func main() {
	defer func() {
		rawDB,err:=db.DB.DB()
		rawDB.Close()
		log.Println(err.Error())
	}()
	go httpPush()
	go broker.HttpMessageForwarding()
	wsPush()
}
func httpPush() {
	var httpPort=conf.CommonSet.HttpPort
	var httpTimeOut=conf.CommonSet.HttpTimeOut
	httpPush:=http.NewServeMux()

	httpPush.HandleFunc("/", router.HttpRouter())
	httpPushTimeOut:=http.TimeoutHandler(httpPush,time.Duration(httpTimeOut)*time.Second,"请求超时")
	log.Printf("http服务器127.0.0.1:%d",httpPort)
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(httpPort)), httpPushTimeOut);err!=nil{
		log.Fatal(err)
	}
}
func wsPush() {
	var wsPort= conf.CommonSet.WsPort
	wsPush:=http.NewServeMux()
	wsPush.HandleFunc("/", router.WsRouter())
	log.Printf("ws服务器127.0.0.1:%d",wsPort)
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(wsPort)), wsPush);err!=nil{
		log.Fatal("main:",err)
	}
}


