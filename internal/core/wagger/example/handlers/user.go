package handlers

// import (
// 	"core-v/internal/core/wagger/models"
// 	"github.com/gin-gonic/gin"
// 	"net/http"
// )

// // In-memory data store (giả lập)
// var usersDB = []models.User{
// 	{ID: "1", Username: "alice", Email: "alice@example.com", Age: 30},
// 	{ID: "2", Username: "bob", Email: "bob@example.com", Age: 25},
// }

// // GetAllUsers retrieves all users with optional filtering
// // @Summary      Get all users
// // @Description  Returns a list of all users, optionally filtered by username
// // @Tags         users
// // @Accept       json
// // @Produce      json
// // @Param        username query string false "Filter by username (partial match)"
// // @Param        limit query int false "Maximum number of results" default(10) minimum(1) maximum(100)
// // @Param        offset query int false "Number of records to skip" default(0) minimum(0)
// // @Success      200 {array} models.User "List of users"
// // @Failure      500 {object} models.ErrorResponse "Internal server error"
// // @Router       /users [get]
// func GetAllUsers(c *gin.Context) {
// 	usernameFilter := c.Query("username")
// 	// limit := c.DefaultQuery("limit", "10")
// 	// offset := c.DefaultQuery("offset", "0")

// 	// Giả lập filter và pagination (có thể parse limit/offset sang int)
// 	filteredUsers := usersDB
// 	if usernameFilter != "" {
// 		temp := []models.User{}
// 		for _, u := range usersDB {
// 			if contains(u.Username, usernameFilter) {
// 				temp = append(temp, u)
// 			}
// 		}
// 		filteredUsers = temp
// 	}

// 	c.JSON(http.StatusOK, filteredUsers)
// }

// // CreateUser handles user creation
// // @Summary      Create a new user
// // @Description  Creates a new user in the system with provided information
// // @Tags         users
// // @Accept       json
// // @Produce      json
// // @Param        request body models.CreateUserRequest true "User creation payload"
// // @Success      201 {object} models.User "User created successfully"
// // @Failure      400 {object} models.ErrorResponse "Invalid input"
// // @Failure      409 {object} models.ErrorResponse "User already exists"
// // @Failure      500 {object} models.ErrorResponse "Internal server error"
// // @Router       /users [post]
// func CreateUser(c *gin.Context) {
// 	var req models.CreateUserRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, models.ErrorResponse{
// 			Code:    400,
// 			Message: "Invalid request payload",
// 		})
// 		return
// 	}

// 	// Logic tạo user (giả lập)
// 	user := models.User{
// 		ID:       "12345",
// 		Username: req.Username,
// 		Email:    req.Email,
// 		Age:      req.Age,
// 	}

// 	c.JSON(http.StatusCreated, user)
// }

// // GetUser retrieves a user by ID
// // @Summary      Get user by ID
// // @Description  Returns detailed information of a specific user
// // @Tags         users
// // @Accept       json
// // @Produce      json
// // @Param        id path string true "User ID"
// // @Success      200 {object} models.User "User found"
// // @Failure      400 {object} models.ErrorResponse "Invalid ID format"
// // @Failure      404 {object} models.ErrorResponse "User not found"
// // @Router       /users/{id} [get]
// func GetUser(c *gin.Context) {
// 	id := c.Param("id")
// 	if id == "" {
// 		c.JSON(http.StatusBadRequest, models.ErrorResponse{
// 			Code:    400,
// 			Message: "ID is required",
// 		})
// 		return
// 	}

// 	// Logic lấy user (giả lập)
// 	for _, u := range usersDB {
// 		if u.ID == id {
// 			c.JSON(http.StatusOK, u)
// 			return
// 		}
// 	}

// 	c.JSON(http.StatusNotFound, models.ErrorResponse{
// 		Code:    404,
// 		Message: "User not found",
// 	})
// }

// // UpdateUser updates existing user information
// // @Summary      Update user
// // @Description  Updates user information by ID
// // @Tags         users
// // @Accept       json
// // @Produce      json
// // @Param        id path string true "User ID"
// // @Param        request body models.CreateUserRequest true "Updated user data"
// // @Success      200 {object} models.User "User updated successfully"
// // @Failure      400 {object} models.ErrorResponse "Invalid input"
// // @Failure      404 {object} models.ErrorResponse "User not found"
// // @Router       /users/{id} [put]
// func UpdateUser(c *gin.Context) {
// 	id := c.Param("id")
// 	var req models.CreateUserRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, models.ErrorResponse{
// 			Code:    400,
// 			Message: "Invalid request payload",
// 		})
// 		return
// 	}

// 	user := models.User{
// 		ID:       id,
// 		Username: req.Username,
// 		Email:    req.Email,
// 		Age:      req.Age,
// 	}

// 	c.JSON(http.StatusOK, user)
// }

// // DeleteUser removes a user from the system
// // @Summary      Delete user
// // @Description  Deletes a user by ID
// // @Tags         users
// // @Accept       json
// // @Produce      json
// // @Param        id path string true "User ID"
// // @Success      204 "No content"
// // @Failure      400 {object} models.ErrorResponse "Invalid ID format"
// // @Failure      404 {object} models.ErrorResponse "User not found"
// // @Router       /users/{id} [delete]
// func DeleteUser(c *gin.Context) {
// 	id := c.Param("id")
// 	if id == "" {
// 		c.JSON(http.StatusBadRequest, models.ErrorResponse{
// 			Code:    400,
// 			Message: "ID is required",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusNoContent, nil)
// }

// // UploadAvatar handles file upload for user avatar
// // @Summary      Upload user avatar
// // @Description  Uploads an avatar image for a specific user
// // @Tags         users
// // @Accept       mpfd
// // @Produce      json
// // @Param        id path string true "User ID"
// // @Param        avatar formData file true "Avatar image file"
// // @Success      200 {object} map[string]string "Avatar uploaded successfully"
// // @Failure      400 {object} models.ErrorResponse "Invalid file or user ID"
// // @Router       /users/{id}/avatar [post]
// func UploadAvatar(c *gin.Context) {
// 	id := c.Param("id")
// 	file, err := c.FormFile("avatar")
// 	if err != nil || id == "" {
// 		c.JSON(http.StatusBadRequest, models.ErrorResponse{
// 			Code:    400,
// 			Message: "Invalid parameters",
// 		})
// 		return
// 	}

// 	// Logic upload file
// 	c.JSON(http.StatusOK, gin.H{
// 		"message":  "Avatar uploaded successfully",
// 		"user_id":  id,
// 		"filename": file.Filename,
// 	})
// }

// // Helper function for partial match (giả lập)
// func contains(s, substr string) bool {
// 	return len(substr) == 0 || (len(s) >= len(substr) && (s == substr || len(s) > 0))
// }
