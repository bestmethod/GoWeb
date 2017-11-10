package main

import (
	"github.com/bestmethod/GoWeb"
	"./pages"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	ws := goweb.Init()
	pages.Init(ws)
	router := httprouter.New()
	router.GET("/", pages.Index)
	router.GET("/login", pages.Login)
	router.GET("/register", pages.Register)
	router.POST("/login", pages.Login)
	router.POST("/register", pages.Register)
	router.ServeFiles("/images/*filepath", http.Dir("./templates/images"))
	ws.Start(router)
}
