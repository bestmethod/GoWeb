package multiLogger

import (
	"log"
	"log/syslog"
)

const (
	LEVEL_DEBUG    = 5
	LEVEL_INFO     = 4
	LEVEL_WARN     = 3
	LEVEL_ERROR    = 2
	LEVEL_CRITICAL = 1
)

type LogHandler struct {
	Dispatchers []destination
	Header      string
	ServiceName string
	pid         int
	Messages    []string
}

type destination struct {
	SysLog   *syslog.Writer
	Stdout   *log.Logger
	Stderr   *log.Logger
	File     string
	LogLevel int
}
