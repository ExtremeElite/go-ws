package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"ws/conf"
)
var upgrader=websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func main(){
	var wsPort= conf.GetConfig().WsPort
	http.HandleFunc("/", echo)
	http.ListenAndServe(":"+strconv.Itoa(int(wsPort)), nil)
}
func echo(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}