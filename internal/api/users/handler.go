package users

import (
	"checkmarks/internal/api/common/access"
	"checkmarks/internal/security"
	"checkmarks/internal/utils"
	"checkmarks/pkg/users"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
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
	router.HandleFunc("/v1/users/signup", h.signup).Methods("POST")
	router.HandleFunc("/v1/users/login", h.login).Methods("POST")
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
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
	isUnique, err := h.service.checkUniqueUsername(r.Context(), user.Username)
	if err != nil {
		fmt.Println("got error when check if username is unique, error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	if !isUnique {
		w.WriteHeader(http.StatusConflict)
		utils.WriteJson("this username already exists", w)
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
	err = h.service.signup(r.Context(), &user)

	if err != nil {
		fmt.Println("add post error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
		return
	}

	tokenString, err := createToken(user.Username, user.Id)
	if err != nil {
		fmt.Println("create token error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
	}

	utils.WriteJson(map[string]string{"token": tokenString, "userId": user.Id.Hex()}, w)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {

	var (
		req users.User
	)

	if decodeErr := json.NewDecoder(r.Body).Decode(&req); decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(decodeErr.Error()))
		return
	}

	req.Prepare()

	if valRes := req.Validate(); len(valRes) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(valRes))
		return
	}

	user, err := h.service.get(r.Context(), &req)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			fmt.Println("invalid username, error: ", err)
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			fmt.Println("get user error, error: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		utils.WriteJson(err, w)
		return
	}

	// Compare hashed password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		fmt.Println("invalid password, error: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		utils.WriteJson(err, w)
		return
	}

	tokenString, err := createToken(user.Username, user.Id)
	if err != nil {
		fmt.Println("create token error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJson(err, w)
	}

	utils.WriteJson(map[string]string{"token": tokenString, "userId": user.Id.Hex()}, w)
}

func createToken(username string, id *primitive.ObjectID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"id":       id.Hex(),
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(security.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
