package common

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"runtime"
	"ws/util"
)

type Mysql struct {
	ServerHost string `validate:"required,ip" label:"数据服务器地址"`
	Port       uint16 `validate:"required,min=0,max=65535" label:"数据服务器地址"`
	User       string `validate:"required" label:"账户"`
	Password   string `validate:"required" label:"密码"`
	Db         string `validate:"required" label:"数据库名称"`
	MaxConnect int    `toml:"maxConnect" validate:"required,max=1000,min=5" label:"最大连接数"`
}
type Common struct {
	Name        string      `validate:"required,min=0,max=32" label:"名称"`
	PidMod      os.FileMode `validate:"required,numeric,oneof=777 755" label:"pid文件权限"`
	LogMod      os.FileMode `validate:"required,numeric,oneof=777 755" label:"log文件权限"`
	WsPort      uint16      `validate:"required,min=0,max=65535" label:"websocket端口"`
	HttpPort    uint16      `validate:"required,min=0,max=65535,nefield=WsPort" label:"Http端口"`
	Env         string      `validate:"required,oneof=dev prod" label:"环境变量"`
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
	util.ValidateStruct(bs.MysqlDB)
	return bs
}
func CheckPort(port int) error {

	return nil
}
func LogDebug(s string) {
	if Debug {
		util.LogUtil(s, "debug", true)
	}
}
func LogInfo(s string) {
	util.LogUtil(s, "info", false)
}

//后台进程守护
func DaemonRun() {
	if runtime.GOOS == "linux" {
		cntxt := &daemon.Context{
			PidFileName: fmt.Sprintf("%v.pid", Setting.Name),
			PidFilePerm: Setting.PidMod,
			LogFileName: fmt.Sprintf("%v.log", Setting.Name),
			LogFilePerm: Setting.LogMod,
			WorkDir:     "./",
			Umask:       022,
			Args:        []string{fmt.Sprintf("[go-daemon %v]", Setting.Name)},
		}
		d,err:=cntxt.Search()
		if d.Pid>0 {
			log.Fatal(fmt.Sprintf("%v is running,pid is %v",Setting.Name,d.Pid))
		}
		children, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if children != nil {
			return
		}
		log.Print("- - - - - - - - - - - - - - -")
		log.Print(fmt.Sprintf("%v started", Setting.Name))
		defer func(cntxt *daemon.Context) {
			_ = cntxt.Release()
		}(cntxt)
	}
}
