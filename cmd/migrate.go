package cmd

import (
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		switch cfg.Env {
		case "development", "test":
		case "production":
			return fmt.Errorf(
				"migrate command is not allowed in production",
			)
		default:
			return fmt.Errorf("invalid ENV: %s", cfg.Env)
		}
		migratePath := utils.GetMigrationsURL(args[1])
		cmd_migrate.WitchesMigrateDrop(cfg.DBUrl, cfg.DBDriver, migratePath)
		return nil
	},
}

var migrateUpCmd = &cobra.Command{
	Use: "up",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)

		switch cfg.Env {
		case "development", "test":
		case "production":
			return fmt.Errorf(
				"migrate command is not allowed in production",
			)
		default:
			return fmt.Errorf("invalid ENV: %s", cfg.Env)
		}
		migratePath := utils.GetMigrationsURL(args[1])
		cmd_migrate.WitchesMigrateUp(cfg.DBUrl, cfg.DBDriver, migratePath)
		return nil
	},
}

var migrateDownCmd = &cobra.Command{
	Use: "down",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		switch cfg.Env {
		case "development", "test":
		case "production":
			return fmt.Errorf(
				"migrate command is not allowed in production",
			)
		default:
			return fmt.Errorf("invalid ENV: %s", cfg.Env)
		}
		migratePath := utils.GetMigrationsURL(args[1])
		cmd_migrate.WitchesMigrateDown(cfg.DBUrl, cfg.DBDriver, migratePath)
		return nil
	},
}

var migrateVersionCmd = &cobra.Command{
	Use: "force",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		switch cfg.Env {
		case "development", "test":
		case "production":
			return fmt.Errorf(
				"migrate command is not allowed in production",
			)
		default:
			return fmt.Errorf("invalid ENV: %s", cfg.Env)
		}
		migratePath := utils.GetMigrationsURL(args[1])
		cmd_migrate.WitchesMigrateForce(cfg.DBUrl, cfg.DBDriver, migratePath, args[0])
		return nil
	},
}
var migrateForceCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := utils.PreloadNotDBURL()
		utils.LoadDbUrl(cfg)
		migratePath := utils.GetMigrationsURL(args[1])
		cmd_migrate.WitchesMigrateVersion(cfg.DBUrl, cfg.DBDriver, migratePath)
	},
}
