package easyjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Struct không implement interface
type SimpleStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Struct implement JSONUnmarshaler
type CustomUnmarshaler struct {
	Value string
}

func (c *CustomUnmarshaler) UnmarshalJSON(data []byte) error {
	// parse manually
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if val, ok := raw["custom"]; ok {
		c.Value = val
		return nil
	}
	return errors.New("missing custom field")
}

// Struct implement JSONMarshaler
type CustomMarshaler struct {
	Value string
}

func (c CustomMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(`{"custom":"` + c.Value + `"}`), nil
}

//  Test UnmarshalJSON

func TestUnmarshalJSON_WithStandardStruct(t *testing.T) {
	data := []byte(`{"name":"Alice","age":30}`)
	var s SimpleStruct
	err := UnmarshalJSON(data, &s)
	require.NoError(t, err)
	assert.Equal(t, "Alice", s.Name)
	assert.Equal(t, 30, s.Age)
}

func TestUnmarshalJSON_WithInvalidJSON(t *testing.T) {
	data := []byte(`{invalid}`)
	var s SimpleStruct
	err := UnmarshalJSON(data, &s)
	assert.Error(t, err)
}

func TestUnmarshalJSON_WithCustomUnmarshaler(t *testing.T) {
	data := []byte(`{"custom":"hello"}`)
	var c CustomUnmarshaler
	err := UnmarshalJSON(data, &c)
	require.NoError(t, err)
	assert.Equal(t, "hello", c.Value)
}

func TestUnmarshalJSON_WithCustomUnmarshalerError(t *testing.T) {
	data := []byte(`{"wrong":"field"}`)
	var c CustomUnmarshaler
	err := UnmarshalJSON(data, &c)
	assert.Error(t, err)
	assert.Equal(t, "missing custom field", err.Error())
}

func TestUnmarshalJSON_NilInput(t *testing.T) {
	var s SimpleStruct
	err := UnmarshalJSON(nil, &s)
	// json.Unmarshal(nil, &s) sẽ trả về lỗi vì EOF hoặc invalid
	assert.Error(t, err)
}

//  Test MarshalJSON

func TestMarshalJSON_WithStandardStruct(t *testing.T) {
	s := SimpleStruct{Name: "Bob", Age: 25}
	data, err := MarshalJSON(s)
	require.NoError(t, err)
	expected := `{"name":"Bob","age":25}`
	assert.JSONEq(t, expected, string(data))
}

func TestMarshalJSON_WithCustomMarshaler(t *testing.T) {
	m := CustomMarshaler{Value: "test"}
	data, err := MarshalJSON(m)
	require.NoError(t, err)
	expected := `{"custom":"test"}`
	assert.JSONEq(t, expected, string(data))
}

func TestMarshalJSON_ErrorCase(t *testing.T) {
	// Truyền vào channel sẽ gây lỗi
	ch := make(chan int)
	_, err := MarshalJSON(ch)
	assert.Error(t, err)
}

//  Test BindJSON

func setupGinContext(body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, w
}

func TestBindJSON_Success(t *testing.T) {
	body := `{"name":"Charlie","age":40}`
	c, _ := setupGinContext(body)

	var s SimpleStruct
	err := BindJSON(c, &s)
	require.NoError(t, err)
	assert.Equal(t, "Charlie", s.Name)
	assert.Equal(t, 40, s.Age)
}

func TestBindJSON_EmptyBody(t *testing.T) {
	c, _ := setupGinContext("")

	var s SimpleStruct
	err := BindJSON(c, &s)
	// io.ReadAll trả về EOF, Unmarshal sẽ báo lỗi
	assert.Error(t, err)
}

func TestBindJSON_InvalidJSON(t *testing.T) {
	body := `{invalid}`
	c, _ := setupGinContext(body)

	var s SimpleStruct
	err := BindJSON(c, &s)
	assert.Error(t, err)
}

func TestBindJSON_WithCustomUnmarshaler(t *testing.T) {
	body := `{"custom":"world"}`
	c, _ := setupGinContext(body)

	var cm CustomUnmarshaler
	err := BindJSON(c, &cm)
	require.NoError(t, err)
	assert.Equal(t, "world", cm.Value)
}

// Test interface assertions

func TestInterfaces(t *testing.T) {
	// Kiểm tra các interface có được định nghĩa đúng không
	var _ JSONUnmarshaler = (*CustomUnmarshaler)(nil)
	var _ JSONMarshaler = (*CustomMarshaler)(nil)
	// Nếu không có lỗi biên dịch, tức là đúng
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, errors.New("read error")
}

func TestBindJSON_ReadAllError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", errorReader{})
	c.Request.Header.Set("Content-Type", "application/json")

	var s SimpleStruct
	err := BindJSON(c, &s)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read error")
}
