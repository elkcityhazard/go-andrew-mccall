package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id       int
	Email    string
	Password []byte
}

func (u *User) InsertIntoDB(db *sql.DB) (sql.Result, error) {

	fmt.Println(u.Email, u.Password)

	stmt := `INSERT INTO users (email, password) VALUES(?,?)`

	res, err := db.Exec(stmt, u.Email, u.Password)

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

	err := row.Scan(&newU.Id, &newU.Email, &newU.Password)

	if err != nil {
		return newU, err
	}

	return newU, nil
}
