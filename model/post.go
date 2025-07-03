package model

type Post struct {
	ID     int    `db:"id" orm:"key:primary"`
	Title  string `db:"title"`
	UserID int    `db:"user_id" orm:"key:foreign"`
	User   *User  `db:"-" orm:"relation:BelongsTo;mainField:UserID;assocField:ID"`
}

func (m *Post) TableName() string {
	return "posts"
}
