package test

type UserGen struct {
	ID   int
	Name string
}

// Implement GenMarker (method rỗng)
func (UserGen) Gen() {}
