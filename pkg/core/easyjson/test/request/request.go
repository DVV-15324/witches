package request

type CreateUser struct {
	Name string
}

func (u *CreateUser) GenEasyJson() {}
