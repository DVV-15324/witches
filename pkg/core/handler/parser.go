package handle_swagger

import (
	"fmt"
	"reflect"
	"strings"
)

type ModelParser struct {
	schemas map[string]Schema
}

func NewModelParser() *ModelParser {
	return &ModelParser{
		schemas: make(map[string]Schema),
	}
}

func (p *ModelParser) Register(model interface{}) string {
	if model == nil {
		return "unknown"
	}

	t := reflect.TypeOf(model)

	//  Nếu là slice hoặc array
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		elemType := t.Elem()
		// Đăng ký element type trước
		elemName := p.registerType(elemType)
		// Tạo schema cho slice
		sliceName := "[]" + elemName
		if _, exists := p.schemas[sliceName]; !exists {
			p.schemas[sliceName] = Schema{
				Type: "array",
				Items: &SchemaRef{
					Ref: "#/definitions/" + elemName,
				},
			}
		}
		return sliceName
	}

	return p.registerType(t)
}

func (p *ModelParser) registerType(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return t.Name()
	}

	name := t.Name()
	if _, exists := p.schemas[name]; exists {
		return name
	}

	schema := Schema{
		Type:       "object",
		Properties: make(map[string]Property),
		Required:   []string{},
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		jsonName := strings.Split(jsonTag, ",")[0]

		prop := Property{
			Type:        p.getType(field.Type),
			Description: field.Tag.Get("description"),
		}

		if example := field.Tag.Get("example"); example != "" {
			prop.Example = example
		}

		if binding := field.Tag.Get("binding"); binding != "" {
			if strings.Contains(binding, "required") {
				schema.Required = append(schema.Required, jsonName)
			}
		}

		if min := field.Tag.Get("minimum"); min != "" {
			var val int
			fmt.Sscanf(min, "%d", &val)
			prop.Minimum = &val
		}
		if max := field.Tag.Get("maximum"); max != "" {
			var val int
			fmt.Sscanf(max, "%d", &val)
			prop.Maximum = &val
		}

		if minLen := field.Tag.Get("minLength"); minLen != "" {
			var val int
			fmt.Sscanf(minLen, "%d", &val)
			prop.MinLength = &val
		}
		if maxLen := field.Tag.Get("maxLength"); maxLen != "" {
			var val int
			fmt.Sscanf(maxLen, "%d", &val)
			prop.MaxLength = &val
		}

		if enum := field.Tag.Get("enum"); enum != "" {
			prop.Enum = strings.Split(enum, ",")
		}

		// ✅ Nếu field là slice, set items
		if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array {
			elemType := field.Type.Elem()
			if elemType.Kind() == reflect.Struct {
				elemName := p.registerType(elemType)
				prop.Type = "array"
				prop.Items = &SchemaRef{
					Ref: "#/definitions/" + elemName,
				}
			}
		}

		schema.Properties[jsonName] = prop
	}

	p.schemas[name] = schema
	return name
}

func (p *ModelParser) getType(t reflect.Type) string {
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		return "array"
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Interface:
		return "object" // or omit the type if you want to allow anything
	case reflect.Struct:
		if t.PkgPath() == "time" && t.Name() == "Time" {
			return "string"
		}
		return "object"
	default:
		return "string"
	}
}

func (p *ModelParser) GetSchemas() map[string]Schema {
	return p.schemas
}
