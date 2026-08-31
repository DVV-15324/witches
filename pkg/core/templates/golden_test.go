package template

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"
)

func getTestdataDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get current test file path")
	}

	return filepath.Join(filepath.Dir(filename), "testdata")
}

func assertGolden(t *testing.T, namespace, actual, relPath string) {
	t.Helper()

	goldenPath := filepath.Join(
		getTestdataDir(),
		namespace,
		relPath+".golden",
	)

	if golden.FlagUpdate() {
		require.NoError(
			t,
			os.MkdirAll(filepath.Dir(goldenPath), 0755),
		)

		require.NoError(
			t,
			os.WriteFile(goldenPath, []byte(actual), 0644),
		)

		t.Logf("updated golden: %s", goldenPath)
		return
	}

	expected, err := os.ReadFile(goldenPath)
	require.NoError(
		t,
		err,
		"golden file not found: %s",
		goldenPath,
	)

	assert.Equal(
		t,
		string(expected),
		actual,
		"golden mismatch: %s",
		relPath,
	)
}

func TestRemoveGoDomain_Golden(t *testing.T) {
	tmpDir := t.TempDir()
	moduleName := "github.com/example/project"

	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(tmpDir, "go.mod"),
			[]byte("module "+moduleName),
			0644,
		),
	)

	routersDir := filepath.Join(
		tmpDir,
		"cmd",
		"server",
		"routers",
	)

	require.NoError(t, os.MkdirAll(routersDir, 0755))

	modulesContent := `package routers

import (
	"github.com/example/project/cmd/server/core"
	"github.com/example/project/internal/auth"
	"github.com/example/project/internal/book"
	"github.com/example/project/internal/refresh"
	"github.com/example/project/internal/user"
)

type Modules struct {
	Auth    *auth.AuthModule
	User    *user.UserModule
	Refresh *refresh.RefreshModule
	Book    *book.BookModule
}

func InitModules(core *core.CoreServices) *Modules {
	userModule := user.NewUserModule(core)
	refreshModule := refresh.NewRefreshModule(core, userModule.Usecase)
	authModule := auth.NewAuthModule(core, userModule.Usecase, refreshModule.Usecase)
	refreshModule.SetAuthUseCase(authModule.Usecase)

	bookModule := book.NewBookModule(core)

	return &Modules{
		Auth:    authModule,
		User:    userModule,
		Refresh: refreshModule,
		Book:    bookModule,
	}
}
`

	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(routersDir, "modules.go"),
			[]byte(modulesContent),
			0644,
		),
	)

	// routers.go với book đã được thêm vào
	routersContent := `package routers

import (
	"fmt"
	"log"

	"github.com/example/project/cmd/server/core"
	"github.com/example/project/internal/shared/middleware"

	w_handl "github.com/DVV-15324/witches/pkg/core/handler"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

func RegisterRoutes(r *gin.Engine, core *core.CoreServices, modules *Modules) {
	// Middleware CORS
	r.Use(w_utils.Cors(core.Config))

	// Rate limit middleware
	rateLimitMiddleware := &middleware.LimitMiddleWare{}

	// Auth middleware
	authMiddleware := middleware.AuthMiddleware(
		modules.Refresh.Usecase,
		core,
	)

	// Rate limit config
	rateLimit := limiter.Rate{
		Period: core.Limiter.Rate.Period,
		Limit:  core.Limiter.Rate.Limit,
	}

	// Khởi tạo Swagger
	gen := initSwagger(r, core, rateLimitMiddleware)
	initModule(modules, gen, rateLimit, authMiddleware)
	if err := gen.Save("swagger.json"); err != nil {
		log.Printf("Error saving swagger.json: %v", err)
	} else {
		log.Println("swagger.json generated successfully!")
	}

	// Swagger UI
	r.GET("/swagger/*any", w_handl.SwaggerUI())
	r.GET("/swagger.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(200, gen.GenerateJSON())
	})
}

func initSwagger(r *gin.Engine, core *core.CoreServices, rateLimitMiddleware w_handl.IRateLimitMiddleware) *w_handl.SwaggerGenerator {
	address := fmt.Sprintf("%s:%d", core.Config.Host, core.Config.Port)
	gen := w_handl.NewSwaggerGenerator("example API", "1.0", address, "/")
	gen.SetEngine(r)
	gen.SetRedisClient(core.Redis.GetClient())
	gen.SetRateLimitMiddleware(rateLimitMiddleware)
	return gen
}

func initModule(modules *Modules, gen *w_handl.SwaggerGenerator, rateLimit limiter.Rate, authMiddleware gin.HandlerFunc) {
	gen.AddTag("v1", "API Version 1")
	gen.AddTag("auth", "Auth endpoints")
	modules.Auth.RegisterProtectedRoutes(gen, &rateLimit, authMiddleware)
	modules.Auth.RegisterPublicRoutes(gen, &rateLimit)
	gen.AddTag("user", "User endpoints")
	modules.User.RegisterProtectedRoutes(gen, &rateLimit, authMiddleware)
	modules.User.RegisterPublicRoutes(gen, &rateLimit)
	gen.AddTag("refresh", "Refresh endpoints")
	modules.Refresh.RegisterPublicRoutes(gen, &rateLimit)
	modules.Refresh.RegisterProtectedRoutes(gen, &rateLimit, authMiddleware)
	gen.AddTag("book", "Book endpoints")
	modules.Book.RegisterPublicRoutes(gen,&rateLimit)
	modules.Book.RegisterProtectedRoutes(gen,&rateLimit,authMiddleware)
}

`

	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(routersDir, "routers.go"),
			[]byte(routersContent),
			0644,
		),
	)

	utilsDir := filepath.Join(
		tmpDir,
		"internal",
		"shared",
		"utils",
	)

	require.NoError(t, os.MkdirAll(utilsDir, 0755))

	keyContent := `package utils

var (
	ObjectUser   int64 = 1
	ObjectAuth   int64 = 2
	ObjectRefresh int64 = 3
	ObjectBook   int64 = 4
)
`

	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(utilsDir, "key_object.go"),
			[]byte(keyContent),
			0644,
		),
	)

	bookDir := filepath.Join(tmpDir, "internal", "book")
	require.NoError(t, os.MkdirAll(filepath.Join(bookDir, "handler"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(bookDir, "model"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(bookDir, "usecase"), 0755))

	handlerContent := `package handler

import (
	"github.com/example/project/internal/book/usecase"
)

type BookHandler struct {
	Usecase *usecase.BookUsecase
}

func NewBookHandler(usecase *usecase.BookUsecase) *BookHandler {
	return &BookHandler{
		Usecase: usecase,
	}
}
`
	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(bookDir, "handler", "handler.go"),
			[]byte(handlerContent),
			0644,
		),
	)

	modelContent := `package model

type Book struct {
	ID     string ` + "`json:\"id\"`" + `
	Title  string ` + "`json:\"title\"`" + `
	Author string ` + "`json:\"author\"`" + `
}
`
	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(bookDir, "model", "model.go"),
			[]byte(modelContent),
			0644,
		),
	)

	usecaseContent := `package usecase

import (
	"github.com/example/project/internal/book/model"
)

type BookUsecase struct{}

func NewBookUsecase() *BookUsecase {
	return &BookUsecase{}
}

func (u *BookUsecase) GetBook(id string) (*model.Book, error) {
	return &model.Book{
		ID:     id,
		Title:  "Test Book",
		Author: "Test Author",
	}, nil
}
`
	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(bookDir, "usecase", "usecase.go"),
			[]byte(usecaseContent),
			0644,
		),
	)

	moduleContent := `package book

import (
	"github.com/example/project/cmd/server/core"
	"github.com/example/project/internal/book/handler"
	"github.com/example/project/internal/book/usecase"

	w_handl "github.com/DVV-15324/witches/pkg/core/handler"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

type BookModule struct {
	Handler  *handler.BookHandler
	Usecase  *usecase.BookUsecase
}

func NewBookModule(core *core.CoreServices) *BookModule {
	usecase := usecase.NewBookUsecase()
	handler := handler.NewBookHandler(usecase)

	return &BookModule{
		Handler: handler,
		Usecase: usecase,
	}
}

func (m *BookModule) RegisterPublicRoutes(gen *w_handl.SwaggerGenerator, rateLimit *limiter.Rate) {
	// Public routes
}

func (m *BookModule) RegisterProtectedRoutes(gen *w_handl.SwaggerGenerator, rateLimit *limiter.Rate, authMiddleware gin.HandlerFunc) {
	// Protected routes
}
`
	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(bookDir, "module.go"),
			[]byte(moduleContent),
			0644,
		),
	)

	domainDir := filepath.Join(tmpDir, "internal", "shared", "domain")
	require.NoError(t, os.MkdirAll(domainDir, 0755))

	domainContent := `package domain

const (
	BookCollection = "books"
)
`
	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(domainDir, "book.go"),
			[]byte(domainContent),
			0644,
		),
	)

	err := RollbackDomain(tmpDir, moduleName, "book")
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(tmpDir, "internal", "book"))
	assert.True(t, os.IsNotExist(err), "book directory should be removed")

	_, err = os.Stat(filepath.Join(tmpDir, "internal", "shared", "domain", "book.go"))
	assert.True(t, os.IsNotExist(err), "shared domain file should be removed")

	expectedFiles := []string{
		"cmd/server/routers/modules.go",
		"cmd/server/routers/routers.go",
		"internal/shared/utils/key_object.go",
	}

	for _, relPath := range expectedFiles {
		actualPath := filepath.Join(tmpDir, relPath)

		actual, err := os.ReadFile(actualPath)
		require.NoError(
			t,
			err,
			"file %s not found",
			relPath,
		)

		assertGolden(t, "remove_domain", string(actual), relPath)
	}
}
