package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "witches",
	Short: "Backend Golang nhanh và có khả năng mở rộng",
	Long: `Witches API được xây dựng bằng Go, được thiết kế để
kiến trúc gọn gàng và phù hợp với phát triển backend cổ điển, hiện đại.`,
	Version: "v1.1",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Scaffold
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(
		&db,
		"db",
		"",
		"database type",
	)

	//install tool
	rootCmd.AddCommand(installCmd)

	rootCmd.AddCommand(initCmd)

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(databaseCmd)

	rootCmd.AddCommand(migrateCmd)

	databaseCmd.AddCommand(databaseDockerUpCmd)
	databaseCmd.AddCommand(databaseDockerDownCmd)

	migrateCmd.AddCommand(migrateDropCmd)

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)

	migrateCmd.AddCommand(migrateVersionCmd)
	migrateCmd.AddCommand(migrateForceCmd)
}
