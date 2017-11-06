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
		fmt.Println("LOADLOGGER: DEBUG Creating log handler")
	}
	l = new(LogHandler)
	l.Dispatchers = []destination{destination{}}
	l.ServiceName = ServiceName
	l.pid = os.Getpid()
	var Destination destination
	if earlyDebug == true {
		fmt.Println("LOADLOGGER: DEBUG Entering loop")
	}
	for i := 0; i < len(LogConf.Logger); i++ {
		if earlyDebug == true {
			fmt.Printf("LOADLOGGER: DEBUG Run %d Setting log level\n", i)
		}
		Destination.SysLog = nil
		Destination.Stderr = nil
		Destination.Stdout = nil
		Destination.File = ""
		switch LogConf.Logger[i].LogLevel {
		case "DEBUG":
			Destination.LogLevel = LEVEL_DEBUG
		case "INFO":
			Destination.LogLevel = LEVEL_INFO
		case "WARN":
			Destination.LogLevel = LEVEL_WARN
		case "ERROR":
			Destination.LogLevel = LEVEL_ERROR
		case "CRITICAL":
			Destination.LogLevel = LEVEL_CRITICAL
		}
		var s *syslog.Writer
		var err error
		if earlyDebug == true {
			fmt.Printf("LOADLOGGER: DEBUG Run %d Checking destination\n", i)
		}
		if LogConf.Logger[i].Destination[0:6] == "tcp://" || LogConf.Logger[i].Destination[0:6] == "udp://" {
			if earlyDebug == true {
				fmt.Printf("LOADLOGGER: DEBUG Run %d Dialling syslog: %s\n", i, LogConf.Logger[i].Destination)
			}
			s, err = syslog.Dial(LogConf.Logger[i].Destination[0:3], LogConf.Logger[i].Destination[6:], syslog.LOG_DAEMON, "Jarvis")
			if err != nil {
				log.Fatalf("ERROR: Could not initialize logger:%s\n", err)
			}
			Destination.SysLog = s
		} else if LogConf.Logger[i].Destination == "stdout" {
			if earlyDebug == true {
				fmt.Printf("LOADLOGGER: DEBUG Run %d Creating stdout logger\n", i)
			}
			Destination.Stdout = log.New(os.Stdout, "", 0)
		} else if LogConf.Logger[i].Destination == "stderr" {
			if earlyDebug == true {
				fmt.Printf("LOADLOGGER: DEBUG Run %d Creating stderr logger\n", i)
			}
			Destination.Stderr = log.New(os.Stderr, "", 0)
		} else if LogConf.Logger[i].Destination == "devlog" {
			if earlyDebug == true {
				fmt.Printf("LOADLOGGER: DEBUG Run %d Dialling syslog:devlog\n", i)
			}
			s, err = syslog.Dial("", "", syslog.LOG_DAEMON, "Jarvis")
			if err != nil {
				log.Fatalf("ERROR: Could not initialize logger:%s\n", err)
			}
			Destination.SysLog = s
		}
		if earlyDebug == true {
			fmt.Printf("LOADLOGGER: DEBUG Run %d Appending array\n", i)
		}
		l.Dispatchers = append(l.Dispatchers, Destination)
	}
	l.loadMessages()
	return l
}
