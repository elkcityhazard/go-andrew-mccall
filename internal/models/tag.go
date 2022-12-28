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

func (t *Tag) GetPostsByTags(db *sql.DB, id int) ([]*Post, error) {
	stmt := `SELECT posts.id, title, summary, expires_at FROM posts, tags WHERE posts.expires_at > UTC_TIMESTAMP() AND tags.id = ? ORDER BY title DESC LIMIT 5`

	rows, err := db.Query(stmt, id)

	if err != nil {
		return nil, err
	}

	posts := []*Post{}

	for rows.Next() {

		p := &Post{}

		err = rows.Scan(&p.Id, &p.Title, &p.Summary, &p.ExpireDate)

		if err != nil {
			return nil, err
		}

		posts = append(posts, p)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil

}
