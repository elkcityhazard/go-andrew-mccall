package models

import (
	"database/sql"
	"errors"
)

var app *AppConfig

type Category struct {
	Id    int
	Name  string
	Slug  string
	Posts []Post
}

func (c *Category) CheckIfCategoryExistsAndReturn(db *sql.DB, identifier ...string) (*sql.Rows, error) {
	stmt := `SELECT * FROM categories WHERE name = ? OR slug = ?`

	cat, err := app.DB.Query(stmt, identifier)

	if err != nil {
		return nil, err
	}

	return cat, nil
}

func (c *Category) CreateNewCategory(db *sql.DB, name, slug string, postId int) (sql.Result, error) {

	rows, err := c.CheckIfCategoryExistsAndReturn(db, name)

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return nil, errors.New("category already exists")
	}

	stmt := `INSERT INTO categories (name, slug, post_id) VALUES(?, ?, ?)`

	result, err := db.Exec(stmt, name, slug, postId)

	if err != nil {
		return nil, err
	}

	return result, nil
}
