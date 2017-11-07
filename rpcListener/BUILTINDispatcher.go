package rpcListener

import (
	"github.com/julienschmidt/httprouter"
)

func dispatch(router *httprouter.Router) {
	logger.Debug(LOG_CONF_DISPATCHERS)
	router.GET("/", index)
	router.GET("/login", login)
	router.POST("/login", login)
	router.GET("/register", register)
	router.POST("/register", register)
}
