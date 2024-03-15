package main

import (
	"checkmarks/internal"
	"checkmarks/internal/api/comments"
	"checkmarks/internal/api/common/access"
	"checkmarks/internal/api/posts"
	"checkmarks/internal/api/users"
	"checkmarks/internal/config"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {

	c := config.GetConfig()

	sdc, err := access.NewServiceDbConnections(c)

	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	appHandlers := []internal.AppHandler{
		posts.NewHandler(sdc),
		comments.NewHandler(sdc),
		users.NewHandler(sdc),
	}

	for i := 0; i < len(appHandlers); i++ {
		appHandlers[i].Init(router)
	}

	fmt.Println("done init all proper handlers")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "ACCEPT", "CONTENT-TYPE", "X-CB-EnvDb", "token"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE"})

	srv := &http.Server{
		Handler: handlers.CORS(headersOk, methodsOk)(
			handlers.CompressHandler(
				router,
			),
		),
		Addr:         ":" + strconv.Itoa(c.Server.Port),
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
	}

	fmt.Println("server listening on: ", c.Server.Port)

	log.Fatal(srv.ListenAndServe())
}
