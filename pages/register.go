package pages

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/russross/meddler"
	"fmt"
	"time"
	"github.com/bestmethod/GoWeb"
)

type register struct {
	Password            string
	Username            string
	Title               *string
	Subtitle            string
	RegisterFailExists  bool
	RegisterFailMissing bool
}

func Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Session := goweb.GetSession(w, r)
	if Session.UserId != 0 {
		http.Redirect(w, r, "/", 303)
		return
	}
	model := new(register)
	model.Title = ws.Config.Website.Name
	model.Subtitle = " - Register"
	if r.Method == "GET" {
		model.Username = "username"
		model.Password = "password"
		model.RegisterFailExists = false
		model.RegisterFailMissing = false
	} else {
		model.Username = r.PostFormValue("username")
		model.Password = r.PostFormValue("password")
		if model.Username == "" || model.Username == "username" || model.Password == "" || model.Password == "password" {
			model.RegisterFailExists = false
			model.RegisterFailMissing = true
		} else {
			ret := runRegister(model.Username, model.Password, Session)
			if ret == true {
				http.Redirect(w, r, "/login?reg=success", 303)
				return
			} else if ret == false {
				model.RegisterFailExists = true
				model.RegisterFailMissing = false
			}
		}
	}
	RpcParse("register", "templates/register.html", w, model)
}

//register user, IF does not exist
func runRegister(user string, pass string, Session *goweb.SessionStruct) bool {
	User := new(UserStruct)
	err := meddler.QueryRow(ws.DbConn, User, "select * from users where username = ?", user, pass)
	var res bool
	if err != nil && err.Error() != "sql: no rows in result set" {
		ws.Logger.Error(fmt.Sprintf("Could not register user in the DB: %s\n", err))
		res = false
	} else if err != nil {
		//user does not exist, creating

		User.Username = user
		User.Password = pass
		User.LastLogin = 0
		User.Registered = time.Now().Unix()
		err := meddler.Insert(ws.DbConn, "users", User)
		if err != nil {
			ws.Logger.Error(fmt.Sprintf("Could not create user in the DB: %s\n", err))
			res = false
		} else {
			res = true
		}
	} else {
		res = false
	}
	goweb.UpdateSession(Session)
	return res
}

