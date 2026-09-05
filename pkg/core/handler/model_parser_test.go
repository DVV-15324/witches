package handle

import (
	"reflect"
	"testing"
	"time"

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

	nameProp := schema.Properties["name"]
	assert.Equal(t, "string", nameProp.Type)
	assert.Equal(t, 3, *nameProp.MinLength)
	assert.Equal(t, 50, *nameProp.MaxLength)

	ageProp := schema.Properties["age"]
	assert.Equal(t, "integer", ageProp.Type)
	assert.Equal(t, 18, *ageProp.Minimum)
	assert.Equal(t, 99, *ageProp.Maximum)

	assert.Equal(t, "User ID", schema.Properties["id"].Description)
}

func TestModelParser_Register_Array(t *testing.T) {
	parser := NewModelParser()

	type Item struct {
		ID int `json:"id"`
	}

	name := parser.Register([2]Item{})

	assert.Equal(t, "[]Item", name)
	assert.Contains(t, parser.schemas, "Item")
	assert.Contains(t, parser.schemas, "[]Item")
	assert.Equal(t, "array", parser.schemas["[]Item"].Type)
}

func TestModelParser_Register_Primitive(t *testing.T) {
	parser := NewModelParser()

	name := parser.Register("test")
	assert.Equal(t, "string", name)

	name = parser.Register(123)
	assert.Equal(t, "int", name)

	name = parser.Register(123.45)
	assert.Equal(t, "float64", name)

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
func TestModelParser_Register_NestedPointer(t *testing.T) {
	parser := NewModelParser()

	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type User struct {
		ID      int      `json:"id"`
		Name    string   `json:"name"`
		Address *Address `json:"address"`
	}

	name := parser.Register(User{})
	assert.Equal(t, "User", name)

	//  Address schema được register
	_, ok := parser.schemas["Address"]
	assert.True(t, ok, "Address schema should be registered")

	//  Kiểm tra User schema
	userSchema, ok := parser.schemas["User"]
	assert.True(t, ok)

	//  Kiểm tra address property có ref đúng
	addrProp, ok := userSchema.Properties["address"]
	assert.True(t, ok)
	assert.Equal(t, "object", addrProp.Type)
	assert.Equal(t, "#/definitions/Address", addrProp.Ref)
}

func TestModelParser_Register_NestedStruct(t *testing.T) {
	parser := NewModelParser()

	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type User struct {
		ID      int     `json:"id"`
		Name    string  `json:"name"`
		Address Address `json:"address"` // Không phải pointer
	}

	name := parser.Register(User{})
	assert.Equal(t, "User", name)

	// Address schema được register
	_, ok := parser.schemas["Address"]
	assert.True(t, ok, "Address schema should be registered")

	userSchema, ok := parser.schemas["User"]
	assert.True(t, ok)

	addrProp, ok := userSchema.Properties["address"]
	assert.True(t, ok)
	assert.Equal(t, "object", addrProp.Type)
	assert.Equal(t, "#/definitions/Address", addrProp.Ref)
}

func TestModelParser_Register_SliceOfPointers(t *testing.T) {
	parser := NewModelParser()

	type Item struct {
		ID int `json:"id"`
	}

	type Order struct {
		Items []*Item `json:"items"`
	}

	parser.Register(Order{})

	items := parser.schemas["Order"].Properties["items"]

	assert.Equal(t, "array", items.Type)
	assert.Equal(t, "#/definitions/Item", items.Items.Ref)
}
func TestModelParser_Register_SliceOfStructs(t *testing.T) {
	parser := NewModelParser()

	type Item struct {
		ID int `json:"id"`
	}

	type Order struct {
		Items []Item `json:"items"`
	}

	parser.Register(Order{})

	_, ok := parser.schemas["Item"]
	assert.True(t, ok)

	items := parser.schemas["Order"].Properties["items"]

	assert.Equal(t, "array", items.Type)
	assert.Equal(t, "#/definitions/Item", items.Items.Ref)
}

func TestModelParser_Register_PointerToSlice(t *testing.T) {
	parser := NewModelParser()

	type Tag struct {
		Name string `json:"name"`
	}

	type Article struct {
		Tags *[]Tag `json:"tags"`
	}

	name := parser.Register(Article{})
	assert.Equal(t, "Article", name)

	_, ok := parser.schemas["Tag"]
	assert.True(t, ok)
}

func TestModelParser_Register_Nil(t *testing.T) {
	parser := NewModelParser()

	name := parser.Register(nil)

	assert.Equal(t, "unknown", name)
	assert.Empty(t, parser.schemas)
}

func TestModelParser_Register_Map(t *testing.T) {
	parser := NewModelParser()

	name := parser.Register(map[string]interface{}{})

	assert.Equal(t, "map", name)
	assert.Empty(t, parser.schemas)
}
func TestModelParser_RegisterType_Pointer(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		Name string `json:"name"`
	}

	name := parser.registerType(reflect.TypeOf(&User{}))

	assert.Equal(t, "User", name)
	assert.Contains(t, parser.schemas, "User")
}

