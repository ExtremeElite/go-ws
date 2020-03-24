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
var NodeList sync.Map
func AddNode(node *Node) (err error){
	if _,ok:=GetNode(node.Name);ok{
		DelNode(node.Name)
	}
	NodeList.Store(node.Name,node)
	return
}
func GetNode(name string)(*Node,bool){
	var node *Node
	if v,ok:=NodeList.Load(name);ok {
		node=v.(*Node)
		return node,true
	}
	return node,false
}
func DelNode(name string)(err error){
	if node,ok:=GetNode(name);ok{
		node.Ws.WsConn.WriteMessage(websocket.TextMessage,[]byte(`你的连接已经断开了`))
		err=errors.New(`你的连接已经断开了`)
		node.Ws.Close()
		NodeList.Delete(name)
	}
	return
}
