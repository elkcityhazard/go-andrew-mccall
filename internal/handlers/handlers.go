package handlers

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"net/mail"
	"path"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

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

	var data = make(map[string]interface{})

	data["SiteTitle"] = "Andrew M McCall - Traverse City Web Design"

	render.RenderTemplate(w, r, "home.tmpl.html", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		var stringMap = map[string]string{}

		stringMap["title"] = "Login"

		render.RenderTemplate(w, r, "login.tmpl.html", &models.TemplateData{
			StringMap: stringMap,
		})
	case "POST":
		err := r.ParseForm()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		email := r.Form.Get("email")
		password := r.Form.Get("password")

		if email == "" || password == "" {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		_, err = mail.ParseAddress(email)

		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		u := models.User{}

		user, err := u.GetUserByEmail(m.AppConfig.DB, email)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		fmt.Println(user.Id)

		app.SessionManager.Put(r.Context(), "authenticatedUserID", strconv.Itoa(user.Id))
		app.SessionManager.Put(r.Context(), "flash", "Successfully Logged in!")

		http.Redirect(w, r, "/posts/", http.StatusSeeOther)

	}

}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	app.SessionManager.Remove(r.Context(), "authenticatedUserID")

	app.SessionManager.Put(r.Context(), "flash", "You have been logged out successfully")

	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	return
}

func (m *Repository) AddPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "POST, GET")

	fmt.Println(app.SessionManager.Exists(r.Context(), "authenticatedUserID"))
	switch r.Method {
	case "GET":

		var data = make(map[string]interface{})

		html := `<script>
		console.log("does this work?")
				var simplemde = new SimpleMDE({ element: document.getElementById("postContent") });
				simplemde.value("Write Markdown Baby");
				</script>`

		data["html"] = template.HTML(fmt.Sprintf("%s", html))

		render.RenderTemplate(w, r, "create-post.tmpl.html", &models.TemplateData{
			Data: data,
		})
	case "POST":
		var uploadTools utils.Tools

		uploadTools.MaxFileSize = 2 << 20
		uploadTools.AllowedFileTypes = []string{"image/jpeg", "image/jpg", "image/png"}

		loggedIn := app.SessionManager.Exists(r.Context(), "authenticatedUserID")

		if !loggedIn {
			err := errors.New("authentication error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userId := app.SessionManager.GetString(r.Context(), "authenticatedUserID")

		intId, err := strconv.Atoi(userId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		currentUser, err := utils.GetAuthor(app.DB, intId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pathToImage, err := uploadTools.UploadSingleFile(r, currentUser.Email, false)

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p := models.Post{}

		p.Title = r.Form.Get("title")
		p.Content = r.Form.Get("content")
		p.Description = r.Form.Get("description")
		p.Summary = r.Form.Get("summary")
		p.FeaturedImage = path.Join("/static/uploads", pathToImage.NewFileName)

		fmt.Println(p.FeaturedImage)

		fmt.Println(r.Form.Get("publishDate"))

		pd, err := time.Parse("2006-01-02", r.Form.Get("publishDate"))
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}
		p.PublishDate = pd
		ed, err := time.Parse("2006-01-02", r.Form.Get("expireDate"))
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}
		p.ExpireDate = ed

		ud, err := time.Parse("2006-01-02", r.Form.Get("expireDate"))
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}
		p.UpdatedDate = ud

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = p.InsertIntoDB(m.AppConfig.DB, userId)

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return

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

		if err != nil {
			fmt.Println(err)
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		email := html.EscapeString(r.Form.Get("email"))
		password := html.EscapeString(r.Form.Get("password"))

		var fileUtil utils.Tools

		fileUtil.AllowedFileTypes = []string{"image/png", "image/jpeg"}

		file, err := fileUtil.UploadSingleFile(r, email, false)

		if err != nil {
			fmt.Println("Error: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(file)

		encrpytedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		u := models.User{}

		u.Email = email
		u.Password = encrpytedPassword
		u.PathToAvatar = path.Join(fmt.Sprintf("%s%s", "./static/uploads", file.NewFileName))

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

		http.Redirect(w, r, "/admin/add-post", http.StatusSeeOther)
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

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v := idToken.Value

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

func (m *Repository) GetListOfPosts(w http.ResponseWriter, r *http.Request) {

	fmt.Println(app.SessionManager.GetString(r.Context(), "authenticatedUserID"))

	postKey := r.URL.Path[len("/posts/"):]

	switch postKey {
	case "":
		post := models.Post{}

		posts, err := post.GetMultiplePosts(m.AppConfig.DB)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var authors []*models.Author
		var avatars []template.HTML

		for _, v := range posts {
			a, err := utils.GetAuthor(m.AppConfig.DB, v.AuthorId)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			authors = append(authors, a)
			avatar := template.HTML(fmt.Sprintf("%s", string(a.PathToAvatar)))
			avatars = append(avatars, avatar)
		}

		data := make(map[string]interface{})

		data["posts"] = posts
		data["authors"] = authors
		data["avatars"] = avatars

		err = render.RenderTemplate(w, r, "list-posts.tmpl.html", &models.TemplateData{
			Data: data,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case postKey:
		p := &models.Post{}

		idKey, err := strconv.Atoi(postKey)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		post, err := p.GetSinglePost(m.AppConfig.DB, idKey)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("User ID: ", post.UserID)

		author, err := utils.GetAuthor(m.AppConfig.DB, post.UserID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := make(map[string]interface{})

		data["post"] = post

		data["author"] = author

		render.RenderTemplate(w, r, "single-post.tmpl.html", &models.TemplateData{
			Data: data,
		})
	}

}

func (m Repository) GetSinglePost(w http.ResponseWriter, r *http.Request) {
	p := &models.Post{}

	params := httprouter.ParamsFromContext(r.Context())

	postID, err := strconv.Atoi(params.ByName("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post, err := p.GetSinglePost(m.AppConfig.DB, postID)

	if err != nil {
		app.SessionManager.Put(r.Context(), "flash", "The page could not be found")
		render.RenderTemplate(w, r, "404.tmpl.html", &models.TemplateData{})
		return
	}

	fmt.Println("User ID: ", post.UserID)

	author, err := utils.GetAuthor(m.AppConfig.DB, post.UserID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})

	data["author"] = author

	data["description"] = post.Description

	data["post"] = post

	data["avatar"] = author.PathToAvatar

	html := template.HTML(fmt.Sprintf("%s", post.Content))

	data["html"] = html

	render.RenderTemplate(w, r, "single-post.tmpl.html", &models.TemplateData{
		Data: data,
	})
}
