package cmd

import (
	"fmt"
	"log"
	"os"

	godotenv "github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "witches",
	Short: "Backend Golang nhanh và có khả năng mở rộng",
	Long: `Witches API được xây dựng bằng Go, được thiết kế để đạt hiệu suất cao,
kiến trúc gọn gàng và phù hợp với phát triển backend cổ điển, hiện đại.`,
	Version: "v1.1",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Witches API đang chạy...")
	},
}

var runCmd = &cobra.Command{
	Use:   "start",
	Short: "Cài đặt các phụ thuộc cần thiết",
	Long:  `Cài easyjson && swag && migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		WitchesRun()
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Cài đặt các phụ thuộc cần thiết",
	Long:  `Cài easyjson && swag && migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		WitchesInstall()
	},
}

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
		db_profile := os.Getenv("DB_PROFILE")
		fmt.Println(db_profile)
	},
}

var databaseUpCmd = &cobra.Command{
	Use:   "up",
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
		db_profile := os.Getenv("DB_PROFILE")
		WitchesDatabaseUp(db_profile)
	},
}

var databaseDownCmd = &cobra.Command{
	Use:   "down",
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
		WitchesDatabaseDown()
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Cai dat Migrate",
	Long:  `Quan ly cac tac vu thuc hien trong Database`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		db_url := os.Getenv("DB_URL")
		fmt.Println(db_url)

	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Cai dat Migrate",
	Long:  `Quan ly cac tac vu thuc hien trong Database`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("missing project")
			return
		}

		WitchesInit(args[0])
	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Up Migrate",
	Long:  `Quan ly cac tac vu thuc hien trong Database`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		db_url := os.Getenv("DB_URL")
		WitchesMigrateUp(db_url)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Down Migrate",
	Long:  `Quan ly cac tac vu thuc hien trong Database`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		db_url := os.Getenv("DB_URL")
		WitchesMigrateDown(db_url)
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:   "force",
	Short: "Force Migrate",
	Long:  `Quan ly cac tac vu thuc hien trong Database`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		db_url := os.Getenv("DB_URL")
		WitchesMigrateVersion(db_url)
	},
}
var migrateForceCmd = &cobra.Command{
	Use:   "version",
	Short: "Version Migrate",
	Long:  `Quan ly cac tac vu thuc hien trong Database`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		db_url := os.Getenv("DB_URL")
		WitchesMigrateForce(db_url)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(databaseCmd)
	rootCmd.AddCommand(initCmd)

	databaseCmd.AddCommand(databaseUpCmd)
	databaseCmd.AddCommand(databaseDownCmd)

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)

	migrateCmd.AddCommand(migrateVersionCmd)
	migrateCmd.AddCommand(migrateForceCmd)
}
