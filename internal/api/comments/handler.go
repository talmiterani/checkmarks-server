package comments

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
	router.HandleFunc("/v1/{postId}/comments", h.getComments).Methods("GET")
}
func (h *Handler) getComments(w http.ResponseWriter, r *http.Request) {

	var (
		params = mux.Vars(r)
		postId = params["postId"]
	)

	res, err := h.service.getComments(r.Context(), postId)
	if err != nil {
		fmt.Println("get movies error: ", err)
		w.WriteHeader(http.StatusInternalServerError)

		utils.WriteJson(err, w)
		return
	}

	utils.WriteJson(res, w)
}
