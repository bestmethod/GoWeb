package goweb

import (
	"./configurator"
	"./multiLogger"
	"database/sql"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/julienschmidt/httprouter"
	"log"
	"os"
	"time"
)

type Webserver struct {
	Logger       *multiLogger.LogHandler
	DbConn       *sql.DB
	Config       *configurator.Config
	rpcWorker    chan int
	rpcWorkerVal int
	newLogger    *multiLogger.LogHandler
	newLogger2   *multiLogger.LogHandler
	rpcWebLogger *multiLogger.LogHandler
	cleaner      chan int
	cleanerVal   int
	startup      bool
}

func Init() *Webserver {
	ws := new(Webserver)
	// check if arg with config file was provided
	fmt.Println(multiLogger.EARLY_LOADING_ARGS)
	if len(os.Args) < 2 {
		log.Fatalf(multiLogger.EARLY_ARGS_FAIL, os.Args[0])
	}
	// load configs and logger
	earlyDebug := false
	cfile := ""
	if len(os.Args) == 3 {
		if os.Args[1] == "--early-debug" || os.Args[2] == "--early-debug" {
			earlyDebug = true
		}
		if os.Args[1] == "--early-debug" {
			cfile = os.Args[2]
		} else {
			cfile = os.Args[1]
		}
	} else if len(os.Args) == 2 {
		cfile = os.Args[1]
	} else {
		log.Fatalf(multiLogger.EARLY_ARGS_FAIL, os.Args[0])
	}
	if earlyDebug == true {
		fmt.Printf(multiLogger.EARLY_DEBUG_ON, cfile)
	}
	fmt.Println(multiLogger.EARLY_LOADING_CONF)
	config := configurator.Load(cfile, earlyDebug)
	fmt.Println(multiLogger.EARLY_LOADING_LOGGER)
	logger := multiLogger.Init(config.Loggers, earlyDebug, config.General.ServiceName)
	logger.Header = multiLogger.LOGHEADER_MAIN
	logger.Debug(multiLogger.LOG_CONFDUMP)
	logger.Debug(spew.Sdump(config))
	// inform the user to be excited
	logger.Info(multiLogger.LOG_BASIC_LOAD_DONE)
	startup := true

	//chans for rpcWorker and sqsQueue
	logger.Debug(multiLogger.LOG_CONF_CHANNELS)
	var rpcWorker = make(chan int, 1)
	rpcWorker <- 1
	rpcWorkerVal := 0
	var cleaner = make(chan int, 1)
	cleaner <- 1
	cleanerVal := 0

	//so that gorotine ends up having a different name
	newLogger := new(multiLogger.LogHandler)
	*newLogger = *logger
	newLogger.Header = multiLogger.LOGHEADER_RPC
	newLogger2 := new(multiLogger.LogHandler)
	*newLogger2 = *logger
	newLogger2.Header = multiLogger.LOGHEADER_CLEANER
	rpcWebLogger := new(multiLogger.LogHandler)
	*rpcWebLogger = *logger
	rpcWebLogger.Header = multiLogger.LOGHEADER_RPCWEB
	ws.Logger = rpcWebLogger
	//it seems that DB should have own log header, really
	dbLogger := new(multiLogger.LogHandler)
	*dbLogger = *logger
	dbLogger.Header = multiLogger.LOGHEADER_DB

	//get DB connection for RPC
	logger.Info(multiLogger.LOG_GETTING_DB_CONN)
	var dbConn *sql.DB
	var err error
	dbConn, err = Connect(config.Database, dbLogger)
	if err != nil {
		logger.Fatal(multiLogger.LOG_GOTDBERROR)
	}
	ws.DbConn = dbConn
	ws.Config = config
	ws.rpcWorker = rpcWorker
	ws.rpcWorkerVal = rpcWorkerVal
	ws.newLogger = newLogger
	ws.rpcWebLogger = rpcWebLogger
	ws.cleaner = cleaner
	ws.cleanerVal = cleanerVal
	ws.newLogger2 = newLogger2
	ws.startup = startup
	return ws
}

