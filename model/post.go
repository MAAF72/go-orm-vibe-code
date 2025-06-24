package model

type Post struct {
	ID     int    `db:"id"`
	Title  string `db:"title"`
	UserID int    `db:"user_id"`
	User   *User  `db:"-" orm:"relation:BelongsTo;mainField:UserID;assocField:ID"`
}

func (m *Post) TableName() string {
	return "posts"
}
