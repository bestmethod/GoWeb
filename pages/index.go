package pages

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"strconv"
	"github.com/bestmethod/GoWeb"
)

type index struct {
	Username string
	Title    *string
	Subtitle string
}

func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Session := goweb.GetSession(w, r)
	if r.FormValue("logout") == "true" {
		Session.UserId = 0
		goweb.UpdateSession(Session)
	}
	if Session.UserId == 0 {
		http.Redirect(w, r, "/login", 307)
		return
	}
	model := new(index)
	model.Username = strconv.Itoa(int(Session.UserId))
	model.Title = ws.Config.Website.Name
	model.Subtitle = " - You are logged in!"
	RpcParse("index", "templates/index.html", w, model)
}
