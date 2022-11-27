package models

type Navigation struct {
	Name              string
	URL               string
	HasAuthentication bool
	Weight            int
	Children          []Navigation
}
