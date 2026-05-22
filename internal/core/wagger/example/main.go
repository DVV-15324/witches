package main

// import (
// 	_ "core-v/docs"
// 	"github.com/gin-gonic/gin"
// 	swaggerFiles "github.com/swaggo/files"
// 	ginSwagger "github.com/swaggo/gin-swagger"
// 	"net/http"
// )

// var dataBase = []*ModelUser{
// 	{Name: "Dinh Viet Vu", Age: 22},
// 	{Name: "Dinh Vu 1", Age: 21},
// 	{Name: "Dinh Vu 2", Age: 21},
// 	{Name: "Dinh Vu 3", Age: 21},
// 	{Name: "Dinh Vu 4", Age: 21},
// 	{Name: "Dinh Vu 5", Age: 21},
// }

// type ModelUser struct {
// 	Name string `json:"name" description:"User Name" maxLength:"10" minLength:"0" binding:"required" example:"Dinh Viet Vu"`
// 	Age  int    `json:"age,omitempty" description:"Age User" maximum:"100" minimum:"0" binding:"required" example:"22"`
// }

// // @Summary Hiển thị hết User
// // @Description Hiển thị tất cả User
// // @Tags user
// // @Accept json
// // @Produce json
// // @Param Authorization header string true "Bearer Token"
// // @Success 200 {array} ModelUser
// // @Failure 400 {string} string "Loi"
// // @Router /user [get]
// func GetAllUser(cxt *gin.Context) {
// 	cxt.JSON(http.StatusOK, dataBase)
// }

// // @Summary Hiển thị User bởi Name
// // @Description Hiển thị User bởi name
// // @Tags user
// // @Accept json
// // @Produce json
// // @Param name path string true "User Name"
// // @Param Authorization header string true "Bearer Token"
// // @Success 200 {object} ModelUser
// // @Failure 400 {string} string "Loi"
// // @Router /user/{name} [get]
// // @Security BearerAuth
// func GetUserByName(cxt *gin.Context) {
// 	n := cxt.Param("name")

// 	for i := 0; i < len(dataBase); i++ {
// 		if dataBase[i].Name == n {
// 			cxt.JSON(http.StatusOK, dataBase[i])
// 			return
// 		}
// 	}

// 	cxt.JSON(http.StatusNotFound, gin.H{
// 		"error": "user not found",
// 	})
// }

// // @Summary Tạo User
// // @Description Tạo mới User
// // @Tags user
// // @Accept json
// // @Produce json
// // @Param user body ModelUser true "User Data"
// // @Param Authorization header string true "Bearer Token"
// // @Success 201 {object} ModelUser
// // @Failure 400 {string} string "Bad Request"
// // @Router /user [post]
// // @Security BearerAuth
// func CreateUser(cxt *gin.Context) {
// 	var newUser ModelUser

// 	if err := cxt.ShouldBindJSON(&newUser); err != nil {
// 		cxt.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	dataBase = append(dataBase, &newUser)

// 	cxt.JSON(http.StatusCreated, newUser)
// }

// // @Summary Cập nhật User bởi Name
// // @Description Cập nhật User theo Name
// // @Tags user
// // @Accept json
// // @Produce json
// // @Param name path string true "User Name"
// // @Param user body ModelUser true "Updated User"
// // @Param Authorization header string true "Bearer Token"
// // @Success 200 {object} ModelUser
// // @Failure 400 {string} string "Loi"
// // @Failure 404 {string} string "User Not Found"
// // @Router /user/{name} [put]
// // @Security BearerAuth
// func UpdateUserByName(cxt *gin.Context) {
// 	name := cxt.Param("name")

// 	var updateUser ModelUser

// 	if err := cxt.ShouldBindJSON(&updateUser); err != nil {
// 		cxt.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	for i := 0; i < len(dataBase); i++ {
// 		if dataBase[i].Name == name {
// 			dataBase[i].Name = updateUser.Name
// 			dataBase[i].Age = updateUser.Age

// 			cxt.JSON(http.StatusOK, dataBase[i])
// 			return
// 		}
// 	}

// 	cxt.JSON(http.StatusNotFound, gin.H{
// 		"error": "user not found",
// 	})
// }

// // @Summary Xóa User bởi Name
// // @Description Xóa User theo Name
// // @Tags user
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param name path string true "User Name"
// // @Param Authorization header string true "Bearer Token"
// // @Success 200 {string} string "Delete Success"
// // @Failure 404 {string} string "User Not Found"
// // @Router /user/{name} [delete]
// func DeleteUserById(cxt *gin.Context) {
// 	name := cxt.Param("name")

// 	for i := 0; i < len(dataBase); i++ {
// 		if dataBase[i].Name == name {

// 			dataBase = append(dataBase[:i], dataBase[i+1:]...)

// 			cxt.JSON(http.StatusOK, gin.H{
// 				"message": "delete success",
// 			})
// 			return
// 		}
// 	}

// 	cxt.JSON(http.StatusNotFound, gin.H{
// 		"error": "user not found",
// 	})
// }

// // @title API Dinh Viet Vu
// // @version 1.0
// // @description API cua Dinh Viet Vu
// // @host localhost:8080
// // @BasePath /api/v1
// func main() {

// 	r := gin.Default()
// 	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
// 	api := r.Group("/api/v1")
// 	{
// 		api.GET("/user", GetAllUser)
// 		api.GET("/user/:name", GetUserByName)
// 		api.POST("/user", CreateUser)
// 		api.PUT("/user/:name", UpdateUserByName)
// 		api.DELETE("/user/:name", DeleteUserById)
// 	}

// 	r.Run(":8080")
// }
