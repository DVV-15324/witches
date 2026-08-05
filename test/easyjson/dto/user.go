package dto

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (u *User) TableName() string {
	return "user"
}
