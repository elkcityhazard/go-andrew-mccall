package models

import (
	"github.com/justinas/nosurf"
	"net/http"
)

type DefaultTemplateData struct {
	SiteTitle       string
	Navigation      []Navigation
	SocialMedia     []SocialLink
	Nap             NAP
	IsAuthenticated bool
	CSRFToken       string
	FlashMessage    string
}

type TemplateData struct {
	StringMap           map[string]string
	IntMap              map[string]string
	Data                map[string]interface{}
	Categories          []*Category
	Tags                []*Tag
	DefaultTemplateData DefaultTemplateData
}

type SocialLink struct {
	Name     string
	Username string
	Url      string
	FilePath string
}

type NAP struct {
	Name      string
	Phone     string
	PhoneHref string
	Street    string
	City      string
	State     string
	ZipCode   int
}

func (t *DefaultTemplateData) AddNoSurf(r *http.Request) string {
	return nosurf.Token(r)
}

func (t *DefaultTemplateData) AddNap() NAP {
	return NAP{
		Name:      "Andrew McCall - Traverse City Web Design",
		Phone:     "(231) 299-0217",
		PhoneHref: "tel:+12312990217",
		Street:    "4889 Silver Pines Rd",
		City:      "Traverse City",
		State:     "Michigan",
		ZipCode:   49685,
	}
}

func (t *DefaultTemplateData) AddSocial() []SocialLink {
	return []SocialLink{
		{
			Name:     "twitter",
			Username: "elkcityhazard",
			Url:      "https://twitter.com/elkcityhazard",
			FilePath: "/static/images/social-media-icons/twitter.png",
		},
		{
			Name:     "mastodon",
			Username: "elkcityhazard",
			Url:      "https://indieweb.social/@elkcityhazard",
			FilePath: "/static/images/social-media-icons/mastodon.png",
		},
		{
			Name:     "youtube",
			Username: "elkcityhazard",
			Url:      "https://www.youtube.com/user/elkcityhazard/featured",
			FilePath: "/static/images/social-media-icons/youtube.png",
		},
		{
			Name:     "instagram",
			Username: "elkcityhazard",
			Url:      "https://instagram.com/elkcityhazard",
			FilePath: "/static/images/social-media-icons/instagram.png",
		},
	}
}

func (t *DefaultTemplateData) AddFlashMessage(app *AppConfig, r *http.Request) string {
	msg := app.SessionManager.PopString(r.Context(), "flash")
	return msg
}
