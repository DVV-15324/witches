package template

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	"github.com/example/project/internal/book"
)

type Modules struct {
	Book    *book.BookModule
}

func InitModules(core *core.CoreServices) *Modules {
	bookModule := book.NewBookModule(core)

	return &Modules{
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
	w_handl "github.com/DVV-15324/witches/pkg/core/handler"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
	"log"
	"new_example/cmd/server/core"
	"new_example/internal/shared/middleware"
)

func RegisterRoutes(r *gin.Engine, core *core.CoreServices, modules *Modules) {
	// Middleware CORS
	r.Use(w_utils.Cors(core.Config))
	// Rate limit middleware
	rateLimitMiddleware := &middleware.LimitMiddleWare{}
	// Khởi tạo Swagger
	gen := initSwagger(r, core, rateLimitMiddleware)
	initModule(modules, gen)
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

func initSwagger(r *gin.Engine, core *core.CoreServices, rateLimitMiddleware *middleware.LimitMiddleWare) *w_handl.SwaggerGenerator {
	address := fmt.Sprintf("%s:%d", core.Config.Host, core.Config.Port)
	gen := w_handl.NewSwaggerGenerator("new_example API", "1.0", address, "/")
	gen.SetEngine(r)
	gen.SetRateLimitMiddleware(rateLimitMiddleware)
	gen.SetRedisClient(core.Redis.GetClient())
	return gen
}

func initModule(modules *Modules, gen *w_handl.SwaggerGenerator) {
	gen.AddTag("v1", "API Version 1")
	gen.AddTag("book", "Book endpoints")
	modules.Book.RegisterPublicRoutes(gen)
	modules.Book.RegisterProtectedRoutes(gen)

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
func TestRollbackDomain_WarningBranches(t *testing.T) {
	t.Run("rollbackKeyObject error", func(t *testing.T) {
		tmpDir := t.TempDir()
		moduleName := "github.com/example/project"

		routersDir := filepath.Join(tmpDir, "cmd", "server", "routers")
		utilsDir := filepath.Join(tmpDir, "internal", "shared", "utils")

		require.NoError(t, os.MkdirAll(routersDir, 0755))
		require.NoError(t, os.MkdirAll(utilsDir, 0755))

		// key_object.go tồn tại nhưng không có ObjectBook
		keyContent := `package utils

var (
	ObjectUser int64 = 1
)
`

		require.NoError(t, os.WriteFile(
			filepath.Join(utilsDir, "key_object.go"),
			[]byte(keyContent),
			0644,
		))

		// modules.go tối thiểu để các bước sau vẫn chạy
		modulesContent := `package routers

type Modules struct {
	Book interface{}
}

func InitModules() *Modules {
	bookModule := 1

	return &Modules{
		Book: bookModule,
	}
}
`

		require.NoError(t, os.WriteFile(
			filepath.Join(routersDir, "modules.go"),
			[]byte(modulesContent),
			0644,
		))

		routersContent := `package routers

func initModule(modules *Modules) {
	modules.Book.AddTag("book")
}
`

		require.NoError(t, os.WriteFile(
			filepath.Join(routersDir, "routers.go"),
			[]byte(routersContent),
			0644,
		))

		err := RollbackDomain(tmpDir, moduleName, "book")

		// RollbackDomain vẫn thành công vì rollbackKeyObject
		// chỉ warning chứ không return error.
		require.NoError(t, err)
	})

	t.Run("rollbackModuleField error", func(t *testing.T) {
		tmpDir := t.TempDir()
		moduleName := "github.com/example/project"

		routersDir := filepath.Join(tmpDir, "cmd", "server", "routers")
		utilsDir := filepath.Join(tmpDir, "internal", "shared", "utils")

		require.NoError(t, os.MkdirAll(routersDir, 0755))
		require.NoError(t, os.MkdirAll(utilsDir, 0755))

		require.NoError(t, os.WriteFile(
			filepath.Join(utilsDir, "key_object.go"),
			[]byte(`package utils

var ObjectBook int64 = 1
`),
			0644,
		))

		// Không có struct Modules
		modulesContent := `package routers

func InitModules() {}
`

		require.NoError(t, os.WriteFile(
			filepath.Join(routersDir, "modules.go"),
			[]byte(modulesContent),
			0644,
		))

		require.NoError(t, os.WriteFile(
			filepath.Join(routersDir, "routers.go"),
			[]byte(`package routers

func initModule(modules interface{}) {}
`),
			0644,
		))

		err := RollbackDomain(tmpDir, moduleName, "book")
		require.NoError(t, err)
	})

	t.Run("rollbackModuleInit error", func(t *testing.T) {
		tmpDir := t.TempDir()
		moduleName := "github.com/example/project"

		routersDir := filepath.Join(tmpDir, "cmd", "server", "routers")
		utilsDir := filepath.Join(tmpDir, "internal", "shared", "utils")

		require.NoError(t, os.MkdirAll(routersDir, 0755))
		require.NoError(t, os.MkdirAll(utilsDir, 0755))

		require.NoError(t, os.WriteFile(
			filepath.Join(utilsDir, "key_object.go"),
			[]byte(`package utils

var ObjectBook int64 = 1
`),
			0644,
		))

		// Có Modules nhưng không có InitModules
		modulesContent := `package routers

type Modules struct {
	Book interface{}
}
`

		require.NoError(t, os.WriteFile(
			filepath.Join(routersDir, "modules.go"),
			[]byte(modulesContent),
			0644,
		))

		require.NoError(t, os.WriteFile(
			filepath.Join(routersDir, "routers.go"),
			[]byte(`package routers

func initModule(modules interface{}) {}
`),
			0644,
		))

		err := RollbackDomain(tmpDir, moduleName, "book")
		require.NoError(t, err)
	})

	t.Run("rollbackRouteRegistration error", func(t *testing.T) {
		tmpDir := t.TempDir()
		moduleName := "github.com/example/project"

		routersDir := filepath.Join(tmpDir, "cmd", "server", "routers")
		utilsDir := filepath.Join(tmpDir, "internal", "shared", "utils")

		require.NoError(t, os.MkdirAll(routersDir, 0755))
		require.NoError(t, os.MkdirAll(utilsDir, 0755))

		require.NoError(t, os.WriteFile(
			filepath.Join(utilsDir, "key_object.go"),
			[]byte(`package utils

var ObjectBook int64 = 1
`),
			0644,
		))

		require.NoError(t, os.WriteFile(
			filepath.Join(routersDir, "modules.go"),
			[]byte(`package routers

type Modules struct{}
`),
			0644,
		))

		// Không có initModule
		require.NoError(t, os.WriteFile(
			filepath.Join(routersDir, "routers.go"),
			[]byte(`package routers
`),
			0644,
		))

		err := RollbackDomain(tmpDir, moduleName, "book")
		require.NoError(t, err)
	})
}
func TestRollbackDomain_RemoveAllError(t *testing.T) {
	original := removeAll
	defer func() {
		removeAll = original
	}()

	removeAll = func(string) error {
		return errors.New("remove all failed")
	}

	err := RollbackDomain(t.TempDir(), "github.com/example/project", "book")

	require.EqualError(
		t,
		err,
		"failed to remove domain directory: remove all failed",
	)
}
func TestRollbackDomain_RemoveSharedFileError(t *testing.T) {
	originalRemove := remove
	defer func() {
		remove = originalRemove
	}()

	remove = func(string) error {
		return errors.New("remove failed")
	}

	tmpDir := t.TempDir()

	err := RollbackDomain(
		tmpDir,
		"github.com/example/project",
		"book",
	)

	require.EqualError(
		t,
		err,
		"failed to remove shared domain file: remove failed",
	)
}
func TestRollbackDomain_SharedFileNotExist(t *testing.T) {
	originalRemove := remove
	defer func() {
		remove = originalRemove
	}()

	remove = func(string) error {
		return os.ErrNotExist
	}

	tmpDir := t.TempDir()

	err := RollbackDomain(
		tmpDir,
		"github.com/example/project",
		"book",
	)

	require.NoError(t, err)
}
