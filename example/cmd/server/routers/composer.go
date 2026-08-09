package routers

import (
	"context"
	config "example/cmd/server/config"
	handleAuth "example/internal/auth-service/handler"
	responsitoryAuth "example/internal/auth-service/repository"
	usecaseAuth "example/internal/auth-service/usecase"
	handleRefresh "example/internal/refresh-service/handler"
	responsitoryRefresh "example/internal/refresh-service/repository"
	usecaseRefresh "example/internal/refresh-service/usecase"
	"example/internal/shared/utils"
	handleUser "example/internal/user-service/handler"
	responsitoryUser "example/internal/user-service/repository"
	usecaseUser "example/internal/user-service/usecase"
	pkg_redis "example/pkg/redis"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	w_db "github.com/DVV-15324/witches/pkg/core/database/connect/sql"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/DVV-15324/witches/pkg/core/response/logger"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	redis_driver "github.com/ulule/limiter/v3/drivers/store/redis"
)

type IRefreshUseCase interface {
	IntrospectToken(ctx context.Context, accessToken string) (*w_utils.JwtClaims, *w_resp.AppError)
}

type IRefreshHandle interface {
	HandleRefreshToken() func(c *gin.Context)
}

type IAuthHandle interface {
	HandleLogin() func(c *gin.Context)
	HandleRegister() func(c *gin.Context)
	HandleLogout() func(c *gin.Context)
	HandleGoogleLogin() func(c *gin.Context)
}
type IUserHandle interface {
	HandleGetAllUser() func(c *gin.Context)
	HandleGetUserById() func(c *gin.Context)
}

type HandleServices struct {
	UsecaseRefresh IRefreshUseCase
	HandleRefresh  IRefreshHandle
	HandleAuth     IAuthHandle
	HandleUser     IUserHandle
	Logger         *logger.EntityLogger
	Session        *w_utils.SessionService
	Blacklist      *w_utils.BlacklistService
	RateLimit      *limiter.Limiter
	RedisClient    *redis.Client
}

// En: The constructor in routes is responsible for creating and connecting all the application components together.
//
// Vi: Hàm khởi tạo trong routes, có nhiệm vụ tạo và kết nối tất cả các thành phần của ứng dụng lại với nhau.
func Services() *HandleServices {
	// 1.
	// En: Load configuration from the .env file.
	// Vi: Nạp cấu hình từ file .env
	config := config.Load()

	// 2.
	// En: Get current path
	// Vi: Lấy đường dẫn thư mục hiện tại
	currentPath, _ := os.Getwd()

	// 3.
	// En: Combine the following paths: ./logs/logs.log
	// Vi: Ghép đường dẫn: ./logs/logs.log
	path := filepath.Join(currentPath, "logs", "logs.log")

	// 4.
	// En: Create a file logger with rotation (1MB, 20 backups, 30 days)
	// Vi: Tạo file logger với rotation (1MB, 20 backups, 30 ngày)
	logg, err := logger.NewFileLogger(path, 1, 20, 30)
	if err != nil {
		log.Println(err)
	}

	// 5.
	// En: Connect the database to the driver and DNS.
	// Vi: Kết nối database với driver và DSN
	db, err := w_db.NewDatabaseInstance(config.DBDriver, config.DBURL, logg, 2*time.Second, utils.ReqKey)

	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	if db == nil {
		log.Fatal(" Database connection is nil")
	}
	if db.DB == nil {
		log.Fatal(" Database DB instance is nil")
	}

	// 6.1
	// En: Connect to Redis cache
	// Vi: Kết nối Redis cache
	address := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
	redis, err := pkg_redis.NewRedisClient(address, config.RedisPassword)
	if err != nil {
		log.Println("Error: Connect database fail")
	}
	//6.2
	durationRateLimitPeriod := time.Duration(config.RateLimitPeriod) * time.Second
	// 100 requests / phút
	rate := limiter.Rate{
		Period: durationRateLimitPeriod,
		Limit:  config.RateLimitMax,
	}

	store, _ := redis_driver.NewStore(redis.GetClient())
	instanceLimit := limiter.New(store, rate)

	// 7.
	// En: Create services session and blacklist
	// Vi: Tạo services session và blacklist
	sessionService := w_utils.NewSessionService(redis.GetClient(), config.SessionTTL, config.IdleTimeout)
	blacklistService := w_utils.NewBlacklistService(redis.GetClient(), config.RevokedTTL)

	//8.
	// En: Create transaction manager
	// Vi: Tạo transaction manager
	tx := w_utils.NewTxManager(db.DB)

	// 9.
	// En: Create repositories
	// Vi: Tạo repositories
	rAuth := responsitoryAuth.NewAuthRepository(db.DB, tx, redis.GetClient())
	rUser := responsitoryUser.NewUserRepository(db.DB, tx)
	rRefresh := responsitoryRefresh.NewRefreshTokenRepository(db.DB, redis.GetClient())

	// 10.
	// En: Create usecase
	// Vi: Tạo usecase
	jwt := w_utils.NewJwtService("vu-dep-trai-nhat-the-gioi", 900, 604800)
	hash := new(w_utils.Hash)
	usecaseUser := usecaseUser.NewUserUsecase(rUser, tx)
	usecaseRefresh := usecaseRefresh.NewRefreshUseCase(usecaseUser, rRefresh, jwt, sessionService, blacklistService, config)
	usecaseAuth := usecaseAuth.NewAuthUseCase(jwt, usecaseUser, hash, rAuth, usecaseRefresh, tx, sessionService, blacklistService, config)
	usecaseRefresh.SetAuthUseCase(usecaseAuth)

	// 11.
	// En: Create handle
	// Vi: Tạo handle
	handleUser := handleUser.NewUserHandle(usecaseUser, logg)
	handleAuth := handleAuth.NewAuthHandle(usecaseAuth, logg)
	handleRefresh := handleRefresh.NewRefreshHandle(usecaseRefresh, logg)

	return &HandleServices{
		HandleRefresh:  handleRefresh,
		HandleAuth:     handleAuth,
		HandleUser:     handleUser,
		Logger:         logg,
		UsecaseRefresh: usecaseRefresh,
		Blacklist:      blacklistService,
		Session:        sessionService,
		RateLimit:      instanceLimit,
		RedisClient:    redis.GetClient(),
	}
}
