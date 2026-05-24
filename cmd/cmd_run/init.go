package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// project: chương trình của bạn
// DB_PROFILE: loại database: mysql | mssql | postgres
func WitchesInit(project string, DB_PROFILE string) {
	projectPath := filepath.Join(".", project)
	err := os.MkdirAll(projectPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	// Tạo folder migrate/migrations
	// Với os.ModePerm: cho phép quyền rwxrwxrwx [r: Đọc, x: Thực thi, w: Viết(Tạo file | folder)]
	err = os.MkdirAll(filepath.Join(projectPath, "migrate", "migrations"), os.ModePerm)
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
		"DB_PASSWORD=YOUR_PASSWORD\n" +
			"DATABASE=YOUR_DATABASE\n" +
			fmt.Sprintf("DB_PROFILE=%s\n", DB_PROFILE),
	)
	if err != nil {
		log.Fatal(err)
	}

	mainPath := filepath.Join(projectPath, "main.go")
	// Tạo main.go
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		file, err := os.Create(mainPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = file.WriteString(`package main

import "fmt"

func main() {
	fmt.Println("Hello this is witches")
}
`)
		if err != nil {
			log.Fatal(err)
		}
	}
	cmdModInit := exec.Command("go", "mod", "init", project)
	cmdModInit.Stdout = os.Stdout
	cmdModInit.Stderr = os.Stderr
	cmdModInit.Dir = projectPath
	err = cmdModInit.Run()
	if err != nil {
		log.Fatal(err)
	}
}
