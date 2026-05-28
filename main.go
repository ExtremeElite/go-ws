package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"ws/common"
	"ws/server"

	"github.com/mbndr/figlet4go"
	"github.com/sevlyar/go-daemon"
)

func init() {
	server.InitHttpChan(common.Conf.WebSocket.WriteChan)
	logo()
}

func logo() {
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	s, _ := ascii.RenderOpts(strings.ToUpper(fmt.Sprintf("%v", common.Conf.Name)), options)
	fmt.Println(s)
}

func main() {
	if runtime.GOOS == "linux" {
		ctxt := &daemon.Context{
			PidFileName: fmt.Sprintf("%v.pid", common.Conf.Name),
			PidFilePerm: common.Conf.PidMod,
			LogFileName: fmt.Sprintf("%v.log", common.Conf.Name),
			LogFilePerm: common.Conf.LogMod,
			WorkDir:     ".",
			Umask:       022,
			Args:        []string{fmt.Sprintf("[go-daemon %v]", common.Conf.Name)},
		}
		d, err := ctxt.Search()
		if err == nil && d.Pid > 0 {
			log.Fatalf("%v is running, pid %v", common.Conf.Name, d.Pid)
		}
		children, err := ctxt.Reborn()
		if err != nil {
			log.Fatal("unable to run: ", err)
		}
		if children != nil {
			return
		}
		log.Printf("%v started", common.Conf.Name)
		defer ctxt.Release()
	}

	if !common.Conf.MultiplexPort {
		go server.HttpPush()
	}
	go server.HttpMessageForwarding(4)
	server.WsPush()
}
