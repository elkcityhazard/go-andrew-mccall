package models

type Navigation struct {
	Name     string
	URL      string
	Weight   int
	Children []Navigation
}
