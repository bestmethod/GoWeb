package main

import (
	"github.com/bestmethod/goweb"
	"github.com/julienschmidt/httprouter"
	"net/http"
        "html/template"
        "fmt"
)

var ws *goweb.Webserver

func main() {
	ws = goweb.Init()
	router := httprouter.New()
	router.GET("/", Index)
	ws.Start(router)
}

type index struct {
	Username string
	Title    *string
	Subtitle string
}

func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	model := new(index)
	model.Username = "Robert"
	model.Title = ws.Config.Website.Name
	model.Subtitle = " - You are logged in!"
	RpcParse("index", "index.html", w, model)
}

func RpcParse(tName string, tFile string, w http.ResponseWriter, model interface{}) {
	t := template.New(tName)
	var err error
	t, err = t.ParseFiles(tFile)
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

