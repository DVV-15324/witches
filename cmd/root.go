package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "witches",
	Short: "En: Fast and Scalable Golang Backend Vi: Backend Golang nhanh và có khả năng mở rộng",
	Long: `En: The Witches API is built using Go, designed for a clean architecture and suitable for classic, modern backend development.
Vi: Witches API được xây dựng bằng Go, được thiết kế để kiến trúc gọn gàng và phù hợp với phát triển backend cổ điển, hiện đại.`,
	Version: "v1.0.5",
	Run:     func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	// Scaffold
	rootCmd.AddCommand(createCmd)

	rootCmd.AddCommand(installCmd)

	rootCmd.AddCommand(initCmd)

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(databaseCmd)

	rootCmd.AddCommand(migrateCmd)

	migrateCmd.AddCommand(migrateDropCmd)

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)

	migrateCmd.AddCommand(migrateVersionCmd)
	migrateCmd.AddCommand(migrateForceCmd)
}
