package rpcListener

const (
	//Serve
	LOG_CONFIG_GLOBALS   = "Configuring rpcListener global vars"
	LOG_ROUTER_MAKE      = "Creating new router and setting session data"
	LOG_STARTLISTEN      = "Starting Listener on %s:%d, TLS=%t\n"
	LOG_CONF_DISPATCHERS = "Configuring router dispatchers"

	//Server :: session functions
	LOG_DB_QUERY_FAIL  = "Could not find cookie session data in the DB: %s\n"
	LOG_DB_INSERT_FAIL = "Could not insert cookie session data to the DB: %s\n"
	LOG_DB_UPDATE_FAIL = "Could not update cookie session data in the DB: %s\n"

	//SessionCleaner
	LOG_CLEANER_FART = "Could not cleanup session data in the DB: %s\n"
	LOG_CLEANER_DONE = "Session Cleaner Complete"
)
