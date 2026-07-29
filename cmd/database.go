package cmd

import (
	"log"
	"os"

	cmd_database "github.com/DVV-15324/witches/cmd/database"
	godotenv "github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
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

		cmd_database.WitchesDatabase(os.Getenv("DB_DRIVER"))
	},
}
