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

	mux.Handle("/", withMiddlewares(
		http.HandlerFunc(HandleWS),
		cors(common.Conf.CorsOrigins),
		logging(),
	))

	mux.Handle("/all", withMiddlewares(
		http.HandlerFunc(HandleAllNodes),
		cors(common.Conf.CorsOrigins),
	))

	mux.Handle("/panic", withMiddlewares(
		http.HandlerFunc(HandlePanicStats),
		cors(common.Conf.CorsOrigins),
		method("GET"),
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

	mux.Handle("/revoke", withMiddlewares(
		http.HandlerFunc(HandleRevokeToken),
		method("POST"),
		localOnly(),
	))

	h := http.TimeoutHandler(mux, time.Duration(timeout)*time.Second, "timeout")
	log.Printf("[success] http server on :%d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(port)), h))
}
