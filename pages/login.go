package pages

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/russross/meddler"
	"fmt"
	"time"
	"github.com/bestmethod/GoWeb"
)

type login struct {
	Password        string
	Username        string
	Title           *string
	Subtitle        string
	LoginFail       bool
	RegisterSuccess bool
	KeepMeLoggedIn  bool
}

func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Session := goweb.GetSession(w, r)
	if Session.UserId != 0 {
		http.Redirect(w, r, "/", 303)
		return
	}
	model := new(login)
	model.Title = ws.Config.Website.Name
	model.Subtitle = " - Login"
	model.KeepMeLoggedIn = false
	model.RegisterSuccess = false
	if r.Method == "GET" {
		model.Username = "username"
		model.Password = "password"
		model.LoginFail = false
		if r.FormValue("reg") == "success" {
			model.RegisterSuccess = true
		}
	} else {
		model.Username = r.PostFormValue("username")
		model.Password = r.PostFormValue("password")
		model.LoginFail = true
		if r.PostFormValue("goweb.KeepMeLoggedIn") == "on" {
			model.KeepMeLoggedIn = true
		}
		if checkLogin(model.Username, model.Password, Session, model.KeepMeLoggedIn) == true {
			http.Redirect(w, r, "/", 303)
			return
		}
	}
	RpcParse("login", "templates/login.html", w, model)
}

// login check and update user session
func checkLogin(user string, pass string, Session *goweb.SessionStruct, KeepMeLoggedIn bool) bool {
	User := new(UserStruct)
	err := meddler.QueryRow(ws.DbConn, User, "select * from users where username = ? and password = ?", user, pass)
	var res bool
	if err != nil && err.Error() != "sql: no rows in result set" {
		ws.Logger.Error(fmt.Sprintf("Could not login user in the DB: %s\n", err))
		res = false
	} else if err != nil {
		//user does not exist
		res = false
	} else {
		Session.UserId = User.UserId
		User.LastLogin = time.Now().Unix()
		err := meddler.Update(ws.DbConn, "users", User)
		if err != nil {
			ws.Logger.Error(fmt.Sprintf("Could not update user last login in the DB: %s\n", err))
		}
		goweb.KeepMeLoggedIn(Session, KeepMeLoggedIn)
		res = true
	}
	if res == false {
		goweb.UpdateSession(Session)
	}
	return res
}