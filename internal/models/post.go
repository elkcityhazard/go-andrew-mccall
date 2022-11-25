package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Post struct {
	Id            int
	Title         string
	Description   string
	Summary       string
	AuthorId      int
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

func (p *Post) GetSinglePost(db *sql.DB, id int) (*Post, error) {
	stmt := `SELECT id, title, content, created_at, updated_at, expires_at, featured_image FROM posts WHERE expires < UTC_TIMESTAMP() and id = ?`

	row := db.QueryRow(stmt, id)

	cp := &Post{}

	err := row.Scan(&cp.Id, &cp.Title, &cp.Content, &cp.PublishDate, &cp.UpdatedDate, &cp.ExpireDate, &cp.FeaturedImage)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}

	return cp, nil
}

func (p *Post) GetMultiplePosts(db *sql.DB) ([]*Post, error) {
	stmt := `SELECT id, title, content, author_id, created_at, updated_at, expires_at, featured_image FROM posts WHERE expires_at > UTC_TIMESTAMP()`

	rows, err := db.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		cp := &Post{}

		err := rows.Scan(&cp.Id, &cp.Title, &cp.Content, &cp.AuthorId, &cp.PublishDate, &cp.UpdatedDate, &cp.ExpireDate, &cp.FeaturedImage)

		if err != nil {
			return nil, err
		}

		posts = append(posts, cp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
