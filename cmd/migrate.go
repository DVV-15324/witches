package cmd

import (
	cmd_migrate "github.com/DVV-15324/witches/cmd/migrate"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Run: func(cmd *cobra.Command, args []string) {},
}

var migrateDropCmd = &cobra.Command{
	Use: "drop",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		migratePath := utils.GetMigrationsURL()
		cmd_migrate.WitchesMigrateDrop(cfg.DBUrl, cfg.DBDriver, migratePath)
	},
}

var migrateUpCmd = &cobra.Command{
	Use: "up",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		migratePath := utils.GetMigrationsURL()
		cmd_migrate.WitchesMigrateUp(cfg.DBUrl, cfg.DBDriver, migratePath)
	},
}

var migrateDownCmd = &cobra.Command{
	Use: "down",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		migratePath := utils.GetMigrationsURL()
		cmd_migrate.WitchesMigrateDown(cfg.DBUrl, cfg.DBDriver, migratePath)
	},
}

var migrateVersionCmd = &cobra.Command{
	Use: "force",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		migratePath := utils.GetMigrationsURL()
		cmd_migrate.WitchesMigrateForce(cfg.DBUrl, cfg.DBDriver, migratePath, args[0])
	},
}
var migrateForceCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		migratePath := utils.GetMigrationsURL()
		cmd_migrate.WitchesMigrateVersion(cfg.DBUrl, cfg.DBDriver, migratePath)
	},
}
