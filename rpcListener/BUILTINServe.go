package rpcListener

import (
	"../configurator"
	"../multiLogger"
	"database/sql"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/julienschmidt/httprouter"
	"github.com/leonelquinteros/gorand"
	"github.com/russross/meddler"
	"net/http"
	"time"
)

/*var (
	templateMap = template.FuncMap{
		"Upper": func(s string) string {
			return strings.ToUpper(s)
		},
	}
)*/

//var templates = template.New("templates/login.html").Funcs(templateMap)

//some Globals, because we need them :/
var cookieExpires time.Duration
var dbConn *sql.DB
var logger *multiLogger.LogHandler
var sessionExpires time.Duration
var WebsiteConf *configurator.WebConf
var rpcConf *configurator.RpcConf

//main webserver serve - set routes, global params and listen to connections
func (ws *WebServer) Serve(wc *configurator.WebConf) {
	ws.Logger.Debug(LOG_CONFIG_GLOBALS)
	dbConn = ws.DbConn
	logger = ws.Logger
	WebsiteConf = wc
	rpcConf = ws.RpcConf
	logger.Debug(LOG_ROUTER_MAKE)
	router := httprouter.New()
	dispatch(router)
	cookieExpires = time.Duration(ws.RpcConf.CookieLifetimeSeconds) * time.Second
	sessionExpires = time.Duration(ws.RpcConf.SessionExpireSeconds) * time.Second
	logger.Info(fmt.Sprintf(LOG_STARTLISTEN, ws.RpcConf.ListenIp, ws.RpcConf.ListenPort, ws.RpcConf.UseSSL))
	if ws.RpcConf.UseSSL == false {
		http.ListenAndServe(fmt.Sprintf("%s:%d", ws.RpcConf.ListenIp, ws.RpcConf.ListenPort), router)
	} else {
		http.ListenAndServeTLS(fmt.Sprintf("%s:%d", ws.RpcConf.ListenIp, ws.RpcConf.ListenPort), ws.RpcConf.SSLCrtPath, ws.RpcConf.SSLKeyPath, router)
	}
}

//helper function for endpoints - call at any point to get the session struct for a given cookie (or set the cookie)
func getSession(w http.ResponseWriter, r *http.Request) *SessionStruct {
	var sessionId *http.Cookie
	var err error
	var Session *SessionStruct
	sessionId, err = r.Cookie("session-id")
	if err != nil {
		// no cookie found
		Session = newSession(w, r)
	} else {
		Session = new(SessionStruct)
		err := meddler.QueryRow(dbConn, Session, "select * from session where session_id = ?", sessionId.Value)
		if err != nil && err.Error() != "sql: no rows in result set" {
			logger.Error(fmt.Sprintf(LOG_DB_QUERY_FAIL, err))
			Session = newSession(w, r)
		} else if err != nil {
			// cookie not found in DB
			Session = newSession(w, r)
		} else {
			if Session.Expires < time.Now().Unix() {
				//session expired, start new session
				Session = newSession(w, r)
			} else {
				var expire time.Time
				expire = time.Now().Add(cookieExpires)
				sessionId.Expires = expire
				sessionId.Path = "/"
				http.SetCookie(w, sessionId)
				if Session.UserId == 0 || Session.KeepMeLoggedIn == false {
					Session.Expires = int64(time.Now().Add(sessionExpires).Unix())
				} else {
					//session never expires
					Session.Expires = expire.Unix()
				}
				updateSession(Session)
			}
		}
	}
	return Session
}

//internal for getSession, no need to use otherwise
func newSession(w http.ResponseWriter, r *http.Request) *SessionStruct {
	var sessionId *http.Cookie
	var Session *SessionStruct
	uuid, _ := gorand.UUID()
	sessionId = new(http.Cookie)
	sessionId.Name = "session-id"
	sessionId.Value = uuid
	expires := time.Now().Add(cookieExpires)
	sessionId.Expires = expires
	sessionId.Path = "/"
	http.SetCookie(w, sessionId)
	Session = new(SessionStruct)
	Session.SessionId = uuid
	Session.UserId = 0
	Session.KeepMeLoggedIn = false
	Session.Expires = int64(time.Now().Add(sessionExpires).Unix())
	err := meddler.Insert(dbConn, "session", Session)
	if err != nil {
		logger.Error(fmt.Sprintf(LOG_DB_INSERT_FAIL, err))
	}
	if rpcConf.SessionDebug == true {
		logger.Debug(spew.Sdump(Session))
		logger.Debug(fmt.Sprintf("Expires: %s", time.Unix(Session.Expires, 0)))
	}
	return Session
}

//internal for getSession, no need to use otherwise
func updateSession(Session *SessionStruct) {
	err := meddler.Update(dbConn, "session", Session)
	if err != nil {
		logger.Error(fmt.Sprintf(LOG_DB_UPDATE_FAIL, err))
	}
	if rpcConf.SessionDebug == true {
		logger.Debug(spew.Sdump(Session))
		logger.Debug(fmt.Sprintf("Expires: %s", time.Unix(Session.Expires, 0)))
	}
}

//switch keepmeloggedin to true
func keepMeLoggedIn(Session *SessionStruct, yesNo bool) {
	if yesNo == true {
		Session.KeepMeLoggedIn = true
		Session.Expires = time.Now().Add(cookieExpires).Unix()
	} else {
		Session.KeepMeLoggedIn = false
		Session.Expires = int64(time.Now().Add(sessionExpires).Unix())
	}
	updateSession(Session)
}
