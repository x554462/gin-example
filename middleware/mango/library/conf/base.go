package conf

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

type App struct {
}

type Server struct {
	RunMode      string
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	HttpTimeout  time.Duration
	RuntimePath  string
	LogPath      string
	LogName      string
}

type Database struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}
type Redis struct {
	Host     string
	Port     int
	Password string
}

var (
	ServerConf         = &Server{}
	MasterDatabaseConf = &Database{}
	SlaveDatabaseConf  = &Database{}
	RedisConf          = &Redis{}
)

var cfg *ini.File

func init() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("conf.init: fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("Server", ServerConf)
	mapTo("MasterDatabase", MasterDatabaseConf)
	mapTo("SlaveDatabase", SlaveDatabaseConf)
	mapTo("Redis", RedisConf)

	ServerConf.ReadTimeout = ServerConf.ReadTimeout * time.Second
	ServerConf.WriteTimeout = ServerConf.WriteTimeout * time.Second
	ServerConf.HttpTimeout = ServerConf.HttpTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
