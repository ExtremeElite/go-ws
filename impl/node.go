package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Node struct {
	Ws *Connection
	Name string
	IsOnline bool
}
var NodeList map[string]Node
var mut sync.Mutex
var onec sync.Once
func AddNode(node *Node) (err error){
	mut.Lock()
	if ok:=isExist(node.Name);ok{
		DelNode(node.Name)
	}
	NodeList[node.Name]=*node
	mut.Unlock()
	return
}
func DelNode(name string)(err error){
	if ok:=isExist(name);ok{
		NodeList[name].Ws.WsConn.WriteMessage(websocket.TextMessage,[]byte(`你的连接已经断开了`))
		err=errors.New(`你的连接已经断开了`)
		delete(NodeList,name)
	}
	return
}
func isExist(name string) (ok bool){
	if _,ok=NodeList[name];ok{
	}
	return
}
