package routes

import (
	middleware "arc-golang/internal/middleware"
	//"log"

	"github.com/gin-gonic/gin"
)

func Start(r *gin.Engine) {
	services := Services()

	// GROUP
	v1 := r.Group("v1")

	//CORS
	v1.Use(middleware.Cors())

	//auth
	auth := v1.Group("auth")
	auth.POST("/login", services.HandleAuth.HandleLogin())
	auth.POST("/register", services.HandleAuth.HandleRegister())

	//user
	userV2 := v1.Group("user").Use(middleware.RequiredAuth(services.UsecaseAuth))
	userV2.POST("/get_user", services.HandleUser.HandleGetUserById())
	userV2.POST("/get_user_all", services.HandleUser.HandleGetAllUser())

}
