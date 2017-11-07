package configurator

import "time"

type WebConf struct {
	Name *string
}

type Config struct {
	General    *generic
	Website    *WebConf
	Database   *DbConf
	Listener   *RpcConf
	Loggers    *LoggersConf
	EarlyDebug bool
}

type generic struct {
	MonitorSleepSeconds time.Duration
	DebugMainLoop       bool
	ServiceName         string
}

type DbConf struct {
	Type              string
	Server            string
	User              string
	Password          string
	DbName            string
	UseTLS            bool
	TLSSkipVerify     bool
	Ssl_ca            string
	ListenerConfID    int
	LoadLoggersFromDB bool
}

type RpcConf struct {
	Id                            int
	ListenIp                      string
	ListenPort                    int
	UseSSL                        bool //will only support TLS1.2 by default, not adding support for others.
	SSLCrtPath                    string
	SSLKeyPath                    string
	CookieLifetimeSeconds         int
	SessionExpireSeconds          int
	SessionCleanerRun             bool
	SessionCleanerIntervalSeconds int
	RpcListenerRun                bool
	SessionDebug                  bool
}

type LoggersConf struct {
	Logger []*loggerConf
}

type loggerConf struct {
	Id           int
	LogLevel     string // DEBUG, INFO, etc...
	Destination  string
	RpcLogLevel  string // rpcListener has it's own log level definitions so we can switch on debug in rpc, but not in the main code
	SessionDebug bool
	// udp://1.2.3.4:389, tcp://logger.example.com:1234, devlog, stdout, stderr
	// as such, we support syslog via tcp and udp, direct-to-file logging, /dev/log, stdout, stderr
	// and array of these can be specified! This basically means we can configure more than one destination and loglevel
}
