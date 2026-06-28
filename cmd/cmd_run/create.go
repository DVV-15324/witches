package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// project: chương trình của bạn
// DB_DRIVER: loại database: mysql | mssql | postgres
func WitchesCreate(project string, DB_DRIVER string) {
	projectPath := filepath.Join(".", project)
	err := os.MkdirAll(projectPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	// Tạo .env
	envPath := filepath.Join(projectPath, "witches.env")
	file, err := os.OpenFile(envPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = file.WriteString(
		"APP_PORT=YOUR_PORT_APP\n" +
			"DB_PASSWORD=YOUR_PASSWORD\n" +
			"DB_NAME=YOUR_DATABASE\n" +
			"DB_HOST=DB_HOST\n" +
			"DB_PORT=YOUR_PORT_APP\n" +
			"REDIS_PORT=REDIS_PORT\n" +
			fmt.Sprintf("DB_DRIVER=%s\n", DB_DRIVER),
	)
	if err != nil {
		log.Fatal(err)
	}

}
