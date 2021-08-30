package main

import (
	"fmt"
	"github.com/mbndr/figlet4go"
	"github.com/sevlyar/go-daemon"
	"log"
	"runtime"
	"ws/broker"
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
	if runtime.GOOS == "linux" {
		cntxt := &daemon.Context{
			PidFileName: "go_ws.pid",
			PidFilePerm: 0777,
			LogFileName: "go_ws.log",
			LogFilePerm: 0777,
			WorkDir:     "./",
			Umask:       022,
			Args:        []string{"[go-daemon go_ws]"},
		}

		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if d != nil {
			return
		}
		log.Print("- - - - - - - - - - - - - - -")
		log.Print("go_ws started")
		defer cntxt.Release()
	}
	go router.HttpPush()
	go broker.HttpMessageForwarding()
	router.WsPush()

}
