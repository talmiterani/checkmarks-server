package posts

import (
	"awesomeProject/internal/api/common/access"
	"awesomeProject/internal/utils"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	*access.DbConnections
	service *Service
}

func NewHandler(sdc *access.DbConnections) *Handler {

	h := Handler{
		sdc,
		NewService(sdc),
	}

	return &h
}

func (h *Handler) Init(router *mux.Router) { //, auth security.AuthHandler) {
	router.HandleFunc("/v1/movies", h.getMovies).Methods("GET")
}
func (h *Handler) getMovies(w http.ResponseWriter, r *http.Request) {

	res, err := h.service.getMovies(r.Context())
	if err != nil {
		fmt.Println("get movies error: ", err)
		w.WriteHeader(http.StatusInternalServerError)

		utils.WriteJson(err, w)
		return
	}

	utils.WriteJson(res, w)
}
