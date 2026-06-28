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

// func WitchesRunAll() {
// 	// Chay go run .
// 	cmd := exec.Command("go", "run", ".")
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	if err := cmd.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }
