package model

import "strconv"

type Post struct {
	ID     int    `db:"id"`
	UserID int    `db:"user_id"`
	Title  string `db:"title"`
	User   *User  `db:"-" orm:"relation:BelongsTo;mainField:UserID;assocField:ID"`
}

func (m *Post) TableName() string {
	return "posts"
}

func (m *Post) GetPK() string {
	return strconv.Itoa(m.ID)
}
