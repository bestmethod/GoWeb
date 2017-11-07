package rpcDb

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
