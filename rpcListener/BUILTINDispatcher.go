package rpcListener

import (
	"github.com/julienschmidt/httprouter"
)

func dispatch(router *httprouter.Router) {
	router.GET("/", index)
	router.GET("/login", login)
	router.POST("/login", login)
	router.GET("/register", register)
	router.POST("/register", register)
}
