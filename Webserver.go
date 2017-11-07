package main

import (
	"./configurator"
	"./multiLogger"
	"./rpcDb"
	"./rpcListener"
	"database/sql"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
	"time"
)

func main() {
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

	//it seems that DB should have own log header, really
	dbLogger := new(multiLogger.LogHandler)
	*dbLogger = *logger
	dbLogger.Header = multiLogger.LOGHEADER_DB

	//get DB connection for RPC
	logger.Info(multiLogger.LOG_GETTING_DB_CONN)
	var dbConn *sql.DB
	defer func() {
		dbConn.Close()
	}()
	var err error
	dbConn, err = rpcDb.Connect(config.Database, dbLogger)
	if err != nil {
		logger.Fatal(multiLogger.LOG_GOTDBERROR)
	}

	logger.Debug(multiLogger.LOG_ENTER_MAIN_LOOP)
	//THE LOOP
	for {

		//handle rpc worker
		if config.Listener.RpcListenerRun == true {
			if config.General.DebugMainLoop == true {
				logger.Debug(multiLogger.LOG_MAINLOOP_STARTRPC)
			}
			if len(rpcWorker) > 0 {
				rpcWorkerVal = <-rpcWorker
				if rpcWorkerVal == 1 {
					logger.Info(multiLogger.LOG_START_RPC)
				} else if rpcWorkerVal == 0 {
					logger.Warn(multiLogger.LOG_RESTART_RPC)
				}
				go StartRpcListener(newLogger, config.Listener, rpcWorker, dbConn, config.Website)
			}
		}

		//handle session cleaner
		if config.Listener.SessionCleanerRun == true {
			if config.General.DebugMainLoop == true {
				logger.Debug(multiLogger.LOG_MAINLOOP_STARTCLEANER)
			}
			if len(cleaner) > 0 {
				cleanerVal = <-cleaner
				if cleanerVal == 1 {
					logger.Info(multiLogger.LOG_START_CLEANER)
				} else if cleanerVal == 0 {
					logger.Warn(multiLogger.LOG_RESTART_CLEANER)
				}
				go StartSessionCleaner(newLogger2, dbConn, config.Listener.SessionCleanerIntervalSeconds, cleaner)
			}
		}

		// inform the user to be happy
		if startup == true {
			logger.Info(multiLogger.LOG_INITIAL_DONE)
			startup = false
		}

		// sleep between probes
		if config.General.DebugMainLoop == true {
			logger.Debug(multiLogger.LOG_MAINLOOP_SLEEP)
		}
		time.Sleep(config.General.MonitorSleepSeconds * time.Second)
	}
	//we can use this when needed: dbConn := config.Database.Connect()
}

func StartRpcListener(logger *multiLogger.LogHandler, rpcConf *configurator.RpcConf, rpcWorker chan int, dbConn *sql.DB, webConf *configurator.WebConf) {
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
	ws := new(rpcListener.WebServer)
	newLogger := new(multiLogger.LogHandler)
	*newLogger = *logger
	newLogger.Header = multiLogger.LOGHEADER_RPCWEB
	ws.Logger = newLogger
	ws.RpcConf = rpcConf
	ws.DbConn = dbConn
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
		rpcListener.SessionCleaner(logger, dbConn)
	}
}
