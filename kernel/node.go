package kernel

import (
	"log"
	"sync"
)

type Node struct {
	Ws         *Connection `json:"-"`
	Name       string      `json:"name"`
	RemoteAddr string      `json:"remote_addr"`
}

var (
	Nodes   sync.Map
	nodeMu  sync.Mutex
)

func AddNode(node *Node) {
	nodeMu.Lock()
	defer nodeMu.Unlock()

	if old, ok := Nodes.Load(node.Name); ok {
		oldNode := old.(*Node)
		log.Printf("[info] replacing node %s", node.Name)
		oldNode.Ws.Close()
		Nodes.Delete(node.Name)
	}
	Nodes.Store(node.Name, node)
}

func GetNode(name string) (v *Node, ok bool) {
	if v, ok := Nodes.Load(name); ok {
		return v.(*Node), ok
	}
	return nil, ok
}

func DelNode(name string) {
	nodeMu.Lock()
	defer nodeMu.Unlock()

	if node, ok := Nodes.Load(name); ok {
		n := node.(*Node)
		n.Ws.Close()
		Nodes.Delete(name)
	}
}

func GetAllNodes() ([]Node, int) {
	var nodes []Node
	var count int

	Nodes.Range(func(name, v interface{}) bool {
		count++
		node := v.(*Node)
		nodes = append(nodes, Node{
			Name:       node.Name,
			RemoteAddr: node.RemoteAddr,
		})
		return true
	})
	return nodes, count
}

func SendToNode(name string, data []byte) bool {
	if node, ok := GetNode(name); ok {
		if err := node.Ws.WriteMsg(data); err != nil {
			log.Printf("[failed] send to %s: %v", name, err)
			return false
		}
		return true
	}
	return false
}
