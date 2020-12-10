package main

import (
	"fmt"
	"github.com/mbndr/figlet4go"
	"log"
	"net/http"
	"strconv"
	"time"
	"ws/broker"
	"ws/conf"
	"ws/router"
)

func init()  {
	broker.HttpChan =make(chan broker.PushData,1)
	Logo()
}
func main() {
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

func Logo()  {
	ascii := figlet4go.NewAsciiRender()
	// Adding the colors to RenderOptions
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		// Colors can be given by default ansi color codes...
		figlet4go.ColorGreen,
		figlet4go.ColorYellow,
		figlet4go.ColorCyan,
	}

	renderStr, _ := ascii.RenderOpts("Gorouting", options)
	fmt.Println(renderStr)
}

