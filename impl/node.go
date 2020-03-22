package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Node struct {
	Ws *Connection
	Name string
}
var NodeList map[string]Node
var mut sync.Mutex

func AddNode(node *Node) (err error){
	mut.Lock()
	if ok:=isExist(node);ok{
		DelNode(node)
	}
	NodeList[node.Name]=*node
	mut.Unlock()
	return
}
func DelNode(node *Node)(err error){
	mut.Lock()
	if ok:=isExist(node);ok{
		node.Ws.WsConn.WriteMessage(websocket.TextMessage,[]byte(`你的连接已经断开了`))
		err=errors.New(`你的连接已经断开了`)
		node.Ws.Close()
		delete(NodeList,node.Name)
	}
	mut.Unlock()
	return
}
func isExist(node *Node) (ok bool){
	mut.Lock()
	if _,ok=NodeList[node.Name];ok{
	}
	mut.Unlock()
	return
}
