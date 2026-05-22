package wagger

// import (
// 	docs "core-v/docs"

// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	swaggerFiles "github.com/swaggo/files"

// 	ginSwagger "github.com/swaggo/gin-swagger"
// )

// // @BasePath /api/v1
// // @Summary ping example
// // @Schemes http
// // @Description do ping
// // @Tags example
// // @Accept json
// // @Produce json
// // @Success 200 {string} Helloworld
// // @Failure 404 {string} string "Not found - error message"
// // @Router /example/helloworld [get]
// func Helloworld(g *gin.Context) {
// 	g.JSON(http.StatusOK, "helloworld")
// }

// func SwagInit() {
// 	r := gin.Default()
// 	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
// 	docs.SwaggerInfo.BasePath = "/"
// 	v1 := r.Group("/")
// 	{
// 		eg := v1.Group("/example")
// 		{

// 			eg.GET("/helloworld", Helloworld)
// 		}
// 	}
// 	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
// 	r.Run(":8080")
// }
