package easyjson

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

// UnmarshalJSON - Tự động chọn easyjson hoặc encoding/json
func UnmarshalJSON(data []byte, v interface{}) error {
	// Nếu v implement JSONUnmarshaler -> dùng easyjson
	if u, ok := v.(JSONUnmarshaler); ok {
		return u.UnmarshalJSON(data)
	}
	// Fallback -> encoding/json
	return json.Unmarshal(data, v)
}

// MarshalJSON - Tự động chọn easyjson hoặc encoding/json
func MarshalJSON(v interface{}) ([]byte, error) {
	if m, ok := v.(JSONMarshaler); ok {
		return m.MarshalJSON()
	}
	return json.Marshal(v)
}

// BindJSON - Bind JSON từ Gin context
func BindJSON(c *gin.Context, v interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	return UnmarshalJSON(body, v)
}
