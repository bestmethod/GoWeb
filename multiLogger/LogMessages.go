package multiLogger

// main
const (
	//log header definitions
	LOGHEADER_MAIN    = "MAIN"
	LOGHEADER_RPC     = "RPC-GOROUTINE"
	LOGHEADER_CLEANER = "CLEANER-GOROUTINE"
	LOGHEADER_RPCWEB  = "RPC-WEBSERVER"
	LOGHEADER_CLNSESS = "RPC-SESSIONCLEANER"
	LOGHEADER_DB      = "DBCONNECT"

	//early logging lines - before logger is setup
	EARLY_ARGS_FAIL      = "Early Fail: Usage: %s [--early-debug] {config_file_name}\n" //process name
	EARLY_LOADING_ARGS   = "MAIN: Processing command line parameters"
	EARLY_DEBUG_ON       = "MAIN: DEBUG Early Debug enabled\nMAIN: DEBUG Conf file: %s\n" //config file name
	EARLY_LOADING_CONF   = "MAIN: Loading configuration"
	EARLY_LOADING_LOGGER = "MAIN: Loading logger"

	//logging lines
	LOG_CONF_CHANNELS         = "Configuring channels for goroutines"
	LOG_BASIC_LOAD_DONE       = "Config loaded, database connection tested, logger initialized ; soon there will be cake!"
	LOG_START_RPC             = "Starting RPC"
	LOG_RESTART_RPC           = "Restarting RPC"
	LOG_INITIAL_DONE          = "Now we have icing!"
	LOG_GETTING_DB_CONN       = "Getting Database Connection for Rpc"
	LOG_ENTER_MAIN_LOOP       = "Entering main loop"
	LOG_MAINLOOP_STARTRPC     = "MainLoop, Processing StartRpc"
	LOG_MAINLOOP_STARTCLEANER = "MainLoop, Processing StartCleaner"
	LOG_MAINLOOP_SLEEP        = "Sleeping in main loop"
	LOG_START_CLEANER         = "Starting Session Cleaner"
	LOG_RESTART_CLEANER       = "Restarting Session Cleaner"
	LOG_PANIC_CAPTURED        = "Panic captured: %s!" //panic details
	LOG_RPC_EXIT              = "Rpc Webserver exit gracefully for some reason!"
	LOG_CLEANER_EXIT          = "Session Cleaner exit gracefully for some reason!"
	LOG_CREATE_WEBSERVER      = "Creating WebServer"
	LOG_RUNCLEAN              = "Running Cleaner"
	LOG_GOTDBERROR            = "Received error from db connection, exiting!"
	LOG_CONFDUMP              = "INIT of configurator and multiLogger done, dumping configuration"
)

// multiLogger
const (
	//early logging lines - before logger is setup
	EARLY_CREATING_HANDLER = "LOADLOGGER: DEBUG Creating log handler"
	EARLY_ENTER_CONFLOOP   = "LOADLOGGER: DEBUG Entering configuration loop"
	EARLY_LOOP_POS         = "LOADLOGGER: DEBUG Run %d Setting log level\n"    //run number in for loop
	EARLY_CHECKING_DEST    = "LOADLOGGER: DEBUG Run %d Checking destination\n" //run number in loop
	EARLY_DIALLING_SYSLOG  = "LOADLOGGER: DEBUG Run %d Dialling syslog: %s\n"  //run number in loop, syslog destination
	EARLY_LOGGER_FAIL      = "ERROR: Could not initialize logger:%s\n"         //error details
	EARLY_CREATE_STDOUT    = "LOADLOGGER: DEBUG Run %d Creating stdout logger\n"
	EARLY_CREATE_STDERR    = "LOADLOGGER: DEBUG Run %d Creating stderr logger\n"
	EARLY_SET_DEVLOG       = "LOADLOGGER: DEBUG Run %d Dialling syslog:devlog\n"
	EARLY_APPEND_ARR       = "LOADLOGGER: DEBUG Run %d Appending array\n"
)
