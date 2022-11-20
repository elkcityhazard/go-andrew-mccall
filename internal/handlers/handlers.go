package handlers

import (
	"fmt"
	"html"
	"net/http"
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

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		render.RenderTemplate(w, r, "login.tmpl.html", &models.TemplateData{})

	case "POST":

	default:
	}
}

// CreateUser Creates A New User

func (m *Repository) AddUser(w http.ResponseWriter, r *http.Request) {

	// stmt := `INSERT INTO users (email, password) VALUES(?,?);`

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

	u.InsertIntoDB(m.AppConfig.DB)

}

func (m *Repository) GetJWT(w http.ResponseWriter, r *http.Request) {
	token, err := utils.CreateToken()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
