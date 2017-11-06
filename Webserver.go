package main

import (
	"./configurator"
	"./multiLogger"
	"./rpcListener"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

//TODO: move ALL messages out of Webserver.go, logger and configurator to LogMessages constants
//TODO: add proper debug messages to rpcListener

func main() {
	// check if arg with config file was provided
	fmt.Println("MAIN: Processing command line parameters")
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
		fmt.Printf("MAIN: DEBUG Early Debug enabled\nMAIN: DEBUG Conf file: %s\n", cfile)
	}
	fmt.Println("MAIN: Loading configuration")
	config := configurator.Load(cfile, earlyDebug)
	fmt.Println("MAIN: Loading logger")
	logger := multiLogger.Init(config.Loggers, earlyDebug, config.General.ServiceName)
	logger.Header = "MAIN"
	// inform the user to be excited
	logger.Info(multiLogger.LOG_BASIC_LOAD_DONE)
	startup := true

	//chans for rpcWorker and sqsQueue
	logger.Debug("Configuring channels for goroutines")
	var rpcWorker = make(chan int, 1)
	rpcWorker <- 1
	rpcWorkerVal := 0
	var cleaner = make(chan int, 1)
	cleaner <- 1
	cleanerVal := 0

	//so that gorotine ends up having a different name
	newLogger := new(multiLogger.LogHandler)
	*newLogger = *logger
	newLogger.Header = "RPC-GOROUTINE"
	newLogger2 := new(multiLogger.LogHandler)
	*newLogger2 = *logger
	newLogger2.Header = "CLEANER-GOROUTINE"

	//get DB connection for RPC
	logger.Info("Getting Database Connection for Rpc")
	var dbConn *sql.DB
	defer func() {
		dbConn.Close()
	}()
	dbConn = config.Database.Connect(earlyDebug)

	logger.Debug("Entering main loop")
	//THE LOOP
	for {

		//handle rpc worker
		if config.Listener.RpcListenerRun == true {
			if config.General.DebugMainLoop == true {
				logger.Debug("startRpc true, processing")
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
				logger.Debug("cleaner true, processing")
			}
			if len(cleaner) > 0 {
				cleanerVal = <-cleaner
				if cleanerVal == 1 {
					logger.Info("Starting session cleaner")
				} else if cleanerVal == 0 {
					logger.Warn("Restarting session cleaner")
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
			logger.Debug("Sleeping in main loop")
		}
		time.Sleep(config.General.MonitorSleepSeconds * time.Second)
	}
	//we can use this when needed: dbConn := config.Database.Connect()
}

func StartRpcListener(logger *multiLogger.LogHandler, rpcConf *configurator.RpcConf, rpcWorker chan int, dbConn *sql.DB, webConf *configurator.WebConf) {
	defer func() {
		r := recover()
		if r != nil {
			logger.Error(fmt.Sprintf("Panic captured: %s!", r))
		} else {
			logger.Error(fmt.Sprintf("Rpc Webserver exit gracefully for some reason!"))
		}
		rpcWorker <- 0
	}()
	logger.Debug("Creating WebServer")
	ws := new(rpcListener.WebServer)
	newLogger := new(multiLogger.LogHandler)
	*newLogger = *logger
	newLogger.Header = "RPC-WEBSERVER"
	ws.Logger = newLogger
	ws.RpcConf = rpcConf
	ws.DbConn = dbConn
	ws.Serve(webConf)
}

func StartSessionCleaner(logger *multiLogger.LogHandler, dbConn *sql.DB, cleanerSleep int, cleaner chan int) {
	defer func() {
		r := recover()
		if r != nil {
			logger.Error(fmt.Sprintf("Panic captured: %s!", r))
		} else {
			logger.Debug(fmt.Sprintf("Rpc Cleaner exit gracefully."))
		}
		cleaner <- 0
	}()
	newLogger := new(multiLogger.LogHandler)
	*newLogger = *logger
	newLogger.Header = "RPC-SESSIONCLEANER"
	cleanerSleepTime := time.Duration(cleanerSleep) * time.Second
	for {
		time.Sleep(cleanerSleepTime)
		logger.Debug("Running Cleaner")
		rpcListener.SessionCleaner(logger, dbConn)
	}
}
