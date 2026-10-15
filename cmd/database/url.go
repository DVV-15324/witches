package database

import (
	"fmt"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var readFileWitchesDBURL = os.ReadFile

func WitchesDBURL(DB_DRIVER string, config *utils.Config) error {
	currentPath := utils.GetCurrentPath()
	envPath := filepath.Join(currentPath, "witches.env")
	content, err := readFileWitchesDBURL(envPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read witches.env: %v", err)
	}
	DB_URL, err := GenerateDBURL(DB_DRIVER,
		config.DBUser, config.DBPassword,
		config.DBHost, config.DBName, config.DBPort)
	if err != nil {
		return err
	}
	config.DBUrl = DB_URL
	if len(content) == 0 {
		content = []byte(utils.CreateContent())
	}
	lines := strings.Split(string(content), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, "DB_URL=") {
			lines[i] = "DB_URL=" + DB_URL
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, "DB_URL="+DB_URL)
	}

	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(envPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("write witches.env: %v", err)
	}
	log.Printf("Info: successfully updated DB_URL in witches.env")
	fmt.Printf("\nNext steps:\n")
	fmt.Print("  witches init\n")
	return nil
}

func GenerateDBURL(DB_DRIVER string,
	DB_USER,
	DB_PASSWORD, DB_HOST,
	DB_NAME string, DB_PORT int64) (string, error) {

	switch DB_DRIVER {
	case "mysql":
		DB_URL := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			DB_USER,
			DB_PASSWORD,
			DB_HOST,
			DB_PORT,
			DB_NAME,
		)
		return DB_URL, nil
	case "postgresql", "postgres":
		DB_URL := fmt.Sprintf(
			"%s:%s@%s:%d/%s?sslmode=disable",
			DB_USER,
			DB_PASSWORD,
			DB_HOST,
			DB_PORT,
			DB_NAME,
		)
		return DB_URL, nil
	case "mssql", "sqlserver":
		DB_URL := fmt.Sprintf(
			"%s:%s@%s:%d?database=%s&encrypt=disable",
			DB_USER,
			DB_PASSWORD,
			DB_HOST,
			DB_PORT,
			DB_NAME,
		)
		return DB_URL, nil
	default:
		return "", fmt.Errorf("unsupported database: %s. supported: mysql, postgresql, postgres, mssql, sqlserver", DB_DRIVER)
	}
}
