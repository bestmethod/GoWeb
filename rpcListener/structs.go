package rpcListener

import (
	"../configurator"
	"../multiLogger"
	"database/sql"
)

//webserver basic struct with configs
type WebServer struct {
	Logger  *multiLogger.LogHandler
	RpcConf *configurator.RpcConf
	DbConn  *sql.DB
}

//session struct, to have cookies, session ID, whateva
type SessionStruct struct {
	Id             int64  `meddler:"id,pk"`
	UserId         int64  `meddler:"user_id"`
	SessionId      string `meddler:"session_id"`
	Expires        int64  `meddler:"expires"`
	KeepMeLoggedIn bool   `meddler:"keep_logged_in"`
}

//db struct for users
type UserStruct struct {
	UserId     int64  `meddler:"user_id,pk"`
	Username   string `meddler:"username"`
	Password   string `meddler:"password"`
	Registered int64  `meddler:"registered"`
	LastLogin  int64  `meddler:"last_login"`
}
