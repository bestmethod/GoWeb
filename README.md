# GoWeb
GoWeb framework provide a ncie and easy way to create a webserver. It handles configuration from a text file, sets up listener, database connections, auto-restart of goroutines and an amazing and flexible multi-destination logger. All that remains to be done is a few lines of code to start it and your own templates/handlers using the httprouter.

## Quickstart guide
Feel free to replace `loginExample` in the git clone command with `simpleExample`
```
go get github.com/bestmethod/GoWeb
cd ~
mkdir GoWebLoginExample
cd GoWebLoginExample
git clone -b loginExample https://github.com/bestmethod/GoWeb.git
cd GoWeb
go run main.go config_file.txt
```
Now connect from your browser to http://127.0.0.1:8080


## Session handling builtin
The builtin session handling uses session-id and session-key to ensure safe sessions without the possibility of session stealing. It also uses DB-backed sessions for horizontal scaling. It also has a builtin session cleaner, which periodically removes old sessions from the database.

`goweb.GetSession(w http.ResponseWriter, r *http.Request) *SessionStruct` <- automatically get a session (DB/Cookie based). Create new one if needed

`goweb.UpdateSession(Session *SessionStruct)` <- automatically update the session so that again there is one hour till expiry. Good to use at end of each call if you want to log user out after 1 hour or no activity on the site.

`goweb.NewSession(w http.ResponseWriter, r *http.Request) *SessionStruct` <- should get called by other functions, creates a new session and cookies are set

`goweb.KeepMeLoggedIn(Session *SessionStruct)` <- to be called if you want the session to expire after 20 years instead of 1 hour. Good if user selects "keep me loggee in" box.

## Example code using this import
Example code can be found in this repo, using the following branches:

`simpleExample` <- branch containing a simple example to build on (essentially the below)

`loginExample` <- branch containing an example of an index (with redirect), login, register pages (all working) together with session handling.

## The Webserver struct
goweb.Init() returns a Webserver struct, which contains loads of useful stuff. For most uses, the following is important:

`Webserver.Logger.(Init|Debug|Warn|Error|Critical)` <- pass this function a string to be logged

`Webserver.DbConn` <- an *sql.DB instance

`Webserver.Config.Website.Name` <- configured name of website from the configuration file

It is advisable to use meddler for database handling.

## Example configuration file
```
[General]
ServiceName="SomeServiceName"
DebugMainLoop=false
MonitorSleepSeconds=5
 
[Database]
# Type= sqlite3|MySQL
Type="sqlite3"
# can be an IP:PORT, or domain:port
Server="./exampleDb.sqlite3"
User="None"
Password="None"
DbName="None"
UseTLS=false
# the below 2 are mutually exclusive (and only required if the above is true). To use ssl_ca, ensure TLSSkipVerify=false
TLSSkipVerify=false
# ssl_ca="/some/path"
 
[Listener]
ListenIp="0.0.0.0"
ListenPort=8080
UseSSL=false
# SSLCrtPath="/some/path"
# SSLKeyPath="/some/path"
CookieLifetimeSeconds=315360000
SessionExpireSeconds=3600
SessionCleanerIntervalSeconds=10
SessionCleanerRun=true
RpcListenerRun=true
SessionDebug=false
 
[Website]
Name="SomeWebsiteName"
 
######### LOGGERS #########
 
[[Logger]]
LogLevel="INFO"
RpcLogLevel="DEBUG"
Destination="stdout"
 
# [[Logger]]
# LogLevel="DEBUG"
# RpcLogLevel="DEBUG"
# Destination="devlog"
 
# [[Logger]]
# LogLevel="ERROR"
# RpcLogLevel="ERROR"
# Destination="stderr"
 
# [[Logger]]
# LogLevel="INFO"
# RpcLogLevel="INFO"
# Destination="udp://127.0.0.1:11514"
```
## Database sql
```
CREATE TABLE session (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    session_id string(50),
    session_key string(130),
    expires INTEGER,
    keep_logged_in bool
);
```

## Other dependencies
```
go get github.com/BurntSushi/toml
go get github.com/mattn/go-sqlite3
go get github.com/russross/meddler
go get github.com/go-sql-driver/mysql
go get github.com/davecgh/go-spew/spew
go get github.com/julienschmidt/httprouter
go get github.com/leonelquinteros/gorand
```

## Simple usage example

#### main bits
```
# import the basics
import (
    "github.com/julienschmidt/httprouter"
    "net/http"
    "github.com/bestmethod/goweb"
    "strconv"
    "html/template"
    "fmt"
)
 
# main func
func main() {
 
    # create a goweb instance
    ws := goweb.Init()
 
    # create a http router and add an index page
    router := httprouter.New()
    router.GET("/", indexFunc)
 
    # start the router
    ws.Start(router)
}
```

#### indexFunc and struct - standard httprouter stuff
```
type index struct {
    Username string
    Title    *string
    Subtitle string
}

func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    model := new(index)
    model.Username = "Robert"
    model.Title = ws.Config.Website.Name
    model.Subtitle = " - Index Page!"
    
    t := template.New("index")
    var err error
    t, err = t.ParseFiles("index.html")
    if err != nil {
        ws.Logger.Error(fmt.Sprintf("There was an error serving page template: %s", err))
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    err = t.Execute(w, &model)
    if err != nil {
        ws.Logger.Error(fmt.Sprintf("There was an error executing templates: %s", err))
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

#### index.html to complement the indexFunc
```
{{ define "index" }}
<html><head><Title>{{.Title}}{{.Subtitle}}</Title></head>
<body>
<center>
Hello, {{.Username}}!<br>
<a href="?logout=true">Click here to logout</a>
</center>
</body></html
{{ end }}
```

## TODO
* README.md in *Example branches
* license file
* unit tests
