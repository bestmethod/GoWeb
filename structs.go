package goweb

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
)

//BUILTIN webserver basic struct with configs
type WebServer struct {
	Logger  *LogHandler
	RpcConf *RpcConf
	DbConn  *sql.DB
	Router  *httprouter.Router
}

//session struct, to have cookies, session ID, whateva
type SessionStruct struct {
	Id             int64  `meddler:"id,pk"`
	UserId         int64  `meddler:"user_id"`
	SessionId      string `meddler:"session_id"`
	SessionKey     string `meddler:"session_key"`
	Expires        int64  `meddler:"expires"`
	KeepMeLoggedIn bool   `meddler:"keep_logged_in"`
}
