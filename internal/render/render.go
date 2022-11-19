package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/elkcityhazard/go-andrew-mccall/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
)

var myFuncMap template.FuncMap

var app *models.AppConfig

func NewRenderer(a *models.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	fmt.Println(app)

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

	// need buf to write to first to ensure everything goes okay

	buf := new(bytes.Buffer)

	err := t.Execute(buf, td)

	if err != nil {
		return err
	}

	// now we can write the buffer to the response writer if everything goes okay

	_, err = buf.WriteTo(w)

	if err != nil {
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
