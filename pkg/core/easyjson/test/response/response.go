package response

type User struct {
	ID   string
	Name string
}

func (u *User) GenEasyJson() {}
