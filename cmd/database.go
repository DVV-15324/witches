package cmd

import (
	cmd_database "github.com/DVV-15324/witches/cmd/database"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use: "database",
	Run: func(cmd *cobra.Command, args []string) {},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate DB_URL from witches.env",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := utils.PreloadNotDBURL()
		utils.Validate(cfg)
		cmd_database.WitchesDBURL(cfg.DBDriver, cfg)
	},
}
