package handlers

import (
	"encoding/json"
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
	Tools     utils.Tools
}

var app *models.AppConfig

var Repo *Repository

func NewHandlers(a *models.AppConfig) {
	app = a
}

func NewRepo(app *models.AppConfig) *Repository {
	return &Repository{
		AppConfig: app,
		Tools:     utils.Tools{},
	}
}

func SetRepo(m *Repository) {
	Repo = m

}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	var post = models.Post{}

	posts, err := post.GetPostWithLimitAndOffset(m.AppConfig.DB, 3, 0)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data = make(map[string]interface{})

	title := "Andrew M McCall - <br />Traverse City Web Design"

	data["MostRecent"] = posts

	data["SiteTitle"] = template.HTML(fmt.Sprintf("%s", title))

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

func (m *Repository) UpdateAvatar(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	switch r.Method {

	case "GET":
		// Do something

		exists := app.SessionManager.Exists(r.Context(), "authenticatedUserID")

		if !exists {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		id := app.SessionManager.GetString(r.Context(), "authenticatedUserID")

		if id != params.ByName("id") {
			http.Error(w, "authentication error", http.StatusUnauthorized)
			return
		}

		stringMap := map[string]string{}

		stringMap["ID"] = id

		// render the template
		render.RenderTemplate(w, r, "update-avatar.tmpl.html", &models.TemplateData{
			StringMap: stringMap,
		})

	case "POST":

		id := app.SessionManager.GetString(r.Context(), "authenticatedUserID")

		user := models.User{}

		toInt, err := strconv.Atoi(id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err = user.GetUserById(app.DB, toInt)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(user)

		// Do Something
		m.Tools.MaxFileSize = 1024 * 1024 * 1024
		m.Tools.AllowedFileTypes = []string{"image/png", "image/jpeg", "image/jpg"}
		file, err := m.Tools.UploadSingleFile(r, user.Email, false)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		app.SessionManager.Put(r.Context(), "flash", fmt.Sprintf("You uploaded %s to %s", file.NewFileName, user.PathToAvatar))

		data := make(map[string]interface{})

		data["file"] = file

		result, err := user.UpdateUserAvatar(app.DB, toInt, file.NewFileName)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pathToUse := path.Join("./static/uploads/", file.NewFileName)

		err = m.Tools.ResizeImage(file, pathToUse, &user, w, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = user.UpdateUserAvatar(app.DB, user.Id, pathToUse)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data["result"] = result

		render.RenderTemplate(w, r, "update-avatar.tmpl.html", &models.TemplateData{
			Data: data,
		})

	default:
		http.Error(w, "nothing to see here", http.StatusNotFound)
		return
	}
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	app.SessionManager.Remove(r.Context(), "authenticatedUserID")

	app.SessionManager.Put(r.Context(), "flash", "You have been logged out successfully")

	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	return
}

func (m *Repository) BulkUpload(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		// do something

		var stringMap map[string]string

		render.RenderTemplate(w, r, "bulk-upload.tmpl.html", &models.TemplateData{
			StringMap: stringMap,
		})
	case "POST":
		// do something
	default:
		http.Error(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}
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
		uploadTools.AllowedFileTypes = []string{"image/jpeg", "image/jpg", "image/png", "audio/mpeg"}

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

		var cat models.Category

		cat.Name = r.Form.Get("category")

		catRows, err := cat.CheckIfCategoryExistsAndReturn(app.DB, cat.Name)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(catRows)

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

//	Logging handles displaying the login page and posting the login

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

			// Category

			var catSlice []*models.Category

			c := models.Category{}

			catSlice, err = c.GetCategoryByPostId(m.AppConfig.DB, v.Id)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var catName string

			if len(catSlice) > 0 {
				catName = catSlice[0].Name
			} else {
				catName = "null"
			}

			v.Categories = []string{}

			v.Categories = append(v.Categories, catName)

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

		// Category

		var catSlice []*models.Category

		c := models.Category{}

		catSlice, err = c.GetCategoryByPostId(m.AppConfig.DB, 24)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("catSlice: ", catSlice[0])

		render.RenderTemplate(w, r, "single-post.tmpl.html", &models.TemplateData{
			Data:       data,
			Categories: catSlice,
		})
	}

}

func (m *Repository) GetSinglePost(w http.ResponseWriter, r *http.Request) {
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

	// Category

	var catSlice []*models.Category

	c := models.Category{}

	catSlice, err = c.GetCategoryByPostId(m.AppConfig.DB, postID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("catSlice: ", catSlice[0])

	//	Tags

	var t *models.Tag

	tags, err := t.GetTagById(m.AppConfig.DB, postID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Related By Category

	var catPosts []*models.Post

	catPosts, err = c.GetPostsByCategoryId(m.AppConfig.DB, catSlice[0].Slug)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["relatedByCategory"] = catPosts

	//	Tags

	var relatedByTags []*models.Post

	for _, v := range tags {
		posts, err := v.GetPostsByTags(m.AppConfig.DB, v.Id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, w := range posts {
			relatedByTags = append(relatedByTags, w)
		}

	}

	reducedTags := map[string]*models.Post{}

	for _, v := range relatedByTags {
		if reducedTags[v.Title] == nil {
			reducedTags[v.Title] = v
		}
	}

	data["relatedByTags"] = reducedTags

	render.RenderTemplate(w, r, "single-post.tmpl.html", &models.TemplateData{
		Data:       data,
		Categories: catSlice,
		Tags:       tags,
	})
}

func (m *Repository) GetCategories(w http.ResponseWriter, r *http.Request) {

	// getParams

	qp := r.URL.Query()

	// set up struct fields for json response

	type Categories struct {
		Categories []models.Category `json:"categories"`
	}

	// create a cat var

	var cat *models.Category

	// return category rows from database

	cats, err := cat.CheckIfCategoryExistsAndReturn(m.AppConfig.DB, qp.Get("category"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer cats.Close()

	// creat a slice to hold categories

	var catArray []*models.Category

	//	loop through categories and add them to the payload

	for cats.Next() {
		c := &models.Category{}

		err = cats.Scan(&c.Id, &c.Name, &c.Slug, &c.PostId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		catArray = append(catArray, c)
	}

	if err = cats.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set the content type header to send JSON correctly

	w.Header().Set("Content-Type", "application/json")

	// work some JSON magic

	formattedPayload := map[string][]*models.Category{}

	formattedPayload["categories"] = catArray

	catPayload, err := json.Marshal(formattedPayload)

	if err != nil {
		type Error struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
		}

		errMsg := Error{
			Status:  500,
			Message: err.Error(),
		}

		payload, err := json.Marshal(errMsg)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(payload)
	}

	w.Write(catPayload)
}
