package dto

type User struct {
	ID   int
	Name string
}

func (u *User) TableName() string {
	return "user"
}
