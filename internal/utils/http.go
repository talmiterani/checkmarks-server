package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJson(payload interface{}, w http.ResponseWriter) {
	//"application/json; charset=UTF-8"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
