package test

import (
	"fmt"
	h_s "github.com/DVV-15324/witches/pkg/core/handle_swagger"
	"testing"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    int    `json:"id" description:"User ID" example:"1"`
	Name  string `json:"name" binding:"required" description:"User name" example:"vu" maxLength:"100"`
	Email string `json:"email" binding:"required" description:"User email" example:"vu@example.com"`
	Age   int    `json:"age" description:"User age" example:"25" minimum:"0" maximum:"150"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required" description:"User name" example:"vu"`
	Email    string `json:"email" binding:"required" description:"User email" example:"vu@example.com"`
	Password string `json:"password" binding:"required" description:"User password" example:"secret123"`
	Age      int    `json:"age" description:"User age" example:"25"`
}

type ResponseWrapper struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func getUsersV1(c *gin.Context) {
	c.JSON(200, ResponseWrapper{
		Status:  200,
		Message: "Success (v1)",
		Data: []User{
			{ID: 1, Name: "vu", Email: "vu@example.com", Age: 25},
			{ID: 2, Name: "vu", Email: "vu@example.com", Age: 30},
		},
	})
}

func createUserV1(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ResponseWrapper{
			Status:  400,
			Message: "Bad Request: " + err.Error(),
		})
		return
	}
	c.JSON(201, ResponseWrapper{
		Status:  201,
		Message: "Created (v1)",
		Data:    req,
	})
}

func getUserByIDV1(c *gin.Context) {
	c.JSON(200, ResponseWrapper{
		Status:  200,
		Message: "Success (v1)",
		Data: User{
			ID:    1,
			Name:  "vu",
			Email: "vu@example.com",
			Age:   25,
		},
	})
}

func updateUserV1(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ResponseWrapper{
			Status:  400,
			Message: "Bad Request: " + err.Error(),
		})
		return
	}
	c.JSON(200, ResponseWrapper{
		Status:  200,
		Message: "Updated (v1)",
		Data:    req,
	})
}

func deleteUserV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, ResponseWrapper{
		Status:  200,
		Message: "Deleted (v1)",
		Data:    gin.H{"id": id},
	})
}

func getUsersV2(c *gin.Context) {
	c.JSON(200, ResponseWrapper{
		Status:  200,
		Message: "Success (v2)",
		Data: []User{
			{ID: 1, Name: "vu v2", Email: "vu.v2@example.com", Age: 25},
			{ID: 2, Name: "vu v2", Email: "vu.v2@example.com", Age: 30},
		},
	})
}

func createUserV2(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ResponseWrapper{
			Status:  400,
			Message: "Bad Request (v2): " + err.Error(),
		})
		return
	}
	c.JSON(201, ResponseWrapper{
		Status:  201,
		Message: "Created (v2)",
		Data:    req,
	})
}

func getUserByIDV2(c *gin.Context) {
	c.JSON(200, ResponseWrapper{
		Status:  200,
		Message: "Success (v2)",
		Data: User{
			ID:    1,
			Name:  "vu v2",
			Email: "vu.v2@example.com",
			Age:   25,
		},
	})
}

func updateUserV2(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ResponseWrapper{
			Status:  400,
			Message: "Bad Request (v2): " + err.Error(),
		})
		return
	}
	c.JSON(200, ResponseWrapper{
		Status:  200,
		Message: "Updated (v2)",
		Data:    req,
	})
}

func deleteUserV2(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, ResponseWrapper{
		Status:  200,
		Message: "Deleted (v2)",
		Data:    gin.H{"id": id},
	})
}

