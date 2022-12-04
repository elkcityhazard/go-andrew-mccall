package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/elkcityhazard/go-andrew-mccall/internal/models"
)

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var myFuncMap = template.FuncMap{
	"humanDate": humanDate,
}

var app *models.AppConfig

func NewRenderer(a *models.AppConfig) {
	app = a
}

func AddDefaultTemplateData(r *http.Request) models.DefaultTemplateData {

	td := models.DefaultTemplateData{
		Navigation: []models.Navigation{
			{
				Name:              "Andrew M McCall",
				URL:               "/",
				Weight:            1,
				HasAuthentication: true,
			},
			{
				Name:              "Andrew M McCall",
				URL:               "/",
				Weight:            1,
				HasAuthentication: false,
			},
			{
				Name:              "About",
				URL:               "/about",
				Weight:            2,
				HasAuthentication: false,
			},
			{
				Name:              "Blog",
				URL:               "/posts",
				Weight:            3,
				HasAuthentication: false,
			},
			{
				Name:              "Blog",
				URL:               "/posts",
				Weight:            3,
				HasAuthentication: true,
			},
			{
				Name:              "Login",
				URL:               "/admin/login",
				Weight:            4,
				HasAuthentication: false,
			},
			{
				Name:              "Logout",
				URL:               "/admin/logout",
				Weight:            5,
				HasAuthentication: true,
			},
			{
				Name:              "New Post",
				URL:               "/admin/add-post",
				Weight:            6,
				HasAuthentication: true,
			},
			{
				Name:              "Bulk Upload",
				URL:               "/admin/bulk-upload",
				Weight:            7,
				HasAuthentication: true,
			},
		},
		IsAuthenticated: app.IsAuthenticated(r),
	}

	td.SiteTitle = "Andrew McCall - Traverse City Web Design"
	td.SocialMedia = td.AddSocial()
	td.Nap = td.AddNap()
	td.CSRFToken = td.AddNoSurf(r)
	td.FlashMessage = td.AddFlashMessage(app, r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	if app.IsProduction {
		tc = app.TemplateCache
	} else {
		newTC, err := CreateTemplateCache()

		if err != nil {
			return err
		}

		tc = newTC
	}

	t, ok := tc[tmpl]

	if !ok {
		err := errors.New("error parsing the template set")
		return err
	}

	td.DefaultTemplateData = AddDefaultTemplateData(r)

	// need buf to write to first to ensure everything goes okay

	buf := new(bytes.Buffer)

	err := t.Execute(buf, td)

	if err != nil {
		fmt.Println(err)
		return err
	}

	// now we can write the buffer to the response writer if everything goes okay

	_, err = buf.WriteTo(w)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {

	// Create a Map As Cache

	myTemplateCache := map[string]*template.Template{}

	// go out and fetch all of the pages

	pages, err := filepath.Glob("./templates/pages/*.tmpl.html")

	if err != nil {
		return nil, err
	}

	// range through the pages

	for _, page := range pages {

		// get the file name

		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(myFuncMap).ParseFiles(page)

		if err != nil {
			return nil, err
		}

		// check for associated layouts to page

		layoutMatches, err := filepath.Glob("./templates/layouts/*.tmpl.html")

		if err != nil {
			return nil, err
		}

		if len(layoutMatches) > 0 {
			ts, err := ts.ParseGlob("./templates/layouts/*.tmpl.html")
			if err != nil {
				return nil, err
			}

			myTemplateCache[name] = ts
		}

		// check for associated layouts to page

		partialsMatches, err := filepath.Glob("./templates/partials/*.tmpl.html")

		if err != nil {
			return nil, err
		}

		if len(partialsMatches) > 0 {
			ts, err := ts.ParseGlob("./templates/partials/*.tmpl.html")
			if err != nil {
				return nil, err
			}

			myTemplateCache[name] = ts
		}

	}

	// if everything checks out, add an entry to the templateSet Cache map

	// return the cache and error
	return myTemplateCache, nil
}
