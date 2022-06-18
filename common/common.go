package common

import (
	"log"
	"os"
	"ws/util"

	"github.com/BurntSushi/toml"
)

type mysql struct {
	ServerHost string `validate:"required,ip" label:"数据服务器地址"`
	Port       uint16 `validate:"required,min=0,max=65535" label:"数据服务器地址" `
	User       string `validate:"required" label:"账户" `
	Password   string `validate:"required" label:"密码"`
	Db         string `validate:"required" label:"数据库名称"`
	MaxConnect int    `toml:"maxConnect" validate:"required,max=1000,min=5" label:"最大连接数"`
}
type db struct {
	Defalut string `validate:"required oneof=mysql" label:"默认数据库"`
	Mysql   mysql  `toml:"mysql"`
}
type http struct {
	HttpPort    uint16 `validate:"required,min=0,max=65535,nefield=WsPort" label:"Http端口"`
	HttpTimeOut int    `validate:"required,min=5,max=30" label:"http请求超时时间"`
}
type websocket struct {
	WsPort    uint16 `validate:"required,min=0,max=65535" label:"websocket端口"`
	Pong      bool   `validate:"-" label:"pong"`
	WsTimeOut int    `validate:"required,min=5,max=300" label:"websocket连接超时"`
	ReadChan  int    `validate:"required,min=2,max=10000" label:"读协程"`
	WriteChan int    `validate:"required,min=2,max=10000" label:"写协程"`
}
type common struct {
	SignKey        string      `validate:"" label:"环境变量"`
	Name           string      `validate:"required,min=0,max=32" label:"名称"`
	MaxBody        int         `validate:"required,min=5,max=100000" label:"请求体"`
	PidMod         os.FileMode `validate:"required,numeric,oneof=777 755" label:"pid文件权限"`
	LogMod         os.FileMode `validate:"required,numeric,oneof=777 755" label:"log文件权限"`
	MultiplexPort  bool        `validate:"-" label:"端口复用"`
	Env            string      `validate:"required,oneof=dev prod" label:"环境变量"`
	ValidateMethod validateMethod
	WebSocket      websocket
	Http           http
}
type validateDetail struct {
	Mold  int    `validate:"numeric,oneof=0 1 2" label:"启用什么类型的验证"`
	Name  string `validate:"required,min=0,max=32" label:"验证名称"`
	Query string `validate:"" label:"数据库验证查询"`
}
type validateMethod struct {
	ValidateHttp validateDetail
	ValidateWs   validateDetail
}

type BaseServer struct {
	Common         common
	DB             db
	ValidateMethod validateMethod
}

var (
	Common common
	DB     db
	Debug  bool
	Http   http
	Ws     websocket
	bs     BaseServer
)

func init() {
	var bs = Config()
	Http = bs.Common.Http
	Ws = bs.Common.WebSocket
	Common = bs.Common
	DB = bs.DB
	Debug = bs.Common.Env == "dev"
}
func Config() BaseServer {
	var configPath = util.PathToEveryOne(`config/config.toml`)
	_, err := toml.DecodeFile(configPath, &bs)
	if err != nil {
		log.Fatal("please check config/config.toml", err.Error())
	}
	LogDebug(bs)
	util.ValidateStruct(bs.Common.Http)
	util.ValidateStruct(bs.Common.WebSocket)
	util.ValidateStruct(bs.DB.Mysql)
	return bs
}
func CheckPort(port int) error {

	return nil
}
func LogDebug(s interface{}) {
	var Debug = bs.Common.Env == "dev"
	if Debug {
		util.LogUtil(s, "debug", Debug)
	}
}
func LogInfo(s string) {
	util.LogUtil(s, "info", false)
}
func LogInfoSuccess(s string) {
	LogInfo(util.MessageHeaderSuccess + s)
}
func LogInfoFailed(s string) {
	LogInfo(util.MessageHeaderFailed + s)
}
