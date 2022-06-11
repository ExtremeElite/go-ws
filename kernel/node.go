package kernel

import (
	"errors"
	"sync"
	"ws/common"

	"github.com/gorilla/websocket"
)

type Node struct {
	Ws         *Connection `json:"-"`
	Name       string      `json:"name"`
	RemoteAddr string      `json:"remote_addr"`
}

var Nodes sync.Map

func AddNode(node *Node) {
	if _, ok := GetNode(node.Name); ok {
		_ = DelNode(node.Name)
	}
	Nodes.Store(node.Name, node)
}
func GetNode(name string) (v *Node, ok bool) {
	var node *Node
	if v, ok := Nodes.Load(name); ok {
		return v.(*Node), ok
	}
	return node, ok
}
func DelNode(name string) (err error) {
	if node, ok := GetNode(name); ok {
		if err = node.Ws.WsConn.WriteMessage(websocket.TextMessage, []byte(common.ConnectClosed)); err != nil {
			Nodes.Delete(name)
			return
		}
		err = errors.New(common.ConnectClosed)
		node.Ws.Close()
		Nodes.Delete(name)
	}
	return
}
func GetAllNode() ([]Node, int) {
	var count int
	var _Node []Node
	Nodes.Range(func(name, v interface{}) bool {
		count++
		node := v.(*Node)
		_Node = append(_Node, *node)
		return true
	})
	return _Node, count
}
