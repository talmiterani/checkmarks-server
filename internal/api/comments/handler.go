package comments

import (
	"checkmarks/internal/api/comments/models"
	"checkmarks/internal/api/common/access"
	"checkmarks/internal/utils"
	"encoding/json"
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
	router.HandleFunc("/v1/comments/{postId}", h.getAll).Methods("GET")
	router.HandleFunc("/v1/comments", h.add).Methods("POST")
	router.HandleFunc("/v1/comments", h.update).Methods("PUT")
	router.HandleFunc("/v1/comments/{commentId}", h.delete).Methods("DELETE")
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {

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

func (h *Handler) add(w http.ResponseWriter, r *http.Request) {

	var (
		comment models.Comment
	)

	if decodeErr := json.NewDecoder(r.Body).Decode(&comment); decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(decodeErr.Error()))
		return
	}

	comment.Prepare()

	if valRes := comment.Validate(true, true, false); len(valRes) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(valRes))
		return
	}

	res, err := h.service.add(r.Context(), &comment)

	if err != nil {
		fmt.Println("add comment error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	utils.WriteJson(res, w)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {

	var (
		comment models.Comment
	)

	if decodeErr := json.NewDecoder(r.Body).Decode(&comment); decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(decodeErr.Error()))
		return
	}

	comment.Prepare()
	if valRes := comment.Validate(false, false, true); len(valRes) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(valRes))
		return
	}

	res, err := h.service.update(r.Context(), &comment)

	if err != nil {
		fmt.Println("update comment error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	utils.WriteJson(res, w)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {

	var (
		params    = mux.Vars(r)
		commentId = params["commentId"]
	)

	if err := h.service.delete(r.Context(), commentId); err != nil {
		fmt.Println("delete post error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	w.Write(utils.OK)
}
