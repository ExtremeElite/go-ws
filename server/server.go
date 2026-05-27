package server

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"ws/common"
)

func WsPush() {
	port := common.Conf.WebSocket.WsPort
	mux := http.NewServeMux()

	// 每个路由独立选择中间件：启用/停用
	mux.Handle("/", withMiddlewares(
		http.HandlerFunc(HandleWS),
		cors(),
		logging(),
		// authMold(func() int { return common.VM.WebSocket.Mold }),
	))

	mux.Handle("/all", withMiddlewares(
		http.HandlerFunc(HandleAllNodes),
		cors(),
		// 不加 auth：公开接口
	))

	log.Printf("[success] ws server on :%d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(port)), mux))
}

func HttpPush() {
	port := common.Conf.Http.HttpPort
	timeout := common.Conf.Http.HttpTimeOut
	mux := http.NewServeMux()

	mux.Handle("/", withMiddlewares(
		http.HandlerFunc(HandlePush),
		method("POST"),
		maxBody(common.Conf.MaxBody),
	))

	mux.Handle("/token", withMiddlewares(
		http.HandlerFunc(HandleToken),
		method("GET"),
		localOnly(),
	))

	h := http.TimeoutHandler(mux, time.Duration(timeout)*time.Second, "timeout")
	log.Printf("[success] http server on :%d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(port)), h))
}
