package utils

import (
	"net/http"
	"time"

	"github.com/elkcityhazard/go-andrew-mccall/internal/models"
	"github.com/golang-jwt/jwt/v4"
)

var app *models.AppConfig

var secret = []byte(app.JWTSecret)

func NewUtils(a *models.AppConfig) {
	app = a
}

func CreateToken() (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["expiration"] = time.Now().Add(time.Hour).Unix()

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//	ValidateJWT is middleware

func ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["token"] != nil {
			token, err := jwt.Parse(r.Header["token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)

				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
				}
				return secret, nil
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
