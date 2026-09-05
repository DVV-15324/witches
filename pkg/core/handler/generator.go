package handle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

var newRedisStore = sredis.NewStore

var getwd = os.Getwd
var mkdirAll = os.MkdirAll
var writeFile = os.WriteFile

type IRateLimitMiddleware interface {
	CreateRateLimitMiddleware(rateLimit *limiter.Limiter) gin.HandlerFunc
}

type SwaggerGenerator struct {
	doc                 *SwaggerDoc
	modelParser         *ModelParser
	engine              *gin.Engine
	globalMiddlewares   []gin.HandlerFunc
	globalSecurity      []SecurityRequirement
	storeFactory        func() (limiter.Store, error)
	redisClient         *redis.Client
	rateLimitMiddleware IRateLimitMiddleware
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
		globalSecurity:      []SecurityRequirement{},
		modelParser:         NewModelParser(),
		globalMiddlewares:   []gin.HandlerFunc{},
		rateLimitMiddleware: nil,
	}
}

func (g *SwaggerGenerator) SetEngine(engine *gin.Engine) *SwaggerGenerator {
	g.engine = engine
	return g
}

func (g *SwaggerGenerator) SetRedisClient(redisClient *redis.Client) *SwaggerGenerator {
	g.redisClient = redisClient
	g.storeFactory = func() (limiter.Store, error) {
		if redisClient == nil {
			return nil, fmt.Errorf("redis client not configured")
		}
		return newRedisStore(redisClient)
	}
	return g
}

func (g *SwaggerGenerator) SetRateLimitMiddleware(middleware IRateLimitMiddleware) *SwaggerGenerator {
	g.rateLimitMiddleware = middleware
	return g
}

func (g *SwaggerGenerator) Use(middlewares ...gin.HandlerFunc) *SwaggerGenerator {
	g.globalMiddlewares = append(g.globalMiddlewares, middlewares...)
	if g.engine != nil {
		g.engine.Use(middlewares...)
	}
	return g
}

func (g *SwaggerGenerator) AddTag(name, description string) {
	g.doc.Tags = append(g.doc.Tags, Tag{Name: name, Description: description})
}

func (g *SwaggerGenerator) AddRoute(method, path string, op Operation) {
	if g.doc.Paths[path] == nil {
		g.doc.Paths[path] = make(PathItem)
	}
	g.doc.Paths[path][method] = op
}

func (g *SwaggerGenerator) RegisterModel(model interface{}) string {
	return g.modelParser.Register(model)
}

func (g *SwaggerGenerator) GenerateJSON() string {

	for name, schema := range g.modelParser.GetSchemas() {
		g.doc.Definitions[name] = schema
	}

	data, _ := json.MarshalIndent(g.doc, "", "  ")
	return string(data)
}

func (g *SwaggerGenerator) Save(filename string) error {
	pwd, err := getwd()
	if err != nil {
		return err
	}

	swaggerDir := filepath.Join(pwd, "swagger")

	if err := mkdirAll(swaggerDir, os.ModePerm); err != nil {
		return err
	}

	filePath := filepath.Join(swaggerDir, filename)

	return writeFile(filePath, []byte(g.GenerateJSON()), 0644)
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
