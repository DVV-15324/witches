package cmd

import (
	"fmt"
	utils "github.com/DVV-15324/witches/cmd/cmd_utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
)

func WitchesDatabaseDockerUp(DB_DRIVER string) {
	cwd, _ := os.Getwd()
	fmt.Println(cwd, "hello")
	envPath := filepath.Join(cwd, "witches.env")
	// Load .env file
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Get environment variables with defaults
	APP_PORT := os.Getenv("APP_PORT")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")
	DB_HOST := os.Getenv("DB_HOST")
	REDIS_PORT := os.Getenv("REDIS_PORT")

	if APP_PORT == "" {
		log.Fatal("APP_PORT is not set in environment")
	}

	if DB_PASSWORD == "" {
		log.Fatal("DB_PASSWORD is not set in environment")
	}

	if DB_NAME == "" {
		log.Fatal("DB_NAME is not set in environment")
	}

	if DB_PORT == "" {
		log.Fatal("DB_PORT is not set in environment")
	}

	if DB_HOST == "" {
		log.Fatal("DB_HOST is not set in environment")
	}
	if REDIS_PORT == "" {
		log.Fatal("REDIS_PORT is not set in environment")
	}

	var DB_USER string
	var DB_URL string

	// Validate and set database-specific configurations
	switch DB_DRIVER {
	case "mysql":
		DB_USER = "root"
		DB_URL = fmt.Sprintf(
			"mysql://%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			DB_USER,
			DB_PASSWORD,
			DB_HOST,
			DB_PORT,
			DB_NAME,
		)

	case "postgres":
		DB_USER = "postgres"
		DB_URL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			DB_USER,
			DB_PASSWORD,
			DB_HOST,
			DB_PORT,
			DB_NAME,
		)

	case "mssql":
		DB_USER = "sa"
		DB_URL = fmt.Sprintf(
			"sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
			DB_USER,
			DB_PASSWORD,
			DB_HOST,
			DB_PORT,
			DB_NAME,
		)

	default:
		log.Fatalf("Unsupported database profile: %s. Supported profiles: mysql, postgres, mssql", DB_DRIVER)
	}

	// Create witches.env file
	//envPath := filepath.Join(".", "witches.env")

	file, err := os.OpenFile(
		envPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)
	if err != nil {
		log.Fatalf("Failed to create witches.env: %v", err)
	}
	defer file.Close()
	content := fmt.Sprintf(
		"APP_PORT=%s\n"+
			"DB_PASSWORD=%s\n"+
			"DB_NAME=%s\n"+
			"DB_HOST=%s\n"+
			"DB_PORT=%s\n"+
			"REDIS_PORT=%s\n"+
			"DB_DRIVER=%s\n"+
			"DB_URL=%s\n", APP_PORT, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT, REDIS_PORT, DB_DRIVER, DB_URL,
	)
	if _, err := file.WriteString(content); err != nil {
		log.Fatalf("Failed to write to witches.env: %v", err)

	}

	log.Printf("Successfully created %s", envPath)

	// Get docker compose path
	frameworkPath := utils.GetFrameworkPath()
	if frameworkPath == "" {
		log.Fatal("Framework path is empty. Please check utils.GetFrameworkPath()")
	}

	dockerPath := filepath.Join(
		frameworkPath,
		"pkg",
		"core",
		"database",
	)

	// Check if docker config exists
	if _, err := os.Stat(dockerPath); os.IsNotExist(err) {
		log.Fatalf(
			"Docker config not found at: %s\nPlease ensure witches is installed correctly.\nError: %v",
			dockerPath,
			err,
		)
	}

	// Run docker compose command
	log.Printf("Starting Docker Compose for profile: %s", DB_DRIVER)

	cmd := exec.Command(
		"docker",
		"compose",
		"--env-file", envPath, // Thêm dòng này
		"--profile", DB_DRIVER,
		"--profile", "redis",
		"up",
		"--build",
		"-d",
	)
	cmd.Dir = dockerPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Docker Compose failed: %v", err)
	}

	log.Printf("Successfully started database containers for %s profile", DB_DRIVER)
}
