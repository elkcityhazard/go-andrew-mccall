package main

import (
	"fmt"
	"net/http"

	"github.com/elkcityhazard/go-andrew-mccall/internal/utils"
)

func SetAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.Header.Set("Api", app.APIKey)
		next.ServeHTTP(w, r)
	})
}

func CheckForAPIKey(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Api"] != nil {
			idToken, err := r.Cookie("Id")

			v := idToken.Value

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if r.Header["Api"][0] == app.APIKey {
				token, err := utils.CreateToken(v)

				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}

				w.Header().Set("token", token)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func IsLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := r.Cookie("Token")

		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		fmt.Println(t)
		next.ServeHTTP(w, r)
	})
}
