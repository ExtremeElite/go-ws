package common

import (
	"github.com/BurntSushi/toml"
	"log"
	"ws/util"
)

type Mysql struct {
	ServerHost string
	Port       uint16
	User       string
	Password   string
	Db         string
	MaxConnect int `toml:"maxConnect"`
}
type Common struct {
	WsPort      uint16 `validate:"required,max=65535,min=0" label:"websocket端口"`
	HttpPort    uint16 `validate:"required,max=65535,min=0,nefield=WsPort" label:"Http端口"`
	Env         string `validate:"required" label:"环境变量"`
	SignKey     string
	DefaultDB   string `validate:"required"`
	WsTimeOut   int    `validate:"required,min=5,max=300" label:"websocket连接超时"`
	ReadChan    int    `validate:"required,min=2,max=10000" label:"读协程"`
	WriteChan   int    `validate:"required,min=2,max=10000" label:"写协程"`
	MaxBody     int    `validate:"required,min=5,max=100000" label:"请求体"`
	HttpTimeOut int    `validate:"required,min=5,max=30" label:"http请求超时时间"`
}

type BaseServer struct {
	Common  Common
	MysqlDB Mysql
}

var (
	Setting  Common
	MysqlSet Mysql
	Debug    bool
)

func init() {
	var bs = Config()
	Setting = bs.Common
	MysqlSet = bs.MysqlDB
	Debug = bs.Common.Env == "dev"
}
func Config() BaseServer {
	var bs BaseServer
	var configPath string
	configPath = util.PathToEveryOne(`config/config.toml`)
	_, err := toml.DecodeFile(configPath, &bs)
	if err != nil {
		log.Fatal("please check config/config.toml", err.Error())
	}
	util.ValidateStruct(bs.Common)
	return bs
}
func CheckPort(port int) error {

	return nil
}
func LogDebug(s string) {
	if Debug {
		util.LogUtil(s, "debug")
	}
}
func LogInfo(s string) {
	util.LogUtil(s, "info")
}
