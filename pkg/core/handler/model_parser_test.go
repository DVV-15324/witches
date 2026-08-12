package handle

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelParser_Register_SimpleStruct(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	name := parser.Register(User{})
	assert.Equal(t, "User", name)

	schema, ok := parser.schemas["User"]
	assert.True(t, ok)
	assert.Equal(t, "object", schema.Type)
	assert.Contains(t, schema.Properties, "id")
	assert.Contains(t, schema.Properties, "name")
}

func TestModelParser_Register_Pointer(t *testing.T) {
	parser := NewModelParser()

	type Product struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	name := parser.Register(&Product{})
	assert.Equal(t, "Product", name)

	_, ok := parser.schemas["Product"]
	assert.True(t, ok)
}

func TestModelParser_Register_WithTags(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		ID    int    `json:"id" binding:"required" description:"User ID"`
		Name  string `json:"name" binding:"required" minLength:"3" maxLength:"50"`
		Age   int    `json:"age" minimum:"18" maximum:"99"`
		Email string `json:"email" binding:"required"`
	}

	parser.Register(User{})

	schema := parser.schemas["User"]
	assert.NotEmpty(t, schema.Required)

	// Check name properties
	nameProp := schema.Properties["name"]
	assert.Equal(t, "string", nameProp.Type)
	assert.Equal(t, 3, *nameProp.MinLength)
	assert.Equal(t, 50, *nameProp.MaxLength)

	// Check age properties
	ageProp := schema.Properties["age"]
	assert.Equal(t, "integer", ageProp.Type)
	assert.Equal(t, 18, *ageProp.Minimum)
	assert.Equal(t, 99, *ageProp.Maximum)

	// Check description
	assert.Equal(t, "User ID", schema.Properties["id"].Description)
}

// pkg/core/handler/model_parser_test.go
func TestModelParser_Register_NestedStruct(t *testing.T) {
	parser := NewModelParser()

	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type User struct {
		ID      int     `json:"id"`
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	parser.Register(User{})

	// Check User schema
	userSchema, ok := parser.schemas["User"]
	assert.True(t, ok)
	assert.Contains(t, userSchema.Properties, "address")

	// Check Address schema - có thể không được register riêng
	// Nếu không có, thì address property đã được inline
	if _, ok := parser.schemas["Address"]; !ok {
		// Address không được register riêng, kiểm tra property address
		addrProp := userSchema.Properties["address"]
		assert.NotEmpty(t, addrProp.Type, "Address property should have type")
		t.Log("Address schema is inlined, not registered separately")
	} else {
		t.Log("Address schema registered separately")
	}
}
func TestModelParser_Register_Array(t *testing.T) {
	parser := NewModelParser()

	type Item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// Register slice
	name := parser.Register([]Item{})
	assert.Equal(t, "[]Item", name)

	// Check slice schema
	schema, ok := parser.schemas["[]Item"]
	assert.True(t, ok)
	assert.Equal(t, "array", schema.Type)
	assert.Equal(t, "#/definitions/Item", schema.Items.Ref)
}

func TestModelParser_Register_Primitive(t *testing.T) {
	parser := NewModelParser()

	// String
	name := parser.Register("test")
	assert.Equal(t, "string", name)

	// Int
	name = parser.Register(123)
	assert.Equal(t, "int", name)

	// Float
	name = parser.Register(123.45)
	assert.Equal(t, "float64", name)

	// Bool
	name = parser.Register(true)
	assert.Equal(t, "bool", name)
}

func TestModelParser_GetType(t *testing.T) {
	parser := NewModelParser()

	assert.Equal(t, "string", parser.getType(reflect.TypeOf("")))
	assert.Equal(t, "integer", parser.getType(reflect.TypeOf(1)))
	assert.Equal(t, "number", parser.getType(reflect.TypeOf(1.5)))
	assert.Equal(t, "boolean", parser.getType(reflect.TypeOf(true)))
	assert.Equal(t, "array", parser.getType(reflect.TypeOf([]string{})))
	assert.Equal(t, "object", parser.getType(reflect.TypeOf(struct{}{})))
}

func TestModelParser_GetSchemas(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	parser.Register(User{})

	schemas := parser.GetSchemas()
	assert.Contains(t, schemas, "User")
	assert.Len(t, schemas, 1)
}