func TestModelParser_RegisterType_Map(t *testing.T) {
	parser := NewModelParser()

	name := parser.registerType(
		reflect.TypeOf(map[string]string{}),
	)

	assert.Equal(t, "map", name)
}

func TestModelParser_RegisterType_Interface(t *testing.T) {
	parser := NewModelParser()

	var value interface{}
	tp := reflect.TypeOf((*interface{})(nil)).Elem()

	name := parser.registerType(tp)

	assert.Equal(t, "object", name)
	_ = value
}
func TestModelParser_RegisterType_NamedType(t *testing.T) {
	parser := NewModelParser()

	type UserID int

	name := parser.registerType(reflect.TypeOf(UserID(1)))

	assert.Equal(t, "UserID", name)
}
func TestModelParser_RegisterType_AnonymousType(t *testing.T) {
	parser := NewModelParser()

	tp := reflect.TypeOf([]string{})

	name := parser.registerType(tp)

	assert.Equal(t, "array", name)
}
func TestModelParser_RegisterType_AlreadyExists(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		ID int `json:"id"`
	}

	tp := reflect.TypeOf(User{})

	first := parser.registerType(tp)
	second := parser.registerType(tp)

	assert.Equal(t, "User", first)
	assert.Equal(t, "User", second)
	assert.Len(t, parser.schemas, 1)
}
func TestModelParser_Register_SkipJSONFields(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		ID       int    `json:"id"`
		Password string `json:"-"`
		Internal string
	}

	parser.Register(User{})

	schema := parser.schemas["User"]

	assert.Contains(t, schema.Properties, "id")
	assert.NotContains(t, schema.Properties, "Password")
	assert.NotContains(t, schema.Properties, "Internal")
}
func TestModelParser_Register_ExampleAndEnum(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		Name string `json:"name" example:"John"`
		Role string `json:"role" enum:"admin,user,guest"`
	}

	parser.Register(User{})

	schema := parser.schemas["User"]

	nameProp := schema.Properties["name"]
	assert.Equal(t, "John", nameProp.Example)

	roleProp := schema.Properties["role"]
	assert.Equal(t, []string{"admin", "user", "guest"}, roleProp.Enum)
}
func TestModelParser_Register_NotRequired(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age" binding:"omitempty"`
	}

	parser.Register(User{})

	schema := parser.schemas["User"]

	assert.Empty(t, schema.Required)
}
func TestModelParser_Register_InvalidNumericTags(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		Age      int    `json:"age" minimum:"abc" maximum:"xyz"`
		Username string `json:"username" minLength:"abc" maxLength:"xyz"`
	}

	parser.Register(User{})

	schema := parser.schemas["User"]

	age := schema.Properties["age"]
	assert.NotNil(t, age.Minimum)
	assert.NotNil(t, age.Maximum)
	assert.Equal(t, 0, *age.Minimum)
	assert.Equal(t, 0, *age.Maximum)

	username := schema.Properties["username"]
	assert.NotNil(t, username.MinLength)
	assert.NotNil(t, username.MaxLength)
	assert.Equal(t, 0, *username.MinLength)
	assert.Equal(t, 0, *username.MaxLength)
}
func TestModelParser_GetType_Pointer(t *testing.T) {
	parser := NewModelParser()

	assert.Equal(
		t,
		"string",
		parser.getType(reflect.TypeOf((*string)(nil))),
	)

	assert.Equal(
		t,
		"integer",
		parser.getType(reflect.TypeOf((*int)(nil))),
	)

	assert.Equal(
		t,
		"boolean",
		parser.getType(reflect.TypeOf((*bool)(nil))),
	)
}
func TestModelParser_GetType_IntegerTypes(t *testing.T) {
	parser := NewModelParser()

	assert.Equal(t, "integer", parser.getType(reflect.TypeOf(int8(1))))
	assert.Equal(t, "integer", parser.getType(reflect.TypeOf(int16(1))))
	assert.Equal(t, "integer", parser.getType(reflect.TypeOf(int32(1))))
	assert.Equal(t, "integer", parser.getType(reflect.TypeOf(int64(1))))

	assert.Equal(t, "integer", parser.getType(reflect.TypeOf(uint(1))))
	assert.Equal(t, "integer", parser.getType(reflect.TypeOf(uint8(1))))
	assert.Equal(t, "integer", parser.getType(reflect.TypeOf(uint16(1))))
	assert.Equal(t, "integer", parser.getType(reflect.TypeOf(uint32(1))))
	assert.Equal(t, "integer", parser.getType(reflect.TypeOf(uint64(1))))
}
func TestModelParser_GetType_Time(t *testing.T) {
	parser := NewModelParser()

	got := parser.getType(reflect.TypeOf(time.Time{}))

	assert.Equal(t, "string", got)
}
func TestModelParser_GetType_Default(t *testing.T) {
	parser := NewModelParser()

	assert.Equal(
		t,
		"string",
		parser.getType(reflect.TypeOf(make(chan int))),
	)
}
func TestModelParser_Register_SlicePrimitive(t *testing.T) {
	parser := NewModelParser()

	type Data struct {
		Names []string `json:"names"`
		IDs   []int    `json:"ids"`
	}

	parser.Register(Data{})

	schema := parser.schemas["Data"]

	names := schema.Properties["names"]
	assert.Equal(t, "array", names.Type)

	ids := schema.Properties["ids"]
	assert.Equal(t, "array", ids.Type)
}
func TestModelParser_GetType_Interface(t *testing.T) {
	parser := NewModelParser()

	typ := reflect.TypeOf((*interface{})(nil)).Elem()

	got := parser.getType(typ)

	assert.Equal(t, "object", got)
}
func TestModelParser_RegisterType_NamedPrimitive(t *testing.T) {
	parser := NewModelParser()

	type UserID int

	got := parser.registerType(reflect.TypeOf(UserID(1)))

	assert.Equal(t, "UserID", got)
}
func TestModelParser_RegisterType_ExistingSchema(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		ID int `json:"id"`
	}

	parser.Register(User{})

	got := parser.registerType(reflect.TypeOf(User{}))

	assert.Equal(t, "User", got)
	assert.Len(t, parser.schemas, 1)
}
func TestModelParser_Register_SkipJSONTag(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		ID       int    `json:"id"`
		Password string `json:"-"`
		Internal string
	}

	parser.Register(User{})

	schema := parser.schemas["User"]

	assert.Contains(t, schema.Properties, "id")
	assert.NotContains(t, schema.Properties, "Password")
	assert.NotContains(t, schema.Properties, "Internal")
}
func TestModelParser_Register_Example(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		Name string `json:"name" example:"John"`
	}

	parser.Register(User{})

	prop := parser.schemas["User"].Properties["name"]

	assert.Equal(t, "John", prop.Example)
}
func TestModelParser_Register_Enum(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		Role string `json:"role" enum:"admin,user,guest"`
	}

	parser.Register(User{})

	prop := parser.schemas["User"].Properties["role"]

	assert.Equal(t, []string{"admin", "user", "guest"}, prop.Enum)
}
func TestModelParser_GetType_SwitchBranches(t *testing.T) {
	parser := NewModelParser()

	tests := []struct {
		name string
		typ  reflect.Type
		want string
	}{
		{
			name: "uint",
			typ:  reflect.TypeOf(uint(1)),
			want: "integer",
		},
		{
			name: "uint8",
			typ:  reflect.TypeOf(uint8(1)),
			want: "integer",
		},
		{
			name: "uint16",
			typ:  reflect.TypeOf(uint16(1)),
			want: "integer",
		},
		{
			name: "uint32",
			typ:  reflect.TypeOf(uint32(1)),
			want: "integer",
		},
		{
			name: "uint64",
			typ:  reflect.TypeOf(uint64(1)),
			want: "integer",
		},
		{
			name: "float32",
			typ:  reflect.TypeOf(float32(1)),
			want: "number",
		},
		{
			name: "float64",
			typ:  reflect.TypeOf(float64(1)),
			want: "number",
		},
		{
			name: "bool",
			typ:  reflect.TypeOf(true),
			want: "boolean",
		},
		{
			name: "interface",
			typ:  reflect.TypeOf((*interface{})(nil)).Elem(),
			want: "object",
		},
		{
			name: "struct",
			typ:  reflect.TypeOf(struct{}{}),
			want: "object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parser.getType(tt.typ))
		})
	}
}

func TestModelParser_GetType_Uint(t *testing.T) {
	parser := NewModelParser()

	tests := []struct {
		name string
		typ  reflect.Type
	}{
		{"uint", reflect.TypeOf(uint(1))},
		{"uint8", reflect.TypeOf(uint8(1))},
		{"uint16", reflect.TypeOf(uint16(1))},
		{"uint32", reflect.TypeOf(uint32(1))},
		{"uint64", reflect.TypeOf(uint64(1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, "integer", parser.getType(tt.typ))
		})
	}
}
func TestModelParser_Register_NoJSONTag(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		Name string
	}

	parser.Register(User{})

	schema := parser.schemas["User"]

	assert.NotContains(t, schema.Properties, "Name")
}
func TestModelParser_Register_JSONTagDash(t *testing.T) {
	parser := NewModelParser()

	type User struct {
		Password string `json:"-"`
	}

	parser.Register(User{})

	schema := parser.schemas["User"]

	assert.NotContains(t, schema.Properties, "Password")
}
func TestModelParser_RegisterType_UnnamedType(t *testing.T) {
	parser := NewModelParser()

	typ := reflect.TypeOf(make(chan int))

	got := parser.registerType(typ)

	assert.Equal(t, "string", got)
}
