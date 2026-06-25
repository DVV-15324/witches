package routes

import (
	handleAuth "arc-golang/internal/handler/auth"
	handleUser "arc-golang/internal/handler/user"
	responsitoryAuth "arc-golang/internal/repository/auth"
	responsitoryUser "arc-golang/internal/repository/user"
	usecaseAuth "arc-golang/internal/usecase/auth"
	usecaseUser "arc-golang/internal/usecase/user"
	c_db "arc-golang/internal/utils"
	c_hash "arc-golang/internal/utils"
	c_jwt "arc-golang/internal/utils"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type UsecaseAuth interface {
	UsecaseIntrospectToken(ctx context.Context, accessToken string) (*c_jwt.JwtClaims, error)
}

type IHandleAuth interface {
	HandleLogin() func(c *gin.Context)
	HandleRegister() func(c *gin.Context)
}
type IHandleUser interface {
	HandleGetAllUser() func(c *gin.Context)
	HandleGetUserById() func(c *gin.Context)
}

type HandleServices struct {
	UsecaseAuth UsecaseAuth
	HandleAuth  IHandleAuth
	HandleUser  IHandleUser
}

func Services() *HandleServices {
	err := godotenv.Load()
	if err != nil {
		log.Println("loi load .env")
	}
	DB_URL := os.Getenv("DB_URL")
	db, err := c_db.ConnectPostgreSQL(DB_URL)

	if err != nil {
		log.Println("loi connect database")
	}

	rAuth := responsitoryAuth.NewRepositoryAuth(db)
	rUser := responsitoryUser.NewRepositoryUser(db)
	// usecase
	jwt := c_jwt.NewJwtServer("vu-dep-trai-nhat-the-gioi", 604800)
	hash := new(c_hash.Hash)
	usecaseUser := usecaseUser.NewUsecaseUser(rUser)
	usecaseAuth := usecaseAuth.NewUsecaseAuth(jwt, usecaseUser, hash, rAuth)

	// handle
	handleUser := handleUser.NewHandleUser(usecaseUser)
	handleAuth := handleAuth.NewHandleUser(usecaseAuth)

	return &HandleServices{
		UsecaseAuth: usecaseAuth,
		HandleAuth:  handleAuth,
		HandleUser:  handleUser,
	}
}
