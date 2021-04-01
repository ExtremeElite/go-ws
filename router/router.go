package router

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"ws/common"
)

func WsPush()  {
	var wsPort = common.Setting.WsPort
	wsPush := http.NewServeMux()
	wsPush.HandleFunc("/", WsRouter())
	wsPush.HandleFunc("/all", AllNodeRouter())
	log.Printf("ws服务器0.0.0.0:%d", wsPort)
	if err := http.ListenAndServe(":"+strconv.Itoa(int(wsPort)), wsPush); err != nil {
		log.Fatal("main:", err)
	}
}
func HttpPush()  {
	var httpPort = common.Setting.HttpPort
	var httpTimeOut = common.Setting.HttpTimeOut
	httpPush := http.NewServeMux()
	httpPush.HandleFunc("/", HttpRouter())
	httpPushTimeOut := http.TimeoutHandler(httpPush, time.Duration(httpTimeOut)*time.Second, "请求超时")
	log.Printf("http服务器0.0.0.0:%d", httpPort)
	if err := http.ListenAndServe(":"+strconv.Itoa(int(httpPort)), httpPushTimeOut); err != nil {
		log.Fatal(err)
	}
}