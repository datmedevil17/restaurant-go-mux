package middlewares

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/datmedevil17/restaurant-management/helpers"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("token")
		if clientToken == "" {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "No Authorization header provided"})
			return
		}

		claims, msg := helpers.ValidateToken(clientToken)
		if msg != "" {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": msg})
			return
		}

		r.Header.Set("email", claims.Email)
		r.Header.Set("first_name", claims.First_name)
		r.Header.Set("last_name", claims.Last_name)
		r.Header.Set("uid", claims.Uid)

		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
