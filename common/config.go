package common

import (
	"log"
	"os"
	"ws/util"

	"github.com/BurntSushi/toml"
)

type (
	mysqlCfg struct {
		ServerHost string `toml:"serverHost"`
		Port       uint16
		User       string
		Password   string
		Db         string
		MaxConnect int `toml:"maxConnect"`
	}

	dbCfg struct {
		Defalut string
		Mysql   mysqlCfg
	}

	httpCfg struct {
		HttpPort    uint16 `toml:"httpPort"`
		HttpTimeOut int    `toml:"httpTimeOut"`
	}

	wsCfg struct {
		WsPort    uint16 `toml:"wsPort"`
		Pong      bool
		WsTimeOut int `toml:"wsTimeOut"`
		ReadChan  int `toml:"readChan"`
		WriteChan int `toml:"writeChan"`
	}

	commonCfg struct {
		SignKey       string `toml:"signKey"`
		Name          string
		MaxBody       int `toml:"maxBody"`
		PidMod        os.FileMode `toml:"pidMod"`
		LogMod        os.FileMode `toml:"logMod"`
		MultiplexPort bool `toml:"multiplexPort"`
		Env           string
		MessageType   int `toml:"messageType"`
		Http          httpCfg
		WebSocket     wsCfg
	}

	validateDetail struct {
		Mold  int
		Name  string
		Query string
	}

	validateMethod struct {
		Mold      int
		Name      string
		Query     string
		Http      validateDetail
		WebSocket validateDetail
	}

	baseConfig struct {
		Common         commonCfg
		DB             dbCfg
		ValidateMethod validateMethod `toml:"validateMethod"`
	}
)

var (
	Conf  commonCfg
	DB    dbCfg
	Debug bool
	VM    validateMethod
)

func init() {
	var cfg baseConfig
	path := util.PathToEveryOne("config/config.toml")
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		log.Fatal("config/config.toml: ", err)
	}
	Conf = cfg.Common
	DB = cfg.DB
	VM = cfg.ValidateMethod
	VM.load()
	Debug = Conf.Env == "dev"
}

func (vm *validateMethod) load() {
	vm.WebSocket.merge(*vm)
	vm.Http.merge(*vm)
}

func (vd *validateDetail) merge(vm validateMethod) {
	if vd.Name == "" {
		vd.Name = vm.Name
	}
	if vd.Query == "" {
		vd.Query = vm.Query
	}
	if vd.Mold == 0 {
		vd.Mold = vm.Mold
	}
}

func LogInfo(s string) {
	log.Println("[info]", s)
}

func LogInfoSuccess(s string) {
	log.Println("[success]", s)
}

func LogInfoFailed(s string) {
	log.Println("[failed]", s)
}

func LogDebug(s string) {
	if Debug {
		log.Println("[debug]", s)
	}
}