func TestHandleSwagger(t *testing.T) {
	r := gin.Default()

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
	gen.AddSecurityDefinition("OAuth2", h_s.SecurityScheme{
		Type:             "oauth2",
		Flow:             "implicit",
		AuthorizationURL: "https://auth.example.com/authorize",
		Scopes: map[string]string{
			"read":  "Read access",
			"write": "Write access",
		},
	})
	// Thêm tags
	gen.AddTag("v1", "API Version 1")
	gen.AddTag("v2", "API Version 2")

	// Đăng ký models
	gen.RegisterModel(User{})
	gen.RegisterModel(CreateUserRequest{})
	gen.RegisterModel(ResponseWrapper{})

	// GET /api/v1/users
	gen.GET("/v1/users").
		Summary("List all users").
		Description("Returns a list of all users - Version 1").
		Tags("v1").
		QueryParam("page", "Page number", false).
		QueryParam("limit", "Items per page", false).
		Response(200, ResponseWrapper{Data: User{}}, "Success").
		Response(500, nil, "Internal Server Error").
		Security("BearerAuth").
		Handler(getUsersV1). //  Truyền thẳng handler
		Build()

	// POST /api/v1/users
	gen.POST("/v1/users").
		Summary("Create user").
		Description("Creates a new user - Version 1").
		Tags("v1").
		Security("OAuth2", "write").
		Body(CreateUserRequest{}, "User data").
		Response(201, User{}, "Created").
		Response(400, nil, "Bad Request").
		Handler(createUserV1).
		Build()

	// GET /api/v1/users/{id}
	gen.GET("/v1/users/{id}").
		Summary("Get user by ID").
		Description("Returns a single user - Version 1").
		Tags("v1").
		PathParam("id", "User ID", true).
		Response(200, ResponseWrapper{Data: User{}}, "Success").
		Response(404, nil, "Not Found").
		Handler(getUserByIDV1).
		Build()

	// PUT /api/v1/users/{id}
	gen.PUT("/v1/users/{id}").
		Summary("Update user").
		Description("Updates an existing user - Version 1").
		Tags("v1").
		PathParam("id", "User ID", true).
		Body(CreateUserRequest{}, "User data").
		Response(200, User{}, "Updated").
		Response(400, nil, "Bad Request").
		Response(404, nil, "Not Found").
		Handler(updateUserV1).
		Build()

	// DELETE /api/v1/users/{id}
	gen.DELETE("/v1/users/{id}").
		Summary("Delete user").
		Description("Deletes a user - Version 1").
		Tags("v1").
		PathParam("id", "User ID", true).
		Response(204, nil, "No Content").
		Response(404, nil, "Not Found").
		Handler(deleteUserV1).
		Build()

	// GET /api/v2/users
	gen.GET("/v2/users").
		Summary("List all users").
		Description("Returns a list of all users - Version 2").
		Tags("v2").
		QueryParam("page", "Page number", false).
		QueryParam("limit", "Items per page", false).
		Response(200, ResponseWrapper{Data: User{}}, "Success").
		Response(500, nil, "Internal Server Error").
		Security("BearerAuth").
		Handler(getUsersV2).
		Build()

	// POST /api/v2/users
	gen.POST("/v2/users").
		Summary("Create user").
		Description("Creates a new user - Version 2").
		Tags("v2").
		Body(CreateUserRequest{}, "User data").
		Response(201, User{}, "Created").
		Response(400, nil, "Bad Request").
		Handler(createUserV2).
		Build()

	// GET /api/v2/users/{id}
	gen.GET("/v2/users/{id}").
		Summary("Get user by ID").
		Description("Returns a single user - Version 2").
		Tags("v2").
		PathParam("id", "User ID", true).
		Response(200, ResponseWrapper{Data: User{}}, "Success").
		Response(404, nil, "Not Found").
		Handler(getUserByIDV2).
		Build()

	// PUT /api/v2/users/{id}
	gen.PUT("/v2/users/{id}").
		Summary("Update user").
		Description("Updates an existing user - Version 2").
		Tags("v2").
		PathParam("id", "User ID", true).
		Body(CreateUserRequest{}, "User data").
		Response(200, User{}, "Updated").
		Response(400, nil, "Bad Request").
		Response(404, nil, "Not Found").
		Handler(updateUserV2).
		Build()

	// DELETE /api/v2/users/{id}
	gen.DELETE("/v2/users/{id}").
		Summary("Delete user").
		Description("Deletes a user - Version 2").
		Tags("v2").
		PathParam("id", "User ID", true).
		Response(204, nil, "No Content").
		Response(404, nil, "Not Found").
		Handler(deleteUserV2).
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

	fmt.Println("Server: http://localhost:8080")
	fmt.Println("Swagger: http://localhost:8080/swagger/index.html")

	r.Run(":8080")
}
