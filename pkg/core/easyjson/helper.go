package easyjson

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

func UnmarshalJSON(data []byte, v interface{}) error {
	if u, ok := v.(JSONUnmarshaler); ok {
		return u.UnmarshalJSON(data)
	}
	return json.Unmarshal(data, v)
}

func MarshalJSON(v interface{}) ([]byte, error) {
	if m, ok := v.(JSONMarshaler); ok {
		return m.MarshalJSON()
	}
	return json.Marshal(v)
}

func BindJSON(c *gin.Context, v interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	return UnmarshalJSON(body, v)
}
