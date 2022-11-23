package utils

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elkcityhazard/go-andrew-mccall/internal/models"
	"github.com/golang-jwt/jwt/v4"
)

var app *models.AppConfig

func NewUtils(a *models.AppConfig) {
	app = a
}

//var secret = []byte(app.JWTSecret)

func CreateToken(id string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["expiration"] = time.Now().Add(time.Minute).Unix()
	claims["id"] = id

	tokenString, err := token.SignedString([]byte(app.JWTSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//	ValidateJWT is middleware

func ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookieToken, err := r.Cookie("Token")

		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		if cookieToken == nil {
			http.Error(w, fmt.Sprintf("%s", errors.New("no token provided")), http.StatusBadRequest)
			return
		}

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)

				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
				}
				return app.JWTSecret, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized :" + err.Error()))
			}
			if token.Valid {
				next(w, r)
			}
		}
	})
}
