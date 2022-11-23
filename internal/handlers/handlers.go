package handlers

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/elkcityhazard/go-andrew-mccall/internal/models"
	"github.com/elkcityhazard/go-andrew-mccall/internal/render"
	"github.com/elkcityhazard/go-andrew-mccall/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	AppConfig *models.AppConfig
}

var app *models.AppConfig

var Repo *Repository

func NewHandlers(a *models.AppConfig) {
	app = a
}

func NewRepo(app *models.AppConfig) *Repository {
	return &Repository{
		AppConfig: app,
	}
}

func SetRepo(m *Repository) {
	Repo = m
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	pathKey := r.URL.Path[len("/"):]

	fmt.Println(pathKey)

	fmt.Println("is this working?")

	render.RenderTemplate(w, r, "home.tmpl.html", &models.TemplateData{})
}

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		render.RenderTemplate(w, r, "signup.tmpl.html", &models.TemplateData{})
	case "POST":
		err := r.ParseForm()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		email := r.Form.Get("email")
		password := r.Form.Get("password")

		u := models.User{}

		user, err := u.GetUserByEmail(m.AppConfig.DB, email)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		m.GetJWT(w, r)

	}

}

func (m *Repository) AddPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "POST, GET")
	switch r.Method {
	case "GET":
		render.RenderTemplate(w, r, "create-post.tmpl.html", &models.TemplateData{})
	case "POST":
		err := r.ParseMultipartForm(2 << 20)

		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		p := models.Post{}

		p.Title = r.Form.Get("title")
		p.Description = r.Form.Get("description")
		p.Summary = r.Form.Get("summary")
		pd, err := time.Parse("2006-01-01T00:00:00.000Z", r.Form.Get("publishDate"))
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}
		p.PublishDate = pd
		ed, err := time.Parse("2006-01-01T00:00:00.000Z", r.Form.Get("expireDate"))
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}
		p.ExpireDate = ed

		ud, err := time.Parse("2006-01-01T00:00:00.000Z", r.Form.Get("expireDate"))
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}
		p.UpdatedDate = ud

		p.InsertIntoDB(m.AppConfig.DB)
	}

}

//	Loging handles displaying the login page and posting the login

func (m *Repository) Signup(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		w.Header().Set("Api", app.APIKey)
		render.RenderTemplate(w, r, "signup.tmpl.html", &models.TemplateData{})

	case "POST":
		err := r.ParseForm()

		r.Header.Set("Api", app.APIKey)

		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		email := html.EscapeString(r.Form.Get("email"))
		password := html.EscapeString(r.Form.Get("password"))

		encrpytedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		u := models.User{}

		u.Email = email
		u.Password = encrpytedPassword

		_, err = u.InsertIntoDB(m.AppConfig.DB)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		InsertedUser, err := u.GetUserByEmail(app.DB, u.Email)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		idCookie := &http.Cookie{
			Name:     "Id",
			Value:    strconv.Itoa(InsertedUser.Id),
			HttpOnly: true,
			MaxAge:   60,
		}

		http.SetCookie(w, idCookie)

		http.Redirect(w, r, "/admin/get-jwt", http.StatusSeeOther)
	default:
		fmt.Println("default")
	}
}

// CreateUser Creates A New User

func (m *Repository) AddUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	email := html.EscapeString(r.Form.Get("email"))
	password := html.EscapeString(r.Form.Get("password"))

	encrpytedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u := models.User{}

	u.Email = email
	u.Password = encrpytedPassword

	result, err := u.InsertIntoDB(m.AppConfig.DB)

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("Result: ", result)

	http.Redirect(w, r, "/admin/get-jwt", http.StatusSeeOther)

}

func (m *Repository) GetJWT(w http.ResponseWriter, r *http.Request) {

	idToken, err := r.Cookie("Id")

	v := idToken.Value

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	token, err := utils.CreateToken(v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenCookie := &http.Cookie{
		Name:     "Token",
		Value:    token,
		HttpOnly: true,
		MaxAge:   3600,
	}

	http.SetCookie(w, tokenCookie)

	w.Header().Add("Token", tokenCookie.Value)

	_, err = fmt.Fprint(w, token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
