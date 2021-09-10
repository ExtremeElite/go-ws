package main

import (
	"fmt"
	"github.com/mbndr/figlet4go"
	"ws/broker"
	"ws/common"
	"ws/router"
)

func init() {
	broker.HttpChan = make(chan broker.PushData, 1)
	Logo()
}
func Logo() {
	ascii := figlet4go.NewAsciiRender()
	// Adding the colors to RenderOptions
	options := figlet4go.NewRenderOptions()
	renderStr, _ := ascii.RenderOpts("Gorouting", options)
	fmt.Println(renderStr)
}
func main() {
	//后台进程守护
	common.DaemonRun()
	go router.HttpPush()
	go broker.HttpMessageForwarding()
	router.WsPush()

}
