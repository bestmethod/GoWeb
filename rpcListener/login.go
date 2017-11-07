package rpcListener

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/russross/meddler"
	"html/template"
	"net/http"
	"time"
)

type Login struct {
	Password        string
	Username        string
	Title           *string
	Subtitle        string
	LoginFail       bool
	RegisterSuccess bool
	KeepMeLoggedIn  bool
}

func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Session := getSession(w, r)
	if Session.UserId != 0 {
		http.Redirect(w, r, "/", 303)
		return
	}
	model := new(Login)
	model.Title = WebsiteConf.Name
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
		if r.PostFormValue("KeepMeLoggedIn") == "on" {
			model.KeepMeLoggedIn = true
		}
		if checkLogin(model.Username, model.Password, Session, model.KeepMeLoggedIn) == true {
			http.Redirect(w, r, "/", 303)
			return
		}
	}
	var err error
	t := template.New("login")
	t, _ = t.ParseFiles("templates/login.html")
	if err != nil {
		logger.Error(fmt.Sprintf("There was an error serving page template: %s",err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t, err = t.ParseGlob("templates/common/*")
	if err != nil {
		logger.Error(fmt.Sprintf("There was an error serving common template: %s",err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = t.Execute(w, &model)
	if err != nil {
		logger.Error(fmt.Sprintf("There was an error executing templates: %s",err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// login check and update user session
func checkLogin(user string, pass string, Session *SessionStruct, KeepMeLoggedIn bool) bool {
	User := new(UserStruct)
	err := meddler.QueryRow(dbConn, User, "select * from users where username = ? and password = ?", user, pass)
	var res bool
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error(fmt.Sprintf("Could not login user in the DB: %s\n", err))
		res = false
	} else if err != nil {
		//user does not exist
		res = false
	} else {
		Session.UserId = User.UserId
		User.LastLogin = time.Now().Unix()
		err := meddler.Update(dbConn, "users", User)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not update user last login in the DB: %s\n", err))
		}
		keepMeLoggedIn(Session, KeepMeLoggedIn)
		res = true
	}
	if res == false {
		updateSession(Session)
	}
	return res
}
