package cmd

import (
	"log"
	"os"
	"os/exec"
)

func WitchesInit() {
	cmd := exec.Command("make", "migrate-init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat("witches.env"); os.IsNotExist(err) {
		file, err := os.Create("witches.env")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		file.WriteString("DB_URL=LINK_CONNECT_TO_DATABASE\n")
		file.WriteString("DB_PROFILE=TYPE_DATABASE")
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

	cmdMigrate := exec.Command("make", "migrate-init")
	cmdMigrate.Stdout = os.Stdout
	cmdMigrate.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
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
