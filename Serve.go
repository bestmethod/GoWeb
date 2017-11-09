package goweb

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/davecgh/go-spew/spew"
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
var logger *LogHandler
var sessionExpires time.Duration
var WebsiteConf *WebConf
var rpcConf *RpcConf

//main webserver serve - set routes, global params and listen to connections
func (ws *WebServer) Serve(wc *WebConf) {
	ws.Logger.Debug(LOG_CONFIG_GLOBALS)
	dbConn = ws.DbConn
	logger = ws.Logger
	WebsiteConf = wc
	rpcConf = ws.RpcConf
	logger.Debug(LOG_ROUTER_MAKE)
	router := ws.Router
	cookieExpires = time.Duration(ws.RpcConf.CookieLifetimeSeconds) * time.Second
	sessionExpires = time.Duration(ws.RpcConf.SessionExpireSeconds) * time.Second
	logger.Info(fmt.Sprintf(LOG_STARTLISTEN, ws.RpcConf.ListenIp, ws.RpcConf.ListenPort, ws.RpcConf.UseSSL))
	var err error
	if ws.RpcConf.UseSSL == false {
		err = http.ListenAndServe(fmt.Sprintf("%s:%d", ws.RpcConf.ListenIp, ws.RpcConf.ListenPort), router)
	} else {
		err = http.ListenAndServeTLS(fmt.Sprintf("%s:%d", ws.RpcConf.ListenIp, ws.RpcConf.ListenPort), ws.RpcConf.SSLCrtPath, ws.RpcConf.SSLKeyPath, router)
	}
	if err != nil {
		panic(err)
	}
}

//helper function for endpoints - call at any point to get the session struct for a given cookie (or set the cookie)
func GetSession(w http.ResponseWriter, r *http.Request) *SessionStruct {
	var sessionId *http.Cookie
	var sessionKey *http.Cookie
	var err error
	var Session *SessionStruct
	sessionId, err = r.Cookie("session-id")
	if err == nil {
		sessionKey, err = r.Cookie("session-key")
	}
	if err != nil {
		// no cookie found
		Session = NewSession(w, r)
	} else {
		Session = new(SessionStruct)
		err := meddler.QueryRow(dbConn, Session, "select * from session where session_id = ? and session_key = ?", sessionId.Value, sessionKey.Value)
		if err != nil && err.Error() != "sql: no rows in result set" {
			logger.Error(fmt.Sprintf(LOG_DB_QUERY_FAIL, err))
			Session = NewSession(w, r)
		} else if err != nil {
			// cookie not found in DB
			Session = NewSession(w, r)
		} else {
			if Session.Expires < time.Now().Unix() {
				//session expired, start new session
				Session = NewSession(w, r)
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
				UpdateSession(Session)
			}
		}
	}
	return Session
}

//internal for getSession, no need to use otherwise
func NewSession(w http.ResponseWriter, r *http.Request) *SessionStruct {
	var sessionId *http.Cookie
	var sessionKey *http.Cookie
	var Session *SessionStruct
	uuid, _ := gorand.UUID()
	expires := time.Now().Add(cookieExpires)
	sessionId = new(http.Cookie)
	sessionId.Name = "session-id"
	sessionId.Value = uuid
	sessionId.Expires = expires
	sessionId.Path = "/"
	http.SetCookie(w, sessionId)
	key := make([]byte, 64)
	_, err := rand.Read(key)
	if err != nil {
		logger.Error(fmt.Sprintf(LOG_RAND_FAIL, err))
	}
	b64key := base64.StdEncoding.EncodeToString(key)
	sessionKey = new(http.Cookie)
	sessionKey.Name = "session-key"
	sessionKey.Value = b64key
	sessionKey.Expires = expires
	sessionKey.Path = "/"
	http.SetCookie(w, sessionKey)
	Session = new(SessionStruct)
	Session.SessionId = uuid
	Session.SessionKey = b64key
	Session.UserId = 0
	Session.KeepMeLoggedIn = false
	Session.Expires = int64(time.Now().Add(sessionExpires).Unix())
	err = meddler.Insert(dbConn, "session", Session)
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
func UpdateSession(Session *SessionStruct) {
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
func KeepMeLoggedIn(Session *SessionStruct, yesNo bool) {
	if yesNo == true {
		Session.KeepMeLoggedIn = true
		Session.Expires = time.Now().Add(cookieExpires).Unix()
	} else {
		Session.KeepMeLoggedIn = false
		Session.Expires = int64(time.Now().Add(sessionExpires).Unix())
	}
	UpdateSession(Session)
}
