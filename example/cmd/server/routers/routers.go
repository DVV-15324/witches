package routers

import (
	"log"
	"time"
	"example/cmd/server/config"
	dtoAuth "example/internal/auth-service/dto/request"
	dtoRefresh "example/internal/refresh-service/dto/request"
	middleware "example/internal/shared/middleware"

	w_handl "github.com/DVV-15324/witches/pkg/core/handler"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
)

func Start(r *gin.Engine, host string, cfg *config.Config) {

	// 1. Khởi tạo services
	services := Services()
	if services == nil {
		log.Fatal("Failed to initialize services")
	}

	// 1.1. Them cors, limit cho services
	r.Use(middleware.Cors())
	limit := middleware.LimitMiddleWare{}
	rateLimitPeriods := time.Duration(cfg.RateLimitPeriod) * time.Second

	rateLimit := limiter.Rate{
		Period: rateLimitPeriods,
		Limit:  cfg.RateLimitMax,
	}
	// 2. Tạo Swagger Generator
	gen := initSwagger(r, services.RedisClient, &limit, host)
	// 3. Định nghĩa routes
	setupPublicRoutes(gen, services, &rateLimit)
	setupProtectedRoutes(gen, services, &rateLimit)

	// 4. Save swagger.json
	if err := gen.Save("swagger.json"); err != nil {
		log.Printf("Error saving swagger.json: %v", err)
	} else {
		log.Println("swagger.json generated successfully!")
	}

	// 5. Swagger UI
	r.GET("/swagger/*any", w_handl.SwaggerUI())
	r.GET("/swagger.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(200, gen.GenerateJSON())
	})
}

func initSwagger(r *gin.Engine, redisClient *redis.Client, rateLimitMiddleware w_handl.IRateLimitMiddleware, host string) *w_handl.SwaggerGenerator {
	gen := w_handl.NewSwaggerGenerator("API", "1.0", host, "/")
	gen.SetEngine(r)
	gen.SetRedisClient(redisClient)
	gen.SetRateLimitMiddleware(rateLimitMiddleware)
	// Tags
	gen.AddTag("v1", "API Version 1")
	gen.AddTag("auth", "Authentication endpoints")
	gen.AddTag("user", "User management endpoints")

	// Register models
	gen.RegisterModel(dtoAuth.Login{})
	gen.RegisterModel(dtoAuth.Register{})
	gen.RegisterModel(dtoRefresh.RefreshTokenRequest{})
	gen.RegisterModel(w_utils.PaginationRequest{})
	gen.RegisterModel(w_resp.AppResponse{})

	return gen
}
