package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"ws/auth"
	"ws/common"
	"ws/kernel"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		for _, allowed := range common.Conf.CorsOrigins {
			if origin == allowed {
				return true
			}
		}
		return len(common.Conf.CorsOrigins) == 0
	},
}

const helloHTML = `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>ws-rs Go</title><style>body{font-family:system-ui,sans-serif;display:flex;justify-content:center;align-items:center;height:100vh;margin:0;background:#f5f5f5}.card{background:#fff;padding:2rem 3rem;border-radius:12px;box-shadow:0 2px 8px rgba(0,0,0,.1);text-align:center}h1{color:#333;margin:0}.p{color:#666}</style></head><body><div class="card"><h1>ws-rs</h1><p>WebSocket Server</p></div></body></html>`

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("[failed] writeJSON: %v", err)
	}
}

func HandlePush(w http.ResponseWriter, r *http.Request) {
	var d PushData
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"code": 400, "msg": "invalid json"})
		return
	}

	switch d.EventType {
	case Conversation:
		select {
		case HttpChan <- d:
			writeJSON(w, map[string]interface{}{"code": 200, "msg": "ok", "data": ""})
		default:
			w.WriteHeader(http.StatusTooManyRequests)
			writeJSON(w, map[string]interface{}{"code": 429, "msg": "too many requests"})
		}
	case GetOnlineInfo:
		writeJSON(w, map[string]interface{}{"code": 200, "msg": "ok", "data": onlineNames()})
	default:
		writeJSON(w, map[string]interface{}{"code": 200, "msg": "ok", "data": d})
	}
	log.Printf("[success] http push from %s", r.RemoteAddr)
}

func HandleToken(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	connectTypeStr := r.URL.Query().Get("connectType")

	id := 1
	connectType := 1

	if idStr != "" {
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			id = 1
		}
	}
	if connectTypeStr != "" {
		if _, err := fmt.Sscanf(connectTypeStr, "%d", &connectType); err != nil {
			connectType = 1
		}
	}

	writeJSON(w, map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": auth.CreateToken(id, connectType),
	})
}

func HandleRevokeToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		writeJSON(w, map[string]interface{}{"code": 400, "msg": "token required"})
		return
	}
	auth.RevokeToken(token)
	writeJSON(w, map[string]interface{}{"code": 200, "msg": "token revoked"})
}

func HandlePanicStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]interface{}{
		"kernel_panic_count": kernel.GetPanicCount(),
		"forward_panic_count": GetForwardPanicCount(),
	})
}

func HandleAllNodes(w http.ResponseWriter, r *http.Request) {
	nodes, total := kernel.GetAllNodes()
	writeJSON(w, map[string]interface{}{
		"info":  nodes,
		"total": total,
	})
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	if !(common.Conf.MultiplexPort || r.Header.Get("Connection") == "Upgrade") {
		w.Write([]byte(helloHTML))
		return
	}

	conn, name, err := wsUpgrade(w, r)
	if err != nil {
		log.Printf("[failed] ws upgrade: %v", err)
		return
	}

	mold := common.VM.WebSocket.Mold
	if !auth.ValidateToken(name, mold) {
		conn.WriteMsg([]byte(`{"code":401,"msg":"unauthorized"}`))
		conn.Close()
		return
	}

	if err := conn.WriteMsg([]byte(`{"code":200,"msg":"login success"}`)); err != nil {
		log.Printf("[failed] send login success: %v", err)
		conn.Close()
		return
	}
	if common.Conf.WebSocket.Pong {
		go conn.Pong()
	}
	defer conn.Close()

	kernel.AddNode(&kernel.Node{Ws: conn, Name: name, RemoteAddr: r.RemoteAddr})
	log.Printf("[success] ws open: %s name=%s", r.RemoteAddr, name)

	for {
		if err := wsHandle(conn); err != nil && !strings.Contains(err.Error(), "wsMessageForwarding") {
			break
		}
	}
	log.Printf("[success] ws close: %s", name)
}

func wsUpgrade(w http.ResponseWriter, r *http.Request) (*kernel.Connection, string, error) {
	token := r.URL.Query().Get("token")
	if token == "" {
		token = r.URL.Query().Get("sn")
	}
	if token == "" {
		return nil, "", fmt.Errorf("token required")
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, "", err
	}

	conn, err := kernel.BuildConn(wsConn)
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		wsConn.Close()
		return nil, "", err
	}

	return conn, token, nil
}

func wsHandle(conn *kernel.Connection) error {
	timeout := common.Conf.WebSocket.WsTimeOut
	if timeout > 0 {
		conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}

	msg, err := conn.ReadMsg()
	if err != nil {
		return err
	}

	s := strings.ToLower(strings.TrimSpace(string(msg)))
	if s == "ping" {
		return conn.WriteMsg([]byte("Pong"))
	}
	if s == "pong" {
		return conn.WriteMsg([]byte("Ping"))
	}

	log.Printf("[success] ws message: %s", msg)

	var d PushData
	if err := json.Unmarshal(msg, &d); err != nil {
		return conn.WriteMsg([]byte(`{"code":400,"msg":"invalid json"}`))
	}
	switch d.EventType {
	case Conversation:
		d.forward()
	}
	return nil
}
