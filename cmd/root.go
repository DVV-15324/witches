package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "witches",
	Short: "Backend Golang nhanh và có khả năng mở rộng",
	Long: `Witches API được xây dựng bằng Go, được thiết kế để đạt hiệu suất cao,
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
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(
		&db,
		"db",
		"",
		"database type",
	)

	//install tool
	rootCmd.AddCommand(installCmd)

	//run test
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(
		&runAll,
		"all",
		"a",
		false,
		"Chay voi easyjson init && swag init",
	)

	rootCmd.AddCommand(databaseCmd)
	// Khoi tao db docker

	rootCmd.AddCommand(migrateCmd)

	databaseCmd.AddCommand(databaseDockerUpCmd)
	databaseCmd.AddCommand(databaseDockerDownCmd)

	migrateCmd.AddCommand(migrateDropCmd)

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)

	migrateCmd.AddCommand(migrateVersionCmd)
	migrateCmd.AddCommand(migrateForceCmd)
}
