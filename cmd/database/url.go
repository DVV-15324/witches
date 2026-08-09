package database

import (
	"fmt"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"log"
	"os"
	"path/filepath"
)

func WitchesDBURL(DB_DRIVER string, config *utils.Config) {
	currentPath := utils.GetCurrentPath()
	envPath := filepath.Join(currentPath, "witches.env")
	file, err := os.OpenFile(envPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Error: update witches.env: %v", err)
	}
	defer file.Close()

	DB_URL := GenerateDBURL(DB_DRIVER,
		config.DBUser, config.DBPassword,
		config.DBHost, config.DBPort,
		config.DBName)

	content := utils.CreateContentRefreshUsed(config.Port,
		config.MetrictPort, config.Host, DB_DRIVER,
		config.DBUser, config.DBHost,
		config.DBPort, config.DBName,
		config.DBPassword, DB_URL,
		config.RedisHost, config.RedisPort,
		config.RedisPassword, config.AccessTokenTTL,
		config.RefreshTokenTTL, config.SessionTTL,
		config.IdleTimeout, config.RevokedTTL,
		config.RateLimitPeriod, config.RateLimitMax)

	if _, err := file.WriteString(content); err != nil {
		log.Fatalf("Error: write to witches.env: %v", err)
	}
	log.Printf("Info: successfully updated witches.env")
}

func GenerateDBURL(DB_DRIVER string,
	DB_USER,
	DB_PASSWORD, DB_HOST,
	DB_PORT, DB_NAME string) string {

	switch DB_DRIVER {
	case "mysql":
		DB_URL := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			DB_USER,
			DB_PASSWORD,
			DB_HOST,
			DB_PORT,
			DB_NAME,
		)
		return DB_URL
	case "postgresql", "postgres":
		DB_URL := fmt.Sprintf(
			"%s:%s@%s:%s/%s?sslmode=disable",
			DB_USER,
			DB_PASSWORD,
			DB_HOST,
			DB_PORT,
			DB_NAME,
		)
		return DB_URL
	case "mssql", "sqlserver":
		DB_URL := fmt.Sprintf(
			"%s:%s@%s:%s?database=%s&encrypt=disable",
			DB_USER,
			DB_PASSWORD,
			DB_HOST,
			DB_PORT,
			DB_NAME,
		)
		return DB_URL
	default:
		log.Fatalf("Error: unsupported database: %s. supported : mysql, postgresql, postgres, mssql, sqlserver", DB_DRIVER)
	}
	return ""
}
