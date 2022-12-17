package main

import (
	"net/http"

	"github.com/justinas/nosurf"

	"github.com/elkcityhazard/go-andrew-mccall/internal/utils"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You don't
		// need to do this in your own code.
		// w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Security-Policy",
			"default-src 'self' 'unsafe-inline'; style-src 'self' https://maxcdn.bootstrapcdn.com https://cdn.jsdelivr.net https://cdnjs.cloudflare.com fonts.googleapis.com; font-src https://maxcdn.bootstrapcdn.com fonts.gstatic.com; script-src 'unsafe-inline' https://cdn.jsdelivr.net;")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

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
		userId := app.SessionManager.Exists(r.Context(), "authenticatedUserId")

		if !userId {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page and
		// return from the middleware chain so that no subsequent handlers in
		// the chain are executed.
		if !app.IsAuthenticated(r) {
			app.IsLoggedIn = false
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		//	set app.IsLoggedIn to true to validate session

		if app.SessionManager.Exists(r.Context(), "authenticatedUserID") {
			app.IsLoggedIn = true
		}
		// And call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly attributes set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}
