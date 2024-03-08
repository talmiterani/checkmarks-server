package internal

import "github.com/gorilla/mux"

type AppHandler interface {
	Init(router *mux.Router) //, auth security.AuthHandler)
}

type ApiHandler interface {
	Init(router *mux.Router)
}