func (ws *Webserver) Start(router *httprouter.Router) {
	ws.Logger.Debug(multiLogger.LOG_ENTER_MAIN_LOOP)
	//THE LOOP
	for {

		//handle rpc worker
		if ws.Config.Listener.RpcListenerRun == true {
			if ws.Config.General.DebugMainLoop == true {
				ws.Logger.Debug(multiLogger.LOG_MAINLOOP_STARTRPC)
			}
			if len(ws.rpcWorker) > 0 {
				ws.rpcWorkerVal = <-ws.rpcWorker
				if ws.rpcWorkerVal == 1 {
					ws.Logger.Info(multiLogger.LOG_START_RPC)
				} else if ws.rpcWorkerVal == 0 {
					ws.Logger.Warn(multiLogger.LOG_RESTART_RPC)
				}
				go StartRpcListener(ws.newLogger, ws.rpcWebLogger, ws.Config.Listener, ws.rpcWorker, ws.DbConn, ws.Config.Website, router)
			}
		}

		//handle session cleaner
		if ws.Config.Listener.SessionCleanerRun == true {
			if ws.Config.General.DebugMainLoop == true {
				ws.Logger.Debug(multiLogger.LOG_MAINLOOP_STARTCLEANER)
			}
			if len(ws.cleaner) > 0 {
				ws.cleanerVal = <-ws.cleaner
				if ws.cleanerVal == 1 {
					ws.Logger.Info(multiLogger.LOG_START_CLEANER)
				} else if ws.cleanerVal == 0 {
					ws.Logger.Warn(multiLogger.LOG_RESTART_CLEANER)
				}
				go StartSessionCleaner(ws.newLogger2, ws.DbConn, ws.Config.Listener.SessionCleanerIntervalSeconds, ws.cleaner)
			}
		}

		// inform the user to be happy
		if ws.startup == true {
			ws.Logger.Info(multiLogger.LOG_INITIAL_DONE)
			ws.startup = false
		}

		// sleep between probes
		if ws.Config.General.DebugMainLoop == true {
			ws.Logger.Debug(multiLogger.LOG_MAINLOOP_SLEEP)
		}
		time.Sleep(ws.Config.General.MonitorSleepSeconds * time.Second)
	}
	//we can use this when needed: dbConn := config.Database.Connect()
}

func StartRpcListener(logger *multiLogger.LogHandler, newLogger *multiLogger.LogHandler, rpcConf *configurator.RpcConf, rpcWorker chan int, dbConn *sql.DB, webConf *configurator.WebConf, router *httprouter.Router) {
	defer func() {
		r := recover()
		if r != nil {
			logger.Error(fmt.Sprintf(multiLogger.LOG_PANIC_CAPTURED, r))
		} else {
			logger.Error(fmt.Sprintf(multiLogger.LOG_RPC_EXIT))
		}
		rpcWorker <- 0
	}()
	logger.Debug(multiLogger.LOG_CREATE_WEBSERVER)
	ws := new(WebServer)
	ws.Logger = newLogger
	ws.RpcConf = rpcConf
	ws.DbConn = dbConn
	ws.Router = router
	ws.Serve(webConf)
}

func StartSessionCleaner(logger *multiLogger.LogHandler, dbConn *sql.DB, cleanerSleep int, cleaner chan int) {
	defer func() {
		r := recover()
		if r != nil {
			logger.Error(fmt.Sprintf(multiLogger.LOG_PANIC_CAPTURED, r))
		} else {
			logger.Error(fmt.Sprintf(multiLogger.LOG_CLEANER_EXIT))
		}
		cleaner <- 0
	}()
	newLogger := new(multiLogger.LogHandler)
	*newLogger = *logger
	newLogger.Header = multiLogger.LOGHEADER_CLNSESS
	cleanerSleepTime := time.Duration(cleanerSleep) * time.Second
	for {
		time.Sleep(cleanerSleepTime)
		logger.Debug(multiLogger.LOG_RUNCLEAN)
		SessionCleaner(logger, dbConn)
	}
}
