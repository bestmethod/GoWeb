package pages

import (
	"github.com/bestmethod/GoWeb"
	"net/http"
	"fmt"
	"html/template"
)

//db struct for users
type UserStruct struct {
	UserId     int64  `meddler:"user_id,pk"`
	Username   string `meddler:"username"`
	Password   string `meddler:"password"`
	Registered int64  `meddler:"registered"`
	LastLogin  int64  `meddler:"last_login"`
}

var ws *goweb.Webserver

func Init(web *goweb.Webserver) {
	ws = web
}

func RpcParse(tName string, tFile string, w http.ResponseWriter, model interface{}) {
	t := template.New(tName)
	var err error
	t, err = t.ParseFiles(tFile)
	if err != nil {
		ws.Logger.Error(fmt.Sprintf("There was an error serving page template: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t, err = t.ParseGlob("templates/common/*")
	if err != nil {
		ws.Logger.Error(fmt.Sprintf("There was an error serving common template: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = t.Execute(w, &model)
	if err != nil {
		ws.Logger.Error(fmt.Sprintf("There was an error executing templates: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
