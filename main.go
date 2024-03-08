package main

import (
	"awesomeProject/internal"
	"awesomeProject/internal/api/common/access"
	"awesomeProject/internal/api/posts"
	"awesomeProject/internal/config"
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
	}

	for i := 0; i < len(appHandlers); i++ {
		appHandlers[i].Init(router) //, authHandler)
	}

	// api handlers -
	// to be used with no auth by EM APi's
	//apiHandlers := []internal.ApiHandler{
	//	matching.NewCategoriesApiHandler(sdc),
	//}

	fmt.Println("done init all proper handlers")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "ACCEPT", "CONTENT-TYPE", "X-CB-EnvDb"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE"})

	srv := &http.Server{
		Handler: handlers.CORS(headersOk, methodsOk)(
			handlers.CompressHandler(
				//middlewares.EnvMiddleW(
				//	authHandler.Authenticate(
				//		middlewares.Finalize(config, sdc)(
				router,
				//		),
				//	),
				//),
			),
		),
		Addr:         ":" + strconv.Itoa(c.Server.Port),
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
	fmt.Println("server listening on: ", c.Server.Port)
}
