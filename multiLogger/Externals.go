package multiLogger

import (
	"fmt"
	"time"
)

func (l *LogHandler) Debug(m string) {
	l.dispatch(LEVEL_DEBUG, m)
}

func (l *LogHandler) Info(m string) {
	l.dispatch(LEVEL_INFO, m)
}

func (l *LogHandler) Warn(m string) {
	l.dispatch(LEVEL_WARN, m)
}

func (l *LogHandler) Error(m string) {
	l.dispatch(LEVEL_ERROR, m)
}

func (l *LogHandler) Critical(m string) {
	l.dispatch(LEVEL_CRITICAL, m)
}

func (l *LogHandler) dispatch(logLevel int, m string) {
	var mm string
	if logLevel == LEVEL_DEBUG {
		mm = fmt.Sprintf("DEBUG    %s %s", l.Header, m)
	} else if logLevel == LEVEL_INFO {
		mm = fmt.Sprintf("INFO     %s %s", l.Header, m)
	} else if logLevel == LEVEL_WARN {
		mm = fmt.Sprintf("WARN     %s %s", l.Header, m)
	} else if logLevel == LEVEL_ERROR {
		mm = fmt.Sprintf("ERROR    %s %s", l.Header, m)
	} else if logLevel == LEVEL_CRITICAL {
		mm = fmt.Sprintf("CRITICAL %s %s", l.Header, m)
	}
	mm = fmt.Sprintf("%s %s[%d]: %s", time.Now().UTC().Format("Jan 02 15:04:05-0700"), l.ServiceName, l.pid, mm)
	for i := 0; i < len(l.Dispatchers); i++ {
		if l.Dispatchers[i].LogLevel >= logLevel {
			if l.Dispatchers[i].SysLog != nil {
				if logLevel == LEVEL_DEBUG {
					l.Dispatchers[i].SysLog.Debug(mm)
				} else if logLevel == LEVEL_INFO {
					l.Dispatchers[i].SysLog.Info(mm)
				} else if logLevel == LEVEL_WARN {
					l.Dispatchers[i].SysLog.Warning(mm)
				} else if logLevel == LEVEL_ERROR {
					l.Dispatchers[i].SysLog.Err(mm)
				} else if logLevel == LEVEL_CRITICAL {
					l.Dispatchers[i].SysLog.Crit(mm)
				}
			} else if l.Dispatchers[i].Stdout != nil {
				l.Dispatchers[i].Stdout.Println(mm)
			} else if l.Dispatchers[i].Stderr != nil {
				l.Dispatchers[i].Stderr.Println(mm)
			}
		}
	}
}
