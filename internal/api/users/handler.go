package users

import (
	"checkmarks/internal/api/common/access"
	"checkmarks/internal/utils"
	"checkmarks/pkg/users"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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
	router.HandleFunc("/v1/users", h.add).Methods("POST")
	router.HandleFunc("/v1/login", h.login).Methods("POST")
	router.HandleFunc("/v1/loginAuth", h.loginAuth).Methods("POST")

}

func (h *Handler) add(w http.ResponseWriter, r *http.Request) {
	var (
		user users.User
	)

	if decodeErr := json.NewDecoder(r.Body).Decode(&user); decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(decodeErr.Error()))
		return
	}

	user.Prepare()

	if valRes := user.Validate(); len(valRes) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(valRes))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("hashed password error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	user.Password = string(hashedPassword)

	err = h.service.add(r.Context(), &user)

	if err != nil {
		fmt.Println("add post error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	w.Write(utils.OK)

}
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) loginAuth(w http.ResponseWriter, r *http.Request) {

}
