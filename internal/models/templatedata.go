package models

type DefaultTemplateData struct {
	Navigation  []Navigation
	SocialMedia map[string]interface{}
}

type TemplateData struct {
	StringMap           map[string]string
	IntMap              map[string]string
	Data                interface{}
	DefaultTemplateData DefaultTemplateData
}

type SocialLink struct {
	Name     string
	Username string
	Url      string
}
