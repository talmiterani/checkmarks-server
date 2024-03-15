package internal

import "github.com/gorilla/mux"

type AppHandler interface {
	Init(router *mux.Router)
}
