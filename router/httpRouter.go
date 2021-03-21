/**
 * @date:2020\11\29 0029 23:00
 * @email:gorouting@qq.com
 * @author:gorouting
 * @description:
**/
package router

import (
	"net/http"
	"ws/broker"
	"ws/pipeLine"
)

func HttpRouter() http.HandlerFunc {
	return pipeLine.Use(
		broker.HttpHandle,
		pipeLine.Logging(),
		pipeLine.Method("POST"),
		pipeLine.HttpAuthMiddle(),
	)
}
