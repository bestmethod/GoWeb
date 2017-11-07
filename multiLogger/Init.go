package multiLogger

import (
	"../configurator"
	"fmt"
	"log"
	"log/syslog"
	"os"
)

func Init(LogConf *configurator.LoggersConf, earlyDebug bool, ServiceName string) (l *LogHandler) {
	if earlyDebug == true {
		fmt.Println(EARLY_CREATING_HANDLER)
	}
	l = new(LogHandler)
	l.Dispatchers = []destination{destination{}}
	l.ServiceName = ServiceName
	l.pid = os.Getpid()
	var Destination destination
	if earlyDebug == true {
		fmt.Println(EARLY_ENTER_CONFLOOP)
	}
	for i := 0; i < len(LogConf.Logger); i++ {
		if earlyDebug == true {
			fmt.Printf(EARLY_LOOP_POS, i)
		}
		Destination.SysLog = nil
		Destination.Stderr = nil
		Destination.Stdout = nil
		Destination.File = ""
		Destination.LogLevel = logStringToInt(LogConf.Logger[i].LogLevel)
		Destination.RpcLogLevel = logStringToInt(LogConf.Logger[i].RpcLogLevel)
		var s *syslog.Writer
		var err error
		if earlyDebug == true {
			fmt.Printf(EARLY_CHECKING_DEST, i)
		}
		if LogConf.Logger[i].Destination[0:6] == "tcp://" || LogConf.Logger[i].Destination[0:6] == "udp://" {
			if earlyDebug == true {
				fmt.Printf(EARLY_DIALLING_SYSLOG, i, LogConf.Logger[i].Destination)
			}
			s, err = syslog.Dial(LogConf.Logger[i].Destination[0:3], LogConf.Logger[i].Destination[6:], syslog.LOG_DAEMON, l.ServiceName)
			if err != nil {
				log.Fatalf(EARLY_LOGGER_FAIL, err)
			}
			Destination.SysLog = s
		} else if LogConf.Logger[i].Destination == "stdout" {
			if earlyDebug == true {
				fmt.Printf(EARLY_CREATE_STDOUT, i)
			}
			Destination.Stdout = log.New(os.Stdout, "", 0)
		} else if LogConf.Logger[i].Destination == "stderr" {
			if earlyDebug == true {
				fmt.Printf(EARLY_CREATE_STDERR, i)
			}
			Destination.Stderr = log.New(os.Stderr, "", 0)
		} else if LogConf.Logger[i].Destination == "devlog" {
			if earlyDebug == true {
				fmt.Printf(EARLY_SET_DEVLOG, i)
			}
			s, err = syslog.Dial("", "", syslog.LOG_DAEMON, l.ServiceName)
			if err != nil {
				log.Fatalf(EARLY_LOGGER_FAIL, err)
			}
			Destination.SysLog = s
		}
		if earlyDebug == true {
			fmt.Printf(EARLY_APPEND_ARR, i)
		}
		l.Dispatchers = append(l.Dispatchers, Destination)
	}
	return l
}

func logStringToInt(logStr string) int {
	switch logStr {
	case "DEBUG":
		return LEVEL_DEBUG
	case "INFO":
		return LEVEL_INFO
	case "WARN":
		return LEVEL_WARN
	case "ERROR":
		return LEVEL_ERROR
	case "CRITICAL":
		return LEVEL_CRITICAL
	}
	return LEVEL_DEBUG
}
