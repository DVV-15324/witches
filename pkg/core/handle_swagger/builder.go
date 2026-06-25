package handle_swagger

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type RouteBuilder struct {
	gen     *SwaggerGenerator
	method  string
	path    string
	op      Operation
	handler gin.HandlerFunc
}

func (b *RouteBuilder) Summary(s string) *RouteBuilder {
	b.op.Summary = s
	return b
}

func (b *RouteBuilder) Description(d string) *RouteBuilder {
	b.op.Description = d
	return b
}

func (b *RouteBuilder) Tags(tags ...string) *RouteBuilder {
	b.op.Tags = tags
	return b
}

// Security với scopes (cho OAuth2)
func (b *RouteBuilder) Security(name string, scopes ...string) *RouteBuilder {
	security := make(SecurityRequirement)
	security[name] = scopes

	b.op.Security = append(b.op.Security, security)
	return b
}

// BearerAuth tiện lợi cho Bearer token
func (b *RouteBuilder) BearerAuth() *RouteBuilder {
	security := make(SecurityRequirement)
	security["BearerAuth"] = []string{} // Không có scope cho Bearer token

	b.op.Security = append(b.op.Security, security)
	return b
}

func (b *RouteBuilder) PathParam(name, desc string, required bool) *RouteBuilder {
	b.op.Parameters = append(b.op.Parameters, Parameter{
		Name:        name,
		In:          "path",
		Description: desc,
		Required:    required,
		Type:        "string",
	})
	return b
}

func (b *RouteBuilder) QueryParam(name, desc string, required bool) *RouteBuilder {
	b.op.Parameters = append(b.op.Parameters, Parameter{
		Name:        name,
		In:          "query",
		Description: desc,
		Required:    required,
		Type:        "string",
	})
	return b
}

func (b *RouteBuilder) Body(model interface{}, desc string) *RouteBuilder {
	modelName := b.gen.RegisterModel(model)
	b.op.Parameters = append(b.op.Parameters, Parameter{
		Name:        "body",
		In:          "body",
		Description: desc,
		Required:    true,
		Schema: &SchemaRef{
			Ref: "#/definitions/" + modelName,
		},
	})
	return b
}
func (b *RouteBuilder) HeaderParam(name, desc string, required bool) *RouteBuilder {
	b.op.Parameters = append(b.op.Parameters, Parameter{
		Name:        name,
		In:          "header",
		Description: desc,
		Required:    required,
		Type:        "string",
	})
	return b
}
func (b *RouteBuilder) Response(code int, model interface{}, desc string) *RouteBuilder {
	resp := Response{
		Description: desc,
	}

	if model != nil {
		// Kiểm tra nếu là string (không cần schema)
		if _, ok := model.(string); ok {
			// Không thêm schema
		} else {
			modelName := b.gen.RegisterModel(model)
			resp.Schema = &SchemaRef{
				Ref: "#/definitions/" + modelName,
			}
		}
	}

	b.op.Responses[fmt.Sprintf("%d", code)] = resp
	return b
}

// Public đánh dấu route không cần authentication
func (b *RouteBuilder) Public() *RouteBuilder {
	// Security trống sẽ override global security
	b.op.Security = []SecurityRequirement{}
	return b
}
func (b *RouteBuilder) Handler(h gin.HandlerFunc) *RouteBuilder {
	b.handler = h

	// Tạo fullPath = BasePath + path
	fullPath := b.gen.doc.BasePath + b.path

	switch b.method {
	case "get":
		b.gen.engine.GET(fullPath, h)
	case "post":
		b.gen.engine.POST(fullPath, h)
	case "put":
		b.gen.engine.PUT(fullPath, h)
	case "delete":
		b.gen.engine.DELETE(fullPath, h)
	case "patch":
		b.gen.engine.PATCH(fullPath, h)
	case "options":
		b.gen.engine.OPTIONS(fullPath, h)
	}

	return b
}

// Build hoàn tất route
func (b *RouteBuilder) Build() {
	if b.gen.doc.Paths[b.path] == nil {
		b.gen.doc.Paths[b.path] = make(PathItem)
	}
	b.gen.doc.Paths[b.path][b.method] = b.op
}
