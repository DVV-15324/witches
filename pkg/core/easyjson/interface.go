package easyjson

// JSONMarshaler - Interface cho struct có thể marshal JSON
type JSONMarshaler interface {
	MarshalJSON() ([]byte, error)
}

// JSONUnmarshaler - Interface cho struct có thể unmarshal JSON
type JSONUnmarshaler interface {
	UnmarshalJSON([]byte) error
}

// JSONer - Tổng hợp cả 2
type JSONer interface {
	JSONMarshaler
	JSONUnmarshaler
}
