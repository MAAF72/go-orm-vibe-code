package model

type User struct {
	ID           int    `db:"id"`
	Name         string `db:"name"`
	Posts        []Post `db:"-" orm:"relation:HaveMany;mainField:ID;assocField:UserID"`
	SupervisorID *int   `db:"supervisor_id"`
	Supervisor   *User  `db:"-" orm:"relation:BelongsTo;mainField:SupervisorID;assocField:ID"`
}

func (m *User) TableName() string {
	return "users"
}
