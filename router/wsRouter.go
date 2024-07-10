/*
date:2020\11\29 0029 22:58
email:gorouting@qq.com
author:gorouting
description:
*/
package router

import (
	"net/http"

	"ws/broker"
	"ws/pipeLine"
)

func WsRouter() http.HandlerFunc {
	var before = pipeLine.Before(
		broker.WsHandle,
		pipeLine.Cors(),
		pipeLine.HasName("token"),
	//pipeLine.WsAuthMiddle()
	)
	return before
}
