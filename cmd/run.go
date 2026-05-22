package cmd

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func WitchesInit(project string) {
	projectPath := filepath.Join(".", project)

	err := os.MkdirAll(projectPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(filepath.Join(projectPath, "migrate", "migrations"), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	envPath := filepath.Join(projectPath, "witches.env")
	// Create witches.env
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		file, err := os.Create(envPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = file.WriteString(
			"DB_URL=LINK_CONNECT_TO_DATABASE\n" +
				"DB_PROFILE=TYPE_DATABASE\n",
		)
		if err != nil {
			log.Fatal(err)
		}
	}
	mainPath := filepath.Join(projectPath, "main.go")
	// Create main.go
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

func WitchesInstall() {
	cmd := exec.Command("make", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// cmdMigrate := exec.Command("make", "migrate-init")
	// cmdMigrate.Stdout = os.Stdout
	// cmdMigrate.Stderr = os.Stderr
	// err = cmd.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func WitchesRun() {
	cmd := exec.Command("make", "run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
