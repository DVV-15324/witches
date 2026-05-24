package cmd

import (
	utils "core-v/cmd/cmd_utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
)

func WitchesDatabaseDockerUp(DB_PROFILE string) {
	cwd, _ := os.Getwd()
	envPath := filepath.Join(cwd, "witches.env")
	// Load .env file
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Get environment variables with defaults
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	fmt.Println(DB_PASSWORD)
	DATABASE := os.Getenv("DATABASE")

	if DB_PASSWORD == "" {
		log.Fatal("DB_PASSWORD is not set in environment")
	}

	if DATABASE == "" {
		log.Fatal("DATABASE is not set in environment")
	}

	var DB_USER string
	var DB_URL string

	// Validate and set database-specific configurations
	switch DB_PROFILE {
	case "mysql":
		DB_USER = "root"
		DB_URL = fmt.Sprintf(
			"mysql://%s:%s@tcp(localhost:1502)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			DB_USER,
			DB_PASSWORD,
			DATABASE,
		)

	case "postgres":
		DB_USER = "postgres"
		DB_URL = fmt.Sprintf(
			"postgres://%s:%s@localhost:1501/%s?sslmode=disable",
			DB_USER,
			DB_PASSWORD,
			DATABASE,
		)

	case "mssql":
		DB_USER = "sa"
		DB_URL = fmt.Sprintf(
			"sqlserver://%s:%s@localhost:1503?database=%s&encrypt=disable",
			DB_USER,
			DB_PASSWORD,
			DATABASE,
		)

	default:
		log.Fatalf("Unsupported database profile: %s. Supported profiles: mysql, postgres, mssql", DB_PROFILE)
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

	content := fmt.Sprintf("DB_PASSWORD=%s\nDATABASE=%s\nDB_PROFILE=%s\nDB_USER=%s\nDB_URL=%s\n",
		DB_PASSWORD, DATABASE, DB_PROFILE, DB_USER, DB_URL)

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
		"internal",
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
	log.Printf("Starting Docker Compose for profile: %s", DB_PROFILE)

	cmd := exec.Command(
		"docker",
		"compose",
		"--env-file", envPath, // Thêm dòng này
		"--profile", DB_PROFILE,
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

	log.Printf("Successfully started database containers for %s profile", DB_PROFILE)
}
