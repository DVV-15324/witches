package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

type RouteBuilder struct {
	gen         *SwaggerGenerator
	method      string
	path        string
	op          Operation
	handler     gin.HandlerFunc
	middlewares []gin.HandlerFunc
	rateLimit   *limiter.Limiter
}

func (b *RouteBuilder) RateLimit(rate limiter.Rate) *RouteBuilder {
	if b == nil {
		fmt.Printf("Warning: RouteBuilder is nil\n")
		return b
	}
	if b.gen == nil {
		fmt.Printf("Warning: generator is nil\n")
		return b
	}

	if b.gen.storeFactory == nil {
		fmt.Printf("Warning: storeFactory is nil\n")
		return b
	}
	store, err := b.gen.storeFactory()
	if err != nil {
		fmt.Printf("Error: Failed to create rate limit store: %v\n", err)
		return b
	}
	if store == nil {
		fmt.Printf("Warning: store is nil\n")
		return b
	}
	b.rateLimit = limiter.New(store, rate)
	return b
}

func (b *RouteBuilder) Use(middlewares ...gin.HandlerFunc) *RouteBuilder {
	b.middlewares = append(b.middlewares, middlewares...)
	return b
}

func (b *RouteBuilder) Summary(s string) *RouteBuilder {
	if b == nil {
		fmt.Printf("Warning: RouteBuilder is nil\n")
		return b
	}
	b.op.Summary = s
	return b
}

func (b *RouteBuilder) Description(d string) *RouteBuilder {
	if b == nil {
		fmt.Printf("Warning: RouteBuilder is nil\n")
		return b
	}
	b.op.Description = d
	return b
}

func (b *RouteBuilder) Tags(tags ...string) *RouteBuilder {
	if b == nil {
		fmt.Printf("Warning: RouteBuilder is nil\n")
		return b
	}
	b.op.Tags = tags
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
		if _, ok := model.(string); ok {
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

// Handler gán handler
func (b *RouteBuilder) Handler(h gin.HandlerFunc) *RouteBuilder {
	b.handler = h
	return b
}

func (b *RouteBuilder) Build() {
	fullPath := b.gen.doc.BasePath + b.path
	var handlers []gin.HandlerFunc

	if b.gen.globalSecurity != nil {
		handlers = append(handlers, b.gen.globalMiddlewares...)
	}
	if b.rateLimit != nil {
		rateLimitHandler := b.gen.rateLimitMiddleware.CreateRateLimitMiddleware(b.rateLimit)
		handlers = append(handlers, rateLimitHandler)
	}
	if len(b.middlewares) > 0 {
		handlers = append(handlers, b.middlewares...)
	}
	if b.handler != nil {
		handlers = append(handlers, b.handler)
	}
	switch b.method {
	case "get":
		b.gen.engine.GET(fullPath, handlers...)
	case "post":
		b.gen.engine.POST(fullPath, handlers...)
	case "put":
		b.gen.engine.PUT(fullPath, handlers...)
	case "delete":
		b.gen.engine.DELETE(fullPath, handlers...)
	case "patch":
		b.gen.engine.PATCH(fullPath, handlers...)
	case "options":
		b.gen.engine.OPTIONS(fullPath, handlers...)
	}
	if b.gen.doc.Paths[b.path] == nil {
		b.gen.doc.Paths[b.path] = make(PathItem)
	}
	b.gen.doc.Paths[b.path][b.method] = b.op
}
