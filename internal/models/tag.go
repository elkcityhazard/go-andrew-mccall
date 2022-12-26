package models

import "database/sql"

type Tag struct {
	Id     int
	Name   string
	Slug   string
	PostId int
}

func (t *Tag) GetTagById(db *sql.DB, id int) ([]*Tag, error) {

	stmt := `SELECT * FROM tags where post_id = ?`

	rows, err := db.Query(stmt, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []*Tag

	for rows.Next() {
		t := &Tag{}

		err = rows.Scan(&t.Id, &t.Name, &t.Slug, &t.PostId)

		if err != nil {
			return nil, err
		}

		tags = append(tags, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil

}
