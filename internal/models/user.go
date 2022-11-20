package models

import "database/sql"

type User struct {
	Id       int
	Email    string
	Password []byte
}

func (u *User) InsertIntoDB(db *sql.DB) {
	stmt := `INSERT INTO users (email, password) VALUES(?,?);`

	db.Exec(stmt, u.Email, u.Password)

}
