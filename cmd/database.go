package cmd

import (
	//"fmt"

	cmd_database "core-v/cmd/cmd_database"
	godotenv "github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Quản lý database",
	Long: `Container database:
			- postgres
			- mySQL
			- mSSQL
			- redis`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	},
}

var databaseDockerUpCmd = &cobra.Command{
	Use:   "docker-up",
	Short: "Tao database",
	Long: `Tạo container database:
			- PostgreSQL
			- MySQL
			- MSSQL
			Defalt:
			- Migrate
			- Redis`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		cmd_database.WitchesDatabaseDockerUp(os.Getenv("DB_PROFILE"))
	},
}

var databaseDockerDownCmd = &cobra.Command{
	Use:   "docker-down",
	Short: "Tao database",
	Long: `Tạo container database:
			- PostgreSQL
			- MySQL
			- MSSQL
			Defalt:
			- Migrate
			- Redis`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		cmd_database.WitchesDatabaseDockerDown()
	},
}
