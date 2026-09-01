package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	cmd := rootCmd
	assert.Equal(t, "witches", cmd.Use)
	assert.Contains(t, cmd.Version, "v1.1.1")
	assert.NotNil(t, cmd.Run)
}

func TestExecute(t *testing.T) {
	// Test with no args
	// This is tricky to test because it calls os.Exit
	// We can test the command exists
	assert.NotNil(t, rootCmd)
}

func TestInitCommands(t *testing.T) {
	// Kiểm tra các command đã được đăng ký
	commands := rootCmd.Commands()
	assert.Greater(t, len(commands), 0)

	commandNames := make([]string, len(commands))
	for i, cmd := range commands {
		commandNames[i] = cmd.Use
	}

	assert.Contains(t, commandNames, "create")
	assert.Contains(t, commandNames, "database")
	assert.Contains(t, commandNames, "install")
	assert.Contains(t, commandNames, "init")
	assert.Contains(t, commandNames, "add")
	assert.Contains(t, commandNames, "run")
	assert.Contains(t, commandNames, "migrate")
}

// Sửa: dùng *cobra.Command thay vì Command
func TestDatabaseSubCommands(t *testing.T) {
	var dbCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "database" {
			dbCmd = cmd
			break
		}
	}

	assert.NotNil(t, dbCmd)

	subCommands := dbCmd.Commands()
	assert.Greater(t, len(subCommands), 0)

	subNames := make([]string, len(subCommands))
	for i, cmd := range subCommands {
		subNames[i] = cmd.Use
	}
	assert.Contains(t, subNames, "generate")
}

// Sửa: dùng *cobra.Command thay vì Command
func TestMigrateSubCommands(t *testing.T) {
	var migrateCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "migrate" {
			migrateCmd = cmd
			break
		}
	}

	assert.NotNil(t, migrateCmd)

	subCommands := migrateCmd.Commands()
	assert.Greater(t, len(subCommands), 0)

	subNames := make([]string, len(subCommands))
	for i, cmd := range subCommands {
		subNames[i] = cmd.Use
	}

	expected := []string{"drop", "up", "down", "force", "version"}
	for _, exp := range expected {
		assert.Contains(t, subNames, exp)
	}
}

// Test Execute with help command
func TestExecute_Help(t *testing.T) {
	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
}

// Test version command
func TestVersion(t *testing.T) {
	rootCmd.SetArgs([]string{"--version"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
}
