package models

import (
	"database/sql"
	"errors"
	"fmt"
)

var app *AppConfig

type Category struct {
	Id     int
	Name   string
	Slug   string
	PostId int
}

func (c *Category) CheckIfCategoryExistsAndReturn(db *sql.DB, identifier string) (*sql.Rows, error) {
	stmt := `SELECT * FROM categories WHERE name = ?`

	cat, err := db.Query(stmt, identifier)

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

func (c *Category) GetCategoryByPostId(db *sql.DB, id int) ([]*Category, error) {
	stmt := `SELECT * FROM categories WHERE post_id = ?  ORDER BY 'asc'`

	rows, err := db.Query(stmt, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var catSlice []*Category

	for rows.Next() {
		c := &Category{}

		err = rows.Scan(&c.Id, &c.Name, &c.Slug, &c.PostId)

		if err != nil {
			return nil, err
		}

		catSlice = append(catSlice, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	fmt.Sprintf("%v", catSlice)

	return catSlice, nil

}
