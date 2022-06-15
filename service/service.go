package service

import (
	"net/http"
	"ws/pipeLine"
	"ws/util"
)

// package service 为对外开发的接口 配合package router 对业务逻辑进行开发
func HttpGetTokenHandle(w http.ResponseWriter, r *http.Request) {
	var response = util.Response{}
	w.Write(response.Json("success", 200, pipeLine.CreateToken(1, 1)))
}
