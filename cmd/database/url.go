package database

import (
	"fmt"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"log"
	"os"
	"path/filepath"
)

func WitchesDBURL(DB_DRIVER string, config *utils.Config) error {
	currentPath := utils.GetCurrentPath()
	envPath := filepath.Join(currentPath, "witches.env")
	file, err := os.OpenFile(envPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("update witches.env: %v", err)
	}
	defer file.Close()

	DB_URL, err := GenerateDBURL(DB_DRIVER,
		config.DBUser, config.DBPassword,
		config.DBHost, config.DBName, config.DBPort)
	if err != nil {
		return err
	}
	config.DBUrl = DB_URL
	content := utils.CreateContentRefreshUsed(DB_URL, config)

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("write to witches.env: %v", err)
	}
	log.Printf("Info: successfully updated witches.env")
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
