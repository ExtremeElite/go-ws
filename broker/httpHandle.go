package broker

import (
	"net/http"
	"ws/common"
)

var HttpChan chan PushData

func HttpHandle(w http.ResponseWriter, r *http.Request) {
	err := httpBroker(w, r)
	if err != nil {
		common.LogDebug(err.Error())
	}
}
