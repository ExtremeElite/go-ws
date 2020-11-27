package core

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Node struct {
	Ws *Connection
	Name string
}
var Nodes sync.Map
func AddNode(node *Node){
	if _,ok:=GetNode(node.Name);ok{
		DelNode(node.Name)
	}
	Nodes.Store(node.Name,node)
}
func GetNode(name string)(*Node,bool){
	var node *Node
	if v,ok:=Nodes.Load(name);ok {
		node=v.(*Node)
		return node,true
	}
	return node,false
}
func DelNode(name string)(err error){
	if node,ok:=GetNode(name);ok{
		if err=node.Ws.WsConn.WriteMessage(websocket.TextMessage,[]byte(`你的连接已经断开了`));err!=nil{
			Nodes.Delete(name)
			return
		}
		err=errors.New(`你的连接已经断开了`)
		node.Ws.Close()
		Nodes.Delete(name)
	}
	return
}
func GetAllNode()([]Node,int)  {
	var count int
	var _Node []Node
	Nodes.Range(func(name,v interface{}) bool {
		count++
		node:=v.(*Node)
		_Node= append(_Node,*node)
		return true
	})
	return _Node,count
}
