package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     "witches",
	Version: "v1.1.1",
	Short:   "Witches CLI - Go project generator",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	// Scaffold
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(databaseCmd)
	databaseCmd.AddCommand(generateCmd)

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(linkCmd)
	initCmd.AddCommand(initCaptainCmd)
	initCmd.AddCommand(initMemberCmd)

	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateDropCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateVersionCmd)
	migrateCmd.AddCommand(migrateForceCmd)
}
