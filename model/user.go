package model

import "strconv"

type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Posts []Post `db:"-" orm:"relation:HaveMany;mainField:ID;assocField:UserID"`
}

func (m *User) TableName() string {
	return "users"
}

func (m *User) GetPK() string {
	return strconv.Itoa(m.ID)
}
