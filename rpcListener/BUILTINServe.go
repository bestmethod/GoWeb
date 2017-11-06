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

//main webserver serve - set routes, global params and listen to connections
func (ws *WebServer) Serve(wc *configurator.WebConf) {
	dbConn = ws.DbConn
	logger = ws.Logger
	WebsiteConf = wc
	router := httprouter.New()
	dispatch(router)
	cookieExpires = time.Duration(ws.RpcConf.CookieLifetimeSeconds) * time.Second
	sessionExpires = time.Duration(ws.RpcConf.SessionExpireSeconds) * time.Second
	//TODO: parse whole RpcConf (including SSL config) here
	http.ListenAndServe(fmt.Sprintf("%s:%d", ws.RpcConf.ListenIp, ws.RpcConf.ListenPort), router)
}

//helper function for endpoints - call at any point to get the session struct for a given cookie (or set the cookie)
func getSession(w http.ResponseWriter, r *http.Request) *SessionStruct {
	var sessionId *http.Cookie
	var err error
	var Session *SessionStruct
	sessionId, err = r.Cookie("session-id")
	if err != nil {
		logger.Debug("No cookie found")
		Session = newSession(w, r)
	} else {
		Session = new(SessionStruct)
		err := meddler.QueryRow(dbConn, Session, "select * from session where session_id = ?", sessionId.Value)
		if err != nil && err.Error() != "sql: no rows in result set" {
			logger.Error(fmt.Sprintf("Could not update cookie session data in the DB: %s\n", err))
			Session = newSession(w, r)
		} else if err != nil {
			logger.Debug("Cookie not found in DB")
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
	logger.Debug("Entering newSession")
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
		logger.Error(fmt.Sprintf("Could not insert cookie session data to the DB: %s\n", err))
	}
	logger.Debug("Exiting newSession")
	return Session
}

//internal for getSession, no need to use otherwise
func updateSession(Session *SessionStruct) {
	logger.Debug("Entering updateSession")
	err := meddler.Update(dbConn, "session", Session)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not update cookie session data in the DB: %s\n", err))
	}
	logger.Debug(spew.Sdump(Session))
	logger.Debug(fmt.Sprintf("Expires: %s", time.Unix(Session.Expires, 0)))
	logger.Debug("Exiting updateSession")
}

//switch keepmeloggedin to true
func keepMeLoggedIn(Session *SessionStruct, yesNo bool) {
	logger.Debug("Entering KeepMeLoggedIn")
	if yesNo == true {
		Session.KeepMeLoggedIn = true
		Session.Expires = time.Now().Add(cookieExpires).Unix()
	} else {
		Session.KeepMeLoggedIn = false
		Session.Expires = int64(time.Now().Add(sessionExpires).Unix())
	}
	updateSession(Session)
	logger.Debug("Exiting KeepMeLoggedIn")
}
