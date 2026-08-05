package dto

type CreateUser struct {
	Name string
}

func (u *CreateUser) GenEasyJson() {}
