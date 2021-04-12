package broker

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
	"ws/common"
	"ws/kernel"
	"ws/pipeLine"
	"ws/util"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandle(w http.ResponseWriter, r *http.Request) {

	var (
		err  error
		conn *kernel.Connection
		name string
	)
	//普通 HTTP请求
	if r.Header.Get("Connection") != "Upgrade" {
		if _, err := w.Write([]byte(common.HelloWorld)); err != nil && common.Debug {
			log.Println("http error: ", err.Error())
		}
		return
	}
	if conn, name, err = wsBuild(w, r); err != nil {
		log.Println("wsRequest:", err.Error())
		return
	}
	defer conn.Close()
	kernel.AddNode(&kernel.Node{Ws: conn, Name: name, RemoteAddr: r.RemoteAddr})
	for {

		if err = wsWork(conn); err != nil && !strings.Contains(err.Error(), `wsMessageForwarding`) {
			goto Err
		}
	}
Err:
	if conn.IsClose {
		if err := kernel.DelNode(name); err != nil && common.Debug {
			log.Println("connect close:", err.Error())
		}
	}

	log.Println("connect close:", name)
}

//业务逻辑处理
func wsWork(conn *kernel.Connection) (err error) {
	var wsTimeOut = common.Setting.WsTimeOut
	//设置服务器读取超时
	if wsTimeOut > 0 {
		if err = conn.SetReadDeadline(time.Now().Add(time.Duration(wsTimeOut) * time.Second)); err != nil {
			return
		}
	}
	return wsBroker(conn)
}

//创建连接
func wsBuild(w http.ResponseWriter, r *http.Request) (conn *kernel.Connection, name string, err error) {
	var (
		wsConn *websocket.Conn
	)
	if wsConn, err = upgrade.Upgrade(w, r, nil); err != nil {
		return
	}
	if conn, err = kernel.BuildConn(wsConn); err != nil {
		if wsErr := wsConn.WriteMessage(websocket.TextMessage, []byte(err.Error())); wsErr != nil && common.Debug {
			log.Println("wsErr is:", wsErr.Error())
		}
		_ = wsConn.Close()
		return
	}
	name = pipeLine.MiddlewareRequest["token"]
	if len(name)==0 {
		err=errors.New("获取到的name为空")
		return nil, "", err
	}
	var response = util.Response{}
	if err = conn.WriteMsg(response.Json("登录成功", 200, "")); err != nil {
		return nil, "", err
	}
	log.Println("connect open:", r.RemoteAddr, name)
	return
}
