// Đây chỉ là file ví dụ
package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GinServerInit() {
	r := gin.Default()
	r.GET("/init-1", GinInit1)
	r.GET("/init-2", GinInit2)
	r.GET("/init-3", GinInit3)
	r.GET("/init-4/:name", GinInit4)

	r.Run(":8080")

}

// gin.H giup tu dong chuyen hoa du lieu sang json theo kieu key values
func GinInit1(c *gin.Context) {
	j := []int{1, 4, 5}
	c.JSON(http.StatusOK, gin.H{"data": j})

}

type Profile struct {
	Name string `json:"name"`
}

// BinJSON tu convert data reqest sanng JSON
func GinInit2(c *gin.Context) {
	p := &Profile{}
	err := c.BindJSON(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loi data"})
	}
	c.JSON(http.StatusOK, p)

}

// Dung Query param trong rest api
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

// Dung path param trong rest api
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
