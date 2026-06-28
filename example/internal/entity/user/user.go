package user

import (
	"time"
)

type User struct {
	Id        int
	Name      string
	Role      string
	Email     string
	Banned    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) TableName() string {
	return "user"
}
