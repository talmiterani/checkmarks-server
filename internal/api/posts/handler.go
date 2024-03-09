package posts

import (
	"awesomeProject/internal/api/common/access"
	"awesomeProject/internal/api/posts/models"
	"awesomeProject/internal/utils"
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
	router.HandleFunc("/v1/posts", h.getAll).Methods("GET")
	router.HandleFunc("/v1/posts", h.add).Methods("POST")
	router.HandleFunc("/v1/posts", h.update).Methods("PUT")
	router.HandleFunc("/v1/posts/{postId}", h.delete).Methods("DELETE")

}

//	var RegisterBookStoreRoutes = func(router *mux.Router) { //, auth security.AuthHandler) {
//		router.HandleFunc("/v1/posts", CreateBook).Methods("POST")
//	}
//
//	func CreateBook(writer http.ResponseWriter, r *http.Request) {
//		CreateBook := &models.Post{}
//		reqBody, _ := io.ReadAll(r.Body)
//		r.Body.Close()
//		err := json.Unmarshal(reqBody, CreateBook)
//		if err != nil {
//			return
//		}
//	}
func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {

	res, err := h.service.getAll(r.Context())

	if err != nil {
		fmt.Println("get posts error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	utils.WriteJson(res, w)
}

func (h *Handler) add(w http.ResponseWriter, r *http.Request) {

	var (
		post models.Post
	)

	if decodeErr := json.NewDecoder(r.Body).Decode(&post); decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(decodeErr.Error()))
		return
	}

	post.Prepare()
	if valRes := post.Validate(); len(valRes) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(valRes))
		return
	}

	res, err := h.service.add(r.Context(), &post)

	if err != nil {
		fmt.Println("add post error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	utils.WriteJson(res, w)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {

	var (
		post models.Post
	)

	if decodeErr := json.NewDecoder(r.Body).Decode(&post); decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(decodeErr.Error()))
		return
	}

	post.Prepare()
	if valRes := post.Validate(); len(valRes) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(valRes))
		return
	}

	res, err := h.service.update(r.Context(), &post)

	if err != nil {
		fmt.Println("update post error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	utils.WriteJson(res, w)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {

	var (
		params = mux.Vars(r)
		postId = params["postId"]
	)

	if err := h.service.delete(r.Context(), postId); err != nil {
		fmt.Println("delete post error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	w.Write(utils.OK)
}
