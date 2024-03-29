package posts

import (
	"checkmarks/internal/api/common/access"
	commonModels "checkmarks/internal/api/common/models"
	"checkmarks/internal/api/posts/models"
	"checkmarks/internal/security"
	"checkmarks/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	router.HandleFunc("/v1/posts/search", h.search).Methods("POST")
	router.HandleFunc("/v1/posts/post/{postId}", h.get).Methods("GET")
	router.HandleFunc("/v1/posts", security.RequireAuth(h.add)).Methods("POST")
	router.HandleFunc("/v1/posts", security.RequireAuth(h.update)).Methods("PUT")
	router.HandleFunc("/v1/posts/post/{postId}", security.RequireAuth(h.delete)).Methods("DELETE")

}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {

	var (
		req commonModels.SearchReq
	)

	if decodeErr := json.NewDecoder(r.Body).Decode(&req); decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(decodeErr.Error()))
		return
	}

	res, err := h.service.search(r.Context(), &req)

	if err != nil {
		fmt.Println("get posts error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	utils.WriteJson(res, w)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {

	var (
		params = mux.Vars(r)
		postId = params["postId"]
	)

	res, err := h.service.get(r.Context(), postId)

	if err != nil {
		fmt.Printf("get post by id %s error: %s", postId, err)
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
	if valRes := post.Validate(false); len(valRes) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(valRes))
		return
	}

	claims, err := security.ExtractTokenInfo(r)

	if err != nil {
		fmt.Println("fail to parse token error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	idClaim, ok := claims["id"]
	if !ok {
		fmt.Println("id claim not found in token")
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson("id claim not found in token", w)
		return
	}

	idStr, ok := idClaim.(string)
	if !ok {
		fmt.Println("fail to claim id to string, error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		fmt.Println("fail to convert id string to mongo object, error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	post.UserId = &id
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
	if valRes := post.Validate(true); len(valRes) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(valRes))
		return
	}

	res, err := h.service.update(r.Context(), &post)

	if err != nil {
		fmt.Printf("update post by id %s error: %s", post.Id, err)
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
		fmt.Printf("delete post by id %s error: %s", postId, err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	w.Write(utils.OK)
}
