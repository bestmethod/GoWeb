package rpcListener

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"strconv"
)

type Index struct {
	Username string
	Title    *string
	Subtitle string
}

func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Session := getSession(w, r)
	if r.FormValue("logout") == "true" {
		Session.UserId = 0
		updateSession(Session)
	}
	if Session.UserId == 0 {
		http.Redirect(w, r, "/login", 307)
		return
	}
	model := new(Index)
	model.Username = strconv.Itoa(int(Session.UserId))
	model.Title = WebsiteConf.Name
	model.Subtitle = " - You are logged in!"
	var err error
	t := template.New("index")
	t, _ = t.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("Error1:", err)
		return
	}
	t, _ = t.ParseGlob("templates/common/*")
	if err != nil {
		fmt.Println("Error2:", err)
		return
	}
	err = t.Execute(w, &model)
	if err != nil {
		fmt.Println("Error3:", err)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
