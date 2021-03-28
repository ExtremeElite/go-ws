package main

import (
	"fmt"
	"github.com/mbndr/figlet4go"
	"log"
	"net/http"
	"strconv"
	"time"
	"ws/broker"
	"ws/router"

	"ws/common"
)

func init() {
	broker.HttpChan = make(chan broker.PushData, 1)
	Logo()
}
func httpPush() {
	var httpPort = common.Setting.HttpPort
	var httpTimeOut = common.Setting.HttpTimeOut
	httpPush := http.NewServeMux()

	httpPush.HandleFunc("/", router.HttpRouter())
	httpPushTimeOut := http.TimeoutHandler(httpPush, time.Duration(httpTimeOut)*time.Second, "请求超时")
	log.Printf("http服务器0.0.0.0:%d", httpPort)
	if err := http.ListenAndServe(":"+strconv.Itoa(int(httpPort)), httpPushTimeOut); err != nil {
		log.Fatal(err)
	}
}
func wsPush() {
	var wsPort = common.Setting.WsPort
	wsPush := http.NewServeMux()
	wsPush.HandleFunc("/", router.WsRouter())
	log.Printf("ws服务器0.0.0.0:%d", wsPort)
	if err := http.ListenAndServe(":"+strconv.Itoa(int(wsPort)), wsPush); err != nil {
		log.Fatal("main:", err)
	}
}
func Logo() {
	ascii := figlet4go.NewAsciiRender()
	// Adding the colors to RenderOptions
	options := figlet4go.NewRenderOptions()
	renderStr, _ := ascii.RenderOpts("Gorouting", options)
	fmt.Println(renderStr)
}
func main() {
	go httpPush()
	go broker.HttpMessageForwarding()
	wsPush()
}
