package kernel

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
	"ws/common"
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
func GetNode(name string) (*Node, bool) {
	var node *Node
	if v, ok := Nodes.Load(name); ok {
		node = v.(*Node)
		return node, true
	}
	return node, false
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