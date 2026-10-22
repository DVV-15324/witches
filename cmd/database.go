package cmd

import (
	"log"
	"os"
	"path/filepath"

	cmd_database "github.com/DVV-15324/witches/cmd/database"
	cmd_utils "github.com/DVV-15324/witches/cmd/utils"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use: "database",
	Run: func(cmd *cobra.Command, args []string) {},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate DB_URL from witches.env",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentPath := cmd_utils.GetCurrentPath()
		envPath := filepath.Join(currentPath, "witches.env")
		if err := godotenv.Load(envPath); err != nil {
			log.Println("Error: not found witches.env, using default values")
			os.Exit(1)
		}

		cfg := utils.PreloadNotDBURL()
		if err := cmd_database.WitchesDBURL(cfg.DBDriver, cfg); err != nil {
			return err
		}

		return nil
	},
}
