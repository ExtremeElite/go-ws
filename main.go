package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"ws/broker"
	"ws/common"
	"ws/router"

	"github.com/mbndr/figlet4go"
	"github.com/sevlyar/go-daemon"
)

var runName = common.Common.Name

func init() {
	broker.HttpChan = make(chan broker.PushData, 1)
	Logo()
}
func Logo() {
	ascii := figlet4go.NewAsciiRender()
	// Adding the colors to RenderOptions
	options := figlet4go.NewRenderOptions()
	renderStr, _ := ascii.RenderOpts(strings.ToUpper(fmt.Sprintf("%v", runName)), options)
	fmt.Println(renderStr)
}
func main() {
	//后台进程守护
	if runtime.GOOS == "linux" {
		ctxt := &daemon.Context{
			PidFileName: fmt.Sprintf("%v.pid", runName),
			PidFilePerm: common.Common.PidMod,
			LogFileName: fmt.Sprintf("%v.log", runName),
			LogFilePerm: common.Common.LogMod,
			WorkDir:     "./",
			Umask:       022,
			Args:        []string{fmt.Sprintf("[go-daemon %v]", runName)},
		}
		d, err := ctxt.Search()
		if err == nil && d.Pid > 0 {
			log.Fatalf("%v is running,pid is %v", runName, d.Pid)
		}
		children, err := ctxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if children != nil {
			return
		}
		log.Print("- - - - - - - - - - - - - - -")
		log.Printf("%v started", runName)
		defer func(cntxt *daemon.Context) {
			_ = ctxt.Release()
		}(ctxt)
	}
	if !common.Common.MultiplexPort {
		go router.HttpPush()
	}
	go broker.HttpMessageForwarding()
	router.WsPush()

}
