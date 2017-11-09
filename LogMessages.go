package goweb

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
	LOG_RAND_FAIL      = "Could not generate random number for a cookie, this is serious!\n"

	//SessionCleaner
	LOG_CLEANER_FART = "Could not cleanup session data in the DB: %s\n"
	LOG_CLEANER_DONE = "Session Cleaner Complete"
)

// db
const (
	//logging lines for DB - db
	LOG_SQLITE3_OPEN       = "sqlite3 connection, opening"
	LOG_MYSQL_OPEN         = "MySQL connection, opening"
	LOG_MYSQL_TLS_TRUE     = "MySQL tls true"
	LOG_MYSQL_TLS_SKIP     = "MySQL tls skip-verify"
	LOG_MYSQL_VOODOO_START = "MySQL ssl_ca voodoo starting"
	LOG_MYSQL_PEMFAIL      = "Failed to append PEM."
	LOG_MYSQL_VOODOO_DONE  = "MySQL ssl_ca voodoo done, tls custom"
	LOG_MYSQL_CONN         = "MySQL opening connection"
	LOG_WRONG_DBTYPE       = "ERROR: The only supported database types are: MySQL | sqlite3. Currently set: %s\n"
	LOG_MEDDLER_DEFAULTS   = "Configuring meddler defaults"
)
