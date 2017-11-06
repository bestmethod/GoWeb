package multiLogger

func (l *LogHandler) loadMessages() {
	l.Messages = append(l.Messages, EARLY_ARGS_FAIL)
	l.Messages = append(l.Messages, LOG_BASIC_LOAD_DONE)
	l.Messages = append(l.Messages, LOG_START_RPC)
	l.Messages = append(l.Messages, LOG_RESTART_RPC)
	l.Messages = append(l.Messages, LOG_INITIAL_DONE)
}

const (
	EARLY_ARGS_FAIL     = "Early Fail: Usage: %s [--early-debug] {config_file_name}\n" //process name
	LOG_BASIC_LOAD_DONE = "Config loaded, database connection tested, logger initialized ; soon there will be cake!"
	LOG_START_RPC       = "Starting RPC"
	LOG_RESTART_RPC     = "Restarting RPC"
	LOG_INITIAL_DONE    = "Now we have icing!"
)
