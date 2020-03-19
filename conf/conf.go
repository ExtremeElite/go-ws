package conf

import (
	"github.com/BurntSushi/toml"
	"log"
	"ws/util"
)
type Mysql struct {
	ServerHost string
	Port uint16
	User string
	Password string
	Db string
	MaxConnect int `toml:"maxConnect"`
}
type Common struct {
	WsPort uint16
	HttpPort uint16
	Env string
	SignKey string
	DefaultDB string
}

type BaseServer struct {
	Common Common
	MysqlDB Mysql
}
func GetConfig() BaseServer  {
	var bs BaseServer
	var configPath string
	configPath=util.PathToEveryOne(`config/config.toml`)
	_, err := toml.DecodeFile(configPath, &bs)
	if err != nil {
		log.Fatal("please check config/config.toml")
	}
	return bs
}
