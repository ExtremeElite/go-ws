package server

import (
	"encoding/json"
	"log"
	"sync/atomic"
	"ws/kernel"
)

const (
	Conversation  = 1
	Login         = 2
	Logout        = 3
	GetOnlineInfo = 4
)

var forwardPanicCount int64

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
	msg := d.String()
	for _, account := range d.PublishAccount {
		node, ok := kernel.GetNode(account)
		if !ok {
			log.Printf("[warn] forward: node %s not found", account)
			continue
		}
		go func(c *kernel.Connection, target string) {
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt64(&forwardPanicCount, 1)
					log.Printf("[PANIC] forward to %s: %v", target, r)
				}
			}()
			if err := c.WriteMsg([]byte(msg)); err != nil {
				log.Printf("[failed] forward to %s: %v", target, err)
			}
		}(node.Ws, account)
	}
}

var HttpChan = make(chan PushData, 100)

func InitHttpChan(capacity int) {
	if capacity > 0 {
		HttpChan = make(chan PushData, capacity)
	}
}

func HttpMessageForwarding(workers int) {
	if workers <= 0 {
		workers = 4
	}
	for i := 0; i < workers; i++ {
		go func() {
			for d := range HttpChan {
				d.forward()
				log.Println("[success] http push:", d.String())
			}
		}()
	}
}

func GetForwardPanicCount() int64 {
	return atomic.LoadInt64(&forwardPanicCount)
}

func onlineNames() []string {
	nodes, _ := kernel.GetAllNodes()
	names := make([]string, len(nodes))
	for i, n := range nodes {
		names[i] = n.Name
	}
	return names
}
