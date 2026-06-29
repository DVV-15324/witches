package cmd

import (
	cmd_database "github.com/DVV-15324/witches/cmd/cmd_database"
	godotenv "github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// En: Database command group
// Vi: Nhóm lệnh quản lý database

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "En: Database management Vi: Quản lý database",
	Long: `Database containers:
			- PostgreSQL
			- MySQL
			- MSSQL
			Defalt:
			- Redis`,
	Run: func(cmd *cobra.Command, args []string) {
		// En: Load environment variables from witches.env
		// Vi: Kiểm tra biến môi trường từ file witches.env
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error: loading .env file")
		}
	},
}

var databaseDockerUpCmd = &cobra.Command{
	Use:   "docker-up",
	Short: "Deploy database containers Vi: Triển khai database containers",
	Long: `Database containers:
			- PostgreSQL
			- MySQL
			- MSSQL
			Defalt:
			- Redis`,
	Run: func(cmd *cobra.Command, args []string) {
		// En: Load environment variables from witches.env
		// Vi: Kiểm tra biến môi trường từ file witches.env
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error: loading .env file")
		}

		cmd_database.WitchesDatabaseDockerUp(os.Getenv("DB_DRIVER"))
	},
}

var databaseDockerDownCmd = &cobra.Command{
	Use:   "docker-down",
	Short: "Down database containers Vi: Xóa database containers",
	Long: `Database containers:
			- PostgreSQL
			- MySQL
			- MSSQL
			Defalt:
			- Redis`,
	Run: func(cmd *cobra.Command, args []string) {
		// En: Load environment variables from witches.env
		// Vi: Kiểm tra biến môi trường từ file witches.env
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error: loading .env file")
		}

		cmd_database.WitchesDatabaseDockerDown()
	},
}
