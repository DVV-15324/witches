package cmd

import (
	"example/cmd/server/config"
	routes "example/cmd/server/routers"
	"example/internal/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "root",
	Short: "Bắt đầu khởi động phần mềm",
	Run: func(cmd *cobra.Command, args []string) {
		g := gin.Default()
		config := config.Load()
		g.Use(middleware.Cors())
		routes.Start(g)
		g.Run(fmt.Sprintf(":%s", config.Port))
	},
}

func GetExcute() *cobra.Command {
	return root
}
