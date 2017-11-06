package configurator

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
)

func (conf *Config) FromFile(ConfigFile string) {
	// get all configs apart from Loggers
	if conf.EarlyDebug == true {
		fmt.Println("LOADCONFIG:FROMFILE: DEBUG Loading main configuration parts")
	}
	if _, err := toml.DecodeFile(ConfigFile, conf); err != nil {
		log.Fatal(err)
	}
	// get logger configs
	if conf.EarlyDebug == true {
		fmt.Println("LOADCONFIG:FROMFILE: DEBUG Loading loggers")
	}
	var l LoggersConf
	if _, err := toml.DecodeFile(ConfigFile, &l); err != nil {
		log.Fatal(err)
	}
	conf.Loggers = &l
}

func Load(fn string, early bool) (conf *Config) {
	if early == true {
		fmt.Println("LOADCONFIG: DEBUG Initializing config structure")
	}
	conf = new(Config)
	conf.General = new(generic)
	conf.Database = new(DbConf)
	conf.Listener = new(RpcConf)
	conf.EarlyDebug = early
	if conf.EarlyDebug == true {
		fmt.Println("LOADCONFIG: DEBUG Loading config from file")
	}
	conf.FromFile(fn)
	if conf.EarlyDebug == true {
		fmt.Println("LOADCONFIG: DEBUG Checking if we should load config from DB")
	}
	return conf
}
