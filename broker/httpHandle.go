package broker

import (
	"log"
	"net/http"
)

var HttpChan chan PushData

func HttpHandle(w http.ResponseWriter, r *http.Request) {
	err := httpBroker(w, r)
	if err != nil {
		log.Println("broker httpHandle line 13 error:", err.Error())
		return
	}
}
