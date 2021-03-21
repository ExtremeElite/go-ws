package conf

var (
	CommonSet Common
	MysqlSet  Mysql
)

func init() {
	var bs = Config()
	CommonSet = bs.Common
	MysqlSet = bs.MysqlDB
}
