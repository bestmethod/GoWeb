package goweb

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
)

func (conf *Config) FromFile(ConfigFile string) {
	// get all configs apart from Loggers
	if conf.EarlyDebug == true {
		fmt.Println(EARLY_LOAD_MAINCONF)
	}
	if _, err := toml.DecodeFile(ConfigFile, conf); err != nil {
		log.Fatal(err)
	}
	// get logger configs
	if conf.EarlyDebug == true {
		fmt.Println(EARLY_LOAD_LOGCONF)
	}
	var l LoggersConf
	if _, err := toml.DecodeFile(ConfigFile, &l); err != nil {
		log.Fatal(err)
	}
	conf.Loggers = &l
}

func Load(fn string, early bool) (conf *Config) {
	if early == true {
		fmt.Println(EARLY_INIT_CONFSTRUCT)
	}
	conf = new(Config)
	conf.General = new(generic)
	conf.Database = new(DbConf)
	conf.Listener = new(RpcConf)
	conf.EarlyDebug = early
	if conf.EarlyDebug == true {
		fmt.Println(EARLY_LOADING_FILECONF)
	}
	conf.FromFile(fn)
	return conf
}
