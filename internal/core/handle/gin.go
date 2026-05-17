package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GinInit() {
	r := gin.Default()

	r.GET("/init-1", GinInit1)
	r.GET("/init-2", GinInit2)
	r.GET("/init-3", GinInit3)
	r.GET("/init-4/:name", GinInit4)

	r.Run(":8080")

}

func GinInit1(c *gin.Context) {
	j := []int{1, 4, 5}
	//c.BindJSON(j)
	c.JSON(http.StatusOK, gin.H{"data": j})

}

type Profile struct {
	Name string `json:"name"`
}

func GinInit2(c *gin.Context) {
	p := &Profile{}
	err := c.BindJSON(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loi data"})
	}
	c.JSON(http.StatusOK, p)

}

func GinInit3(c *gin.Context) {
	name := c.Query("name")
	p := &Profile{
		Name: name,
	}
	err := c.BindJSON(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loi data"})
	}
	c.JSON(http.StatusOK, gin.H{"data": name})

}

func GinInit4(c *gin.Context) {
	name := c.Param("name")
	p := &Profile{
		Name: name,
	}
	err := c.BindJSON(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loi data"})
	}
	c.JSON(http.StatusOK, gin.H{"data": name})

}
