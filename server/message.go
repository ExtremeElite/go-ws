package server

import (
	"encoding/json"
	"log"
	"ws/kernel"
)

const (
	Conversation  = 1
	Login         = 2
	Logout        = 3
	GetOnlineInfo = 4
)

type PushData struct {
	EventType      int32       `json:"event_type"`
	PublishAccount []string    `json:"publish_account"`
	Data           interface{} `json:"data"`
}

func (d PushData) String() string {
	switch v := d.Data.(type) {
	case string:
		return v
	default:
		b, _ := json.Marshal(d.Data)
		return string(b)
	}
}

func (d PushData) forward() {
	for _, account := range d.PublishAccount {
		if node, ok := kernel.GetNode(account); ok {
			go func(c *kernel.Connection, msg string) {
				_ = c.WriteMsg([]byte(msg))
			}(node.Ws, d.String())
		}
	}
}

var HttpChan = make(chan PushData, 10)

func HttpMessageForwarding() {
	for d := range HttpChan {
		d.forward()
		log.Println("[success] http push:", d.String())
	}
}

func onlineNames() []string {
	nodes, _ := kernel.GetAllNode()
	names := make([]string, len(nodes))
	for i, n := range nodes {
		names[i] = n.Name
	}
	return names
}
