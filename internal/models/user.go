package models

import (
	"database/sql"
	"fmt"
	"path"
)

type User struct {
	Id           int
	Email        string
	Password     []byte
	PathToAvatar string
}

func (u *User) InsertIntoDB(db *sql.DB) (sql.Result, error) {

	stmt := `INSERT INTO users (email, path_to_avatar, password) VALUES(?,?,?)`

	res, err := db.Exec(stmt, u.Email, u.PathToAvatar, u.Password)

	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, err

	}

	return res, nil

}

func (u *User) GetUserByEmail(db *sql.DB, email string) (User, error) {
	stmt := `SELECT * FROM users WHERE email = ?`

	newU := User{}

	row := db.QueryRow(stmt, email)

	err := row.Scan(&newU.Id, &newU.Email, &newU.Password, &newU.PathToAvatar)

	if err != nil {
		return newU, err
	}

	return newU, nil
}

func (u *User) GetUserById(db *sql.DB, id int) (User, error) {
	stmt := `SELECT * FROM users WHERE id = ?`

	newU := User{}

	row := db.QueryRow(stmt, id)

	err := row.Scan(&newU.Id, &newU.Email, &newU.Password, &newU.PathToAvatar)

	if err != nil {
		return newU, err
	}

	return newU, nil
}

func (u *User) UpdateUserAvatar(db *sql.DB, id int, file string) (sql.Result, error) {

	stmt := `UPDATE users SET path_to_avatar = ? WHERE id = ?`

	res, err := db.Exec(stmt, path.Join(fmt.Sprintf("%s", file)), id)

	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, err

	}

	return res, nil
}
