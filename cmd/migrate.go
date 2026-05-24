package cmd

import (
	cmd_migrate "core-v/cmd/cmd_migrate"
	"fmt"
	godotenv "github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Cai dat Migrate",
	Long:  `Quan ly cac tac vu thuc hien trong Database`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

	},
}

var migrateDockerUpCmd = &cobra.Command{
	Use:   "docker-up",
	Short: "Up Migrate",
	Long:  `Quan ly cac tac vu thuc hien trong Database`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load("witches.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		db_url := os.Getenv("DB_URL")
		cmd_migrate.WitchesMigrateDockerUp(db_url)
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
		cmd_migrate.WitchesMigrateUp(db_url)
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
		cmd_migrate.WitchesMigrateDown(db_url)
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
		cmd_migrate.WitchesMigrateVersion(db_url)
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
		if len(args) < 1 {
			fmt.Println("missing project")
			return
		}
		cmd_migrate.WitchesMigrateForce(db_url, args[0])
	},
}
