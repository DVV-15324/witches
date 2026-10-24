package cmd

import (
	"fmt"

	cmd_migrate "github.com/DVV-15324/witches/cmd/migrate"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use: "migrate",
}

func getMigrateConfig(pathArg string) (*utils.Config, string, error) {
	cfg := utils.Load()

	switch cfg.Env {
	case "development", "test":
	case "production":
		return nil, "", fmt.Errorf("migrate command is not allowed in production")
	default:
		return nil, "", fmt.Errorf("invalid ENV: %s", cfg.Env)
	}

	// Nếu pathArg rỗng, hàm GetMigrationsURL của cậu sẽ tự động dùng "migrate/migrations"
	migratePath := utils.GetMigrationsURL(pathArg)
	return cfg, migratePath, nil
}

var migrateUpCmd = &cobra.Command{
	Use: "up [path]",
	RunE: func(cmd *cobra.Command, args []string) error {
		pathArg := ""
		if len(args) > 0 {
			pathArg = args[0] // Lấy args[0] thay vì args[1]
		}
		cfg, migratePath, err := getMigrateConfig(pathArg)
		if err != nil {
			return err
		}
		return cmd_migrate.WitchesMigrateUp(cfg, migratePath)
	},
}

var migrateDownCmd = &cobra.Command{
	Use: "down [path]",
	RunE: func(cmd *cobra.Command, args []string) error {
		pathArg := ""
		if len(args) > 0 {
			pathArg = args[0]
		}
		cfg, migratePath, err := getMigrateConfig(pathArg)
		if err != nil {
			return err
		}
		return cmd_migrate.WitchesMigrateDown(cfg, migratePath)
	},
}

var migrateDropCmd = &cobra.Command{
	Use: "drop [path]",
	RunE: func(cmd *cobra.Command, args []string) error {
		pathArg := ""
		if len(args) > 0 {
			pathArg = args[0]
		}
		cfg, migratePath, err := getMigrateConfig(pathArg)
		if err != nil {
			return err
		}
		return cmd_migrate.WitchesMigrateDrop(cfg, migratePath)
	},
}

var migrateVersionCmd = &cobra.Command{
	Use: "version [path]",
	RunE: func(cmd *cobra.Command, args []string) error {
		pathArg := ""
		if len(args) > 0 {
			pathArg = args[0]
		}
		cfg, migratePath, err := getMigrateConfig(pathArg)
		if err != nil {
			return err
		}
		return cmd_migrate.WitchesMigrateVersion(cfg, migratePath)
	},
}

var migrateForceCmd = &cobra.Command{
	Use: "force <version> [path]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("cú pháp đúng: witches migrate force <version> [path]")
		}
		version := args[0] // args[0] là version (VD: "1")

		pathArg := ""
		if len(args) > 1 {
			pathArg = args[1] // args[1] là path (nếu có)
		}

		cfg, migratePath, err := getMigrateConfig(pathArg)
		if err != nil {
			return err
		}
		return cmd_migrate.WitchesMigrateForce(cfg, migratePath, version)
	},
}

func init() {
	migrateCmd.AddCommand(migrateDropCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateForceCmd)
	migrateCmd.AddCommand(migrateVersionCmd)
	rootCmd.AddCommand(migrateCmd)
}
