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

// TestAddGoDomain_Golden kiểm tra sinh domain mới.
func TestAddGoDomain_Golden(t *testing.T) {
	tmpDir := t.TempDir()
	moduleName := "github.com/example/project"

	// ------------------------------------------------------------
	// 1. Tạo go.mod
	// ------------------------------------------------------------
	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(tmpDir, "go.mod"),
			[]byte("module "+moduleName),
			0644,
		),
	)

	// ------------------------------------------------------------
	// 2. Tạo routers/modules.go và routers/routers.go
	//    LƯU Ý:
	//    Đây phải là trạng thái TRƯỚC KHI thêm book.
	// ------------------------------------------------------------
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
	"github.com/example/project/internal/refresh"
	"github.com/example/project/internal/user"
)

type Modules struct {
	Auth    *auth.AuthModule
	User    *user.UserModule
	Refresh *refresh.RefreshModule
}

func InitModules(core *core.CoreServices) *Modules {
	userModule := user.NewUserModule(core)
	refreshModule := refresh.NewRefreshModule(core, userModule.Usecase)
	authModule := auth.NewAuthModule(core, userModule.Usecase, refreshModule.Usecase)
	refreshModule.SetAuthUseCase(authModule.Usecase)

	return &Modules{
		Auth:    authModule,
		User:    userModule,
		Refresh: refreshModule,
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
		core.Session,
		core.Blacklist,
		core.Logger,
		core.Config,
		core.DeviceHelper,
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

func initSwagger(
	r *gin.Engine,
	core *core.CoreServices,
	rateLimitMiddleware w_handl.IRateLimitMiddleware,
) *w_handl.SwaggerGenerator {
	address := fmt.Sprintf("%s:%d", core.Config.Host, core.Config.Port)

	gen := w_handl.NewSwaggerGenerator(
		"github.com/example/project API",
		"1.0",
		address,
		"/",
	)

	gen.SetEngine(r)
	gen.SetRedisClient(core.Redis.GetClient())
	gen.SetRateLimitMiddleware(rateLimitMiddleware)

	return gen
}

func initModule(
	modules *Modules,
	gen *w_handl.SwaggerGenerator,
	rateLimit limiter.Rate,
	authMiddleware gin.HandlerFunc,
) {
	gen.AddTag("v1", "API Version 1")
	gen.AddTag("auth", "Auth endpoints")

	modules.Auth.RegisterProtectedRoutes(
		gen,
		&rateLimit,
		authMiddleware,
	)

	modules.Auth.RegisterPublicRoutes(
		gen,
		&rateLimit,
	)

	gen.AddTag("user", "User endpoints")

	modules.User.RegisterProtectedRoutes(
		gen,
		&rateLimit,
		authMiddleware,
	)

	modules.User.RegisterPublicRoutes(
		gen,
		&rateLimit,
	)

	gen.AddTag("refresh", "Refresh endpoints")

	modules.Refresh.RegisterPublicRoutes(
		gen,
		&rateLimit,
	)

	modules.Refresh.RegisterProtectedRoutes(
		gen,
		&rateLimit,
		authMiddleware,
	)
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

	// ------------------------------------------------------------
	// 3. Tạo key_object.go
	// ------------------------------------------------------------
	utilsDir := filepath.Join(
		tmpDir,
		"internal",
		"shared",
		"utils",
	)

	require.NoError(t, os.MkdirAll(utilsDir, 0755))

	keyContent := `package utils

var (
	ObjectUser   uint = 1
	ObjectAuth   uint = 2
	ObjectRefresh uint = 3
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

	// ------------------------------------------------------------
	// 4. Gọi AddGoDomain
	// ------------------------------------------------------------
	AddGoDomain(tmpDir, moduleName, "book")

	// ------------------------------------------------------------
	// 5. Debug actual nếu cần
	// ------------------------------------------------------------
	actualModules, err := os.ReadFile(
		filepath.Join(
			tmpDir,
			"cmd",
			"server",
			"routers",
			"modules.go",
		),
	)
	require.NoError(t, err)

	t.Logf("GENERATED MODULES:\n%s", actualModules)

	// ------------------------------------------------------------
	// 6. Golden
	// ------------------------------------------------------------
	expectedFiles := []string{
		"internal/book/handler/handler.go",
		"internal/book/model/model.go",
		"internal/book/usecase/usecase.go",
		"internal/book/module.go",
		"internal/shared/domain/book.go",
		"cmd/server/routers/modules.go",
		"cmd/server/routers/routers.go",
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

		assertGolden(t, "add_domain", string(actual), relPath)
	}
}

// TestCreateProjectStructure_Golden kiểm tra tạo project.
func TestCreateProjectStructure_Golden(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.Chdir(origWd))
	}()

	// Generator hiện tại phụ thuộc cwd.
	require.NoError(t, os.Chdir(tmpDir))

	config := ProjectConfig{
		ModuleName: "github.com/example/project",
	}

	err = createProjectStructure(config, "postgres")
	require.NoError(t, err)

	importantFiles := []string{
		"main.go",
		"go.mod",
		"cmd/root.go",
		"cmd/server/routers/modules.go",
		"cmd/server/routers/routers.go",
		"internal/shared/utils/key_object.go",
		"internal/auth/handler/handler.go",
		"internal/auth/module.go",
		"internal/refresh/module.go",
		"pkg/redis/client.go",
		"migrate/migrations/1_create_table.up.sql",
		"migrate/migrations/1_drop_table.down.sql",
	}

	for _, relPath := range importantFiles {
		actualPath := filepath.Join(tmpDir, relPath)

		actual, err := os.ReadFile(actualPath)
		require.NoError(t, err, "file %s not found", relPath)

		assertGolden(t, "create_project", string(actual), relPath)
	}
}

func TestCreateGoArcRefresh(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, err := os.Getwd()
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.Chdir(origWd))
	}()

	require.NoError(t, os.Chdir(tmpDir))

	projectName := "github.com/example/project"
	dbType := "postgres"

	CreateGoArcRefresh(projectName, dbType)

	importantFiles := []string{
		"main.go",
		"go.mod",
		"cmd/root.go",
		"cmd/server/routers/modules.go",
		"cmd/server/routers/routers.go",
		"internal/shared/utils/key_object.go",
		"internal/auth/module.go",
		"internal/user/module.go",
		"internal/refresh/module.go",
	}

	for _, f := range importantFiles {
		assert.FileExists(
			t,
			filepath.Join(tmpDir, f),
			"missing file: %s",
			f,
		)
	}
}
