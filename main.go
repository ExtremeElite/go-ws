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
	run()
}
func run(){
	impl.HttpChan=make(chan []byte,1)
	go httpPush()
	go getDataFromHttp()
	wsPush()
}
func httpPush() {
	var httpPort=conf.CommonSet.HttpPort
	httpPush:=http.NewServeMux()
	httpPush.HandleFunc("/",impl.HttpHandle)
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(httpPort)), httpPush);err!=nil{
		log.Fatal(err)
	}
}
func wsPush() {
	var wsPort= conf.CommonSet.WsPort
	wsPush:=http.NewServeMux()
	wsPush.HandleFunc("/", impl.WsHandle)
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(wsPort)), wsPush);err!=nil{
		log.Fatal("main:",err)
	}
}
func getDataFromHttp()  {
	for{
		select {
		case data:=<-impl.HttpChan:
			impl.NodeList.Range(func(name, node interface{}) bool {
				node.(*impl.Node).Ws.WriteMsg(data)
				return true
			})
			fmt.Println(string(data))
		}
	}
}
