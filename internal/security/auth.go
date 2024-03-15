package security

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

var SecretKey = []byte("F4BD4860E0E7488DFEBB469A9C8054D08586E7F2638658AEA0CAD7B08A55C727")

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("token")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	}
}

func ExtractTokenInfo(r *http.Request) (jwt.MapClaims, error) {
	tokenString := r.Header.Get("token")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to extract claims")
	}

	return claims, nil
}
