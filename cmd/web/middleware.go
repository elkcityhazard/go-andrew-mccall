package main

import (
	"net/http"

	"github.com/elkcityhazard/go-andrew-mccall/internal/utils"
)

func CheckForAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Api"] != nil {
			if r.Header["Api"][0] == app.APIKey {
				token, err := utils.CreateToken()

				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}

				w.Header().Set("token", token)
			}
		}

		next.ServeHTTP(w, r)
	})
}
