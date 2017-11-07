package rpcListener

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/russross/meddler"
	"html/template"
	"net/http"
	"time"
)

type Register struct {
	Password            string
	Username            string
	Title               *string
	Subtitle            string
	RegisterFailExists  bool
	RegisterFailMissing bool
}

func register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Session := getSession(w, r)
	if Session.UserId != 0 {
		http.Redirect(w, r, "/", 303)
		return
	}
	model := new(Register)
	model.Title = WebsiteConf.Name
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
	var err error
	t := template.New("register")
	t, _ = t.ParseFiles("templates/register.html")
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

//register user, IF does not exist
func runRegister(user string, pass string, Session *SessionStruct) bool {
	User := new(UserStruct)
	err := meddler.QueryRow(dbConn, User, "select * from users where username = ?", user, pass)
	var res bool
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error(fmt.Sprintf("Could not register user in the DB: %s\n", err))
		res = false
	} else if err != nil {
		//user does not exist, creating

		User.Username = user
		User.Password = pass
		User.LastLogin = 0
		User.Registered = time.Now().Unix()
		err := meddler.Insert(dbConn, "users", User)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not create user in the DB: %s\n", err))
			res = false
		} else {
			res = true
		}
	} else {
		res = false
	}
	updateSession(Session)
	return res
}
