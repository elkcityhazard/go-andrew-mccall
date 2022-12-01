package utils

import (
	"database/sql"
	"errors"

	"github.com/elkcityhazard/go-andrew-mccall/internal/models"
)

func GetAuthor(db *sql.DB, id int) (*models.Author, error) {
	stmt := `SELECT * FROM users where id = ?`

	row := db.QueryRow(stmt, id)

	a := &models.Author{}

	err := row.Scan(&a.Id, &a.Email, &a.Password, &a.PathToAvatar)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		} else {
			return nil, err
		}
	}

	return a, nil
}
