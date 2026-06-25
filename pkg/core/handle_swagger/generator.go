package handle_swagger

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type SwaggerGenerator struct {
	doc         *SwaggerDoc
	modelParser *ModelParser
	engine      *gin.Engine

	globalSecurity []SecurityRequirement
}

func NewSwaggerGenerator(title, version, host, basePath string) *SwaggerGenerator {
	return &SwaggerGenerator{
		doc: &SwaggerDoc{
			Swagger: "2.0",
			Info: Info{
				Title:   title,
				Version: version,
			},
			Host:                host,
			BasePath:            basePath,
			Schemes:             []string{"http", "https"},
			Paths:               make(map[string]PathItem),
			Definitions:         make(map[string]Schema),
			SecurityDefinitions: make(map[string]SecurityScheme),
			Tags:                []Tag{},
		},
		globalSecurity: []SecurityRequirement{},
		modelParser:    NewModelParser(),
	}
}

// SetEngine gán engine và tạo group mặc định
func (g *SwaggerGenerator) SetEngine(engine *gin.Engine) *SwaggerGenerator {
	g.engine = engine

	return g
}

// Use thêm middleware vào group hiện tại
func (g *SwaggerGenerator) Use(middleware ...gin.HandlerFunc) *SwaggerGenerator {
	g.engine.Use(middleware...)
	return g
}

// Thêm vào SwaggerGenerator
func (g *SwaggerGenerator) AddBearerAuth(name string) {
	g.doc.SecurityDefinitions[name] = SecurityScheme{
		Type:        "apiKey",
		In:          "header",
		Name:        "Authorization",
		Description: "Bearer token authentication. Example: 'Bearer {token}'",
	}
}

// Hoặc chi tiết hơn
func (g *SwaggerGenerator) AddBearerAuthWithDescription(name, description string) {
	g.doc.SecurityDefinitions[name] = SecurityScheme{
		Type:        "apiKey",
		In:          "header",
		Name:        "Authorization",
		Description: description,
	}
}

// AddSecurityDefinition thêm định nghĩa security
func (g *SwaggerGenerator) AddSecurityDefinition(name string, scheme SecurityScheme) {
	g.doc.SecurityDefinitions[name] = scheme
}

// SetGlobalSecurity thiết lập security mặc định cho tất cả routes
func (g *SwaggerGenerator) SetGlobalSecurity(name string, scopes ...string) {
	security := make(SecurityRequirement)
	security[name] = scopes
	g.globalSecurity = append(g.globalSecurity, security)
}

// AddTag thêm tag để nhóm routes
func (g *SwaggerGenerator) AddTag(name, description string) {
	g.doc.Tags = append(g.doc.Tags, Tag{Name: name, Description: description})
}

// AddRoute thêm route
func (g *SwaggerGenerator) AddRoute(method, path string, op Operation) {
	if g.doc.Paths[path] == nil {
		g.doc.Paths[path] = make(PathItem)
	}
	g.doc.Paths[path][method] = op
}

// RegisterModel đăng ký model để gen schema
func (g *SwaggerGenerator) RegisterModel(model interface{}) string {
	return g.modelParser.Register(model)
}

// GenerateJSON sinh file swagger.json
func (g *SwaggerGenerator) GenerateJSON() string {
	// Merge definitions từ model parser
	for name, schema := range g.modelParser.GetSchemas() {
		g.doc.Definitions[name] = schema
	}

	data, _ := json.MarshalIndent(g.doc, "", "  ")
	return string(data)
}

// Save lưu ra file
func (g *SwaggerGenerator) Save(filename string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	swaggerDir := filepath.Join(pwd, "swagger")

	if err := os.MkdirAll(swaggerDir, os.ModePerm); err != nil {
		return err
	}

	filePath := filepath.Join(swaggerDir, filename)

	return os.WriteFile(filePath, []byte(g.GenerateJSON()), 0644)
}

func (g *SwaggerGenerator) POST(path string) *RouteBuilder {
	return &RouteBuilder{
		gen:    g,
		method: "post",
		path:   path,
		op: Operation{
			Consumes:  []string{"application/json"},
			Produces:  []string{"application/json"},
			Responses: make(map[string]Response),
		},
	}
}

func (g *SwaggerGenerator) GET(path string) *RouteBuilder {
	return &RouteBuilder{
		gen:    g,
		method: "get",
		path:   path,
		op: Operation{
			Produces:  []string{"application/json"},
			Responses: make(map[string]Response),
		},
	}
}

func (g *SwaggerGenerator) PUT(path string) *RouteBuilder {
	return &RouteBuilder{
		gen:    g,
		method: "put",
		path:   path,
		op: Operation{
			Consumes:  []string{"application/json"},
			Produces:  []string{"application/json"},
			Responses: make(map[string]Response),
		},
	}
}

func (g *SwaggerGenerator) DELETE(path string) *RouteBuilder {
	return &RouteBuilder{
		gen:    g,
		method: "delete",
		path:   path,
		op: Operation{
			Produces:  []string{"application/json"},
			Responses: make(map[string]Response),
		},
	}
}
func (g *SwaggerGenerator) PATCH(path string) *RouteBuilder {
	return &RouteBuilder{
		gen:    g,
		method: "patch",
		path:   path,
		op: Operation{
			Consumes:  []string{"application/json"},
			Produces:  []string{"application/json"},
			Responses: make(map[string]Response),
		},
	}
}

func (g *SwaggerGenerator) OPTIONS(path string) *RouteBuilder {
	return &RouteBuilder{
		gen:    g,
		method: "options",
		path:   path,
		op: Operation{
			Produces:  []string{"application/json"},
			Responses: make(map[string]Response),
		},
	}
}
