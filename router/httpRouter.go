/*
date:2020\11\29 0029 23:00
email:gorouting@qq.com
author:gorouting
description:
*/
package router

import (
	"net/http"
	"ws/broker"
	"ws/pipeLine"
	"ws/service"
)

// http常规请求
func HttpRouter() http.HandlerFunc {
	return pipeLine.Next(
		broker.HttpHandle,
		pipeLine.Logging(),
		pipeLine.Method("POST"),
		pipeLine.HttpAuthMiddle(),
	)
}

// 内网获取token
func HttpGetToken() http.HandlerFunc {
	return pipeLine.Next(
		service.HttpGetTokenHandle,
		pipeLine.Logging(),
		pipeLine.Method("get"),
		pipeLine.LocalRequest([]string{"127.0.0.1"}),
	)
}
