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
	impl.NodeList=make(map[string]impl.Node)
	go httpPush()
	go getDataFromHttp()
	wsPush()

}
func httpPush()  {
	var httpPort=conf.Config().Common.HttpPort
	httpPush:=http.NewServeMux()
	httpPush.HandleFunc("/",impl.HttpHandle)
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(httpPort)), httpPush);err!=nil{
		log.Fatal(err)
	}
}
func wsPush() {
	var wsPort= conf.Config().Common.WsPort
	wsPush:=http.NewServeMux()
	wsPush.HandleFunc("/", impl.WsHandle)
	if err:=http.ListenAndServe(":"+strconv.Itoa(int(wsPort)), wsPush);err!=nil{
		log.Fatal(err)
	}
}

func getDataFromHttp()  {
	for{
		select {
		case data:=<-impl.HttpChan:
			for name,_:=range impl.NodeList{
				impl.NodeList[name].Ws.WriteMsg(data)
			}
			fmt.Println(string(data))
		}
	}
}


