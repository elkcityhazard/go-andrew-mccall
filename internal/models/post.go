package models

import (
	"bytes"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
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

func (p *Post) InsertIntoDB(db *sql.DB, id string) (sql.Result, error) {

	newId, err := strconv.Atoi(id)

	if err != nil {
		return nil, err
	}

	p.UserID = newId

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(p.Content), &buf); err != nil {
		panic(err)
	}

	var stmt = `INSERT INTO posts (
				   title,
                   description,
                   summary,
                   created_at, 
                   updated_at, 
                   expires_at, 
                   featured_image, 
                   content, 
                   author_id
				)
				   VALUES(?,?,?,?,?,?,?,?,?)`
	res, err := db.Exec(stmt, p.Title, p.Description, p.Summary, p.PublishDate, time.Now(), p.ExpireDate, p.FeaturedImage, buf.String(), p.UserID)

	if err != nil {
		return nil, err
	}

	return res, nil

}

func (p *Post) GetSinglePost(db *sql.DB, id int) (*Post, error) {
	stmt := `SELECT id, title, content, summary, description, author_id, created_at, updated_at, expires_at, featured_image FROM posts WHERE expires_at > UTC_TIMESTAMP() and id = ?`

	row := db.QueryRow(stmt, id)

	cp := &Post{}

	err := row.Scan(&cp.Id, &cp.Title, &cp.Content, &cp.Summary, &cp.Description, &cp.UserID, &cp.PublishDate, &cp.UpdatedDate, &cp.ExpireDate, &cp.FeaturedImage)

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
	stmt := `SELECT id, title, content, summary, author_id, created_at, updated_at, expires_at, featured_image FROM posts WHERE expires_at > UTC_TIMESTAMP() ORDER BY created_at ASC limit 10`

	rows, err := db.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		cp := &Post{}

		err := rows.Scan(&cp.Id, &cp.Title, &cp.Content, &cp.Summary, &cp.AuthorId, &cp.PublishDate, &cp.UpdatedDate, &cp.ExpireDate, &cp.FeaturedImage)

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

func (p *Post) GetPostWithLimitAndOffset(db *sql.DB, limit int, offset int) ([]*Post, error) {
	stmt := `SELECT id, title, content, summary, author_id, created_at, updated_at, expires_at, featured_image FROM posts WHERE expires_at > UTC_TIMESTAMP() ORDER BY created_at ASC limit ? offset ?`

	rows, err := db.Query(stmt, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		cp := &Post{}

		err := rows.Scan(&cp.Id, &cp.Title, &cp.Content, &cp.Summary, &cp.AuthorId, &cp.PublishDate, &cp.UpdatedDate, &cp.ExpireDate, &cp.FeaturedImage)

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
