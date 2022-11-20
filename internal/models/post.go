package models

import (
	"fmt"
	"time"
)

type Post struct {
	Title         string
	Description   string
	Summary       string
	PublishDate   time.Time
	UpdatedDate   time.Time
	ExpireDate    time.Time
	Tags          []string
	Categories    []string
	FeaturedImage string
	Content       string
}

func (p *Post) InsertIntoDB(app *AppConfig) {

	stmt := `INSERT INTO posts (title, 
                   description, 
                   summary, 
                   publish_date, 
                   update_date, 
                   expire_date, 
                   featured_image, 
                   content, 
                   user_id)
				   VALUES(?,?,?,?,?,?,?,?,?);
			`

	fmt.Println(stmt)

}
