package easyjson

type JSONMarshaler interface {
	MarshalJSON() ([]byte, error)
}

type JSONUnmarshaler interface {
	UnmarshalJSON([]byte) error
}

type JSONer interface {
	JSONMarshaler
	JSONUnmarshaler
}
