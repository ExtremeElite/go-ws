/**
 * @date:2020/8/19 16:50
 * @email:gorouting@qq.com
 * @author:gorouting
 * @description:
**/
package request

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"ws/core"
	"ws/middleware/auth"
)
var upgrader=websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func WsRequest(w http.ResponseWriter,r *http.Request)(conn *core.Connection,name string,err error)  {
	var(
		wsConn *websocket.Conn
	)
	if wsConn, err = upgrader.Upgrade(w, r, nil);err != nil{
		log.Print("upgrade:", err)
		return
	}
	if conn,err= core.BuildConn(wsConn);err!=nil {
		conn.WsConn.WriteMessage(websocket.TextMessage,[]byte(err.Error()))
		conn.Close()
		return
	}
	//登录判断
	if name,err=auth.WsAuth(r);err!=nil {
		log.Println(err.Error())
		conn.WsConn.WriteMessage(websocket.TextMessage,[]byte(err.Error()))
		conn.Close()
	}
	return
}
