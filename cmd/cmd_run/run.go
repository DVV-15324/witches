package cmd

import (
	"log"
	"os"
	"os/exec"
)

func WitchesRun() {
	cmd := exec.Command("go", "run", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func WitchesRunAll() {
	// Chay easyjson
	cmd_easyjson := exec.Command("easyjson", "init")
	cmd_easyjson.Stdout = os.Stdout
	cmd_easyjson.Stderr = os.Stderr
	if err := cmd_easyjson.Run(); err != nil {
		log.Fatalf("easyjson failed: %v", err)
	}

	// Chay swag
	cmd_swag := exec.Command("swag", "init")
	cmd_swag.Stdout = os.Stdout
	cmd_swag.Stderr = os.Stderr
	if err := cmd_swag.Run(); err != nil {
		log.Fatalf("swag failed: %v", err)
	}

	// Chay go run .
	cmd := exec.Command("go", "run", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
