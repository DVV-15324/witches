package routes

import (
	dtoAuth "example/internal/dto/auth/request"
	"fmt"
	//dtoUser "example/internal/dto/user/request"
	middleware "example/internal/middleware"

	h_s "github.com/DVV-15324/witches/pkg/core/handle_swagger"
	u_response "github.com/DVV-15324/witches/pkg/core/response_logger"
	"github.com/gin-gonic/gin"
)

func Start(r *gin.Engine) {

	// Tạo Swagger Generator
	gen := h_s.NewSwaggerGenerator(
		"My API",
		"1.0",
		"localhost:8080",
		"/api",
	)
	gen.SetEngine(r)

	gen.AddSecurityDefinition(
		"BearerAuth",
		h_s.SecurityScheme{
			Type:        "apiKey",
			Name:        "Authorization",
			In:          "header",
			Description: "JWT Authorization header using the Bearer scheme. Example: 'Bearer {token}'",
		},
	)

	// Thêm tags
	gen.AddTag("v1", "API Version 1")

	// Đăng ký models dto
	gen.RegisterModel(dtoAuth.Login{})
	gen.RegisterModel(dtoAuth.Register{})

	// Đăng ký models ResponseHandle
	gen.RegisterModel(u_response.ResponseHandle{})

	services := Services()

	gen.POST("/v1/auth/login").
		Summary("User login").
		Description("Authenticate user and return JWT token").
		Tags("auth").
		Body(dtoAuth.Login{}, "Login credentials").
		Response(200, u_response.ResponseHandle{}, "Login successful").
		Response(401, u_response.ResponseHandle{}, "Unauthorized").
		Handler(services.HandleAuth.HandleLogin()).
		Build()
	gen.POST("/v1/auth/register").
		Summary("User registration").
		Description("Register a new user").
		Tags("auth").
		Body(dtoAuth.Register{}, "User registration data").
		Response(201, u_response.ResponseHandle{}, "Registration successful").
		Response(400, u_response.ResponseHandle{}, "Bad request").
		Build()
	gen.POST("/v1/user/get_user").
		Summary("Get user by ID").
		Description("Get user details by user ID").
		Tags("user").
		Security("BearerAuth").
		Body(map[string]int{"id": 0}, "User ID").
		Handler(services.HandleUser.HandleGetUserById()).
		Response(200, u_response.ResponseHandle{}, "User found").
		Response(404, u_response.ResponseHandle{}, "User not found").
		Use(middleware.RequiredAuth(services.UsecaseAuth, services.Logger)).
		Build()

	gen.POST("/v1/user/get_user_all").
		Summary("Get all users").
		Description("Get list of all users").
		Tags("user").
		Security("BearerAuth").
		Response(200, u_response.ResponseHandle{}, "Users list").
		Use(middleware.RequiredAuth(services.UsecaseAuth, services.Logger)).
		Handler(services.HandleUser.HandleGetAllUser()).
		Build()

		// Save swagger.json
	if err := gen.Save("swagger.json"); err != nil {
		fmt.Println("Lỗi lưu swagger.json:", err)
	} else {
		fmt.Println("swagger.json generated!")
	}

	// Swagger UI
	r.GET("/swagger/*any", h_s.SwaggerUI())

	r.GET("/swagger.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(200, gen.GenerateJSON())
	})

}
