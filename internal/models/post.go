package models

import (
	"database/sql"
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
	UserID        int
}

func (p *Post) InsertIntoDB(db *sql.DB) {

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

	fmt.Println(stmt, p.Title, p.Description, p.Summary, p.PublishDate, p.UpdatedDate, p.ExpireDate, p.FeaturedImage, p.Content, p.UserID)

}
