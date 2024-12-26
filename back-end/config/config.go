package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
)

// MySQL配置信息
type MySQL struct {
	Host      string `toml:"host"` //MySQL服务器主机名或IP地址
	Port      int    `toml:"port"` //MySQL服务器端口号
	Database  string `toml:"database"` //要连接的数据库名称
	Username  string `toml:"username"` //登录MySQL服务器的用户名
	Password  string `toml:"password"` //登录MySQL服务器的密码
	Charset   string `toml:"charset"` //连接使用的字符集
	ParseTime bool   `toml:"parse_time"` //是否将MySQL返回的时间类型解析为Go的本地时间类型
	Local     string `toml:"local"`      //指定本地时区
}

// 服务器配置信息
type Server struct {
	IP   string //服务器IP地址
	Port int    //服务器端口号
}

// 路径配置信息
type Path struct {
	StaticSourcePath string `toml:"static_source_path"` //静态文件路径
}

type Config struct {
	DB     MySQL `toml:"mysql"`
	Server `toml:"server"`
}

var Conf Config

// 初始化函数
// 程序启动时会自动读取配置文件，并确保配置信息的有效性。这样可以避免在代码中硬编码配置信息，使得程序更具可维护性和可扩展性
func init() {
	//使用toml包的DecodeFile()函数解析配置文件config.toml
	if _, err := toml.DecodeFile("./config/config.toml", &Conf); err != nil {
		panic(err)
	}

	//去除左右的空格
	strings.Trim(Conf.Server.IP, " ")
	// strings.Trim(Conf.RDB.IP, " ") Redis 还没学
	strings.Trim(Conf.DB.Host, " ")
}

// 填充得到数据库连接字符串
func DBConnectString() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		Conf.DB.Username, Conf.DB.Password, Conf.DB.Host, Conf.DB.Port, Conf.DB.Database,
		Conf.DB.Charset, Conf.DB.ParseTime, Conf.DB.Local)
	log.Println(dsn)
	return dsn
}
