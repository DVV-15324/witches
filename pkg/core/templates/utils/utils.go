package utils

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type Config interface {
	GetProjectName() string
}

func RenderTemplate(templateFS embed.FS, baseDir, destFile, tmplFile string, config Config) {
	tmpl, err := template.ParseFS(templateFS, tmplFile)
	if err != nil {
		log.Fatalf("parse template %s: %v", tmplFile, err)
	}
	fullPath := filepath.Join(baseDir, destFile)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		log.Fatalf("create folder for %s: %v", fullPath, err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		log.Fatalf("create file %s: %v", fullPath, err)
	}
	defer func() {
		_ = file.Close()
	}()
	if err := tmpl.Execute(file, config); err != nil {
		log.Fatalf("execute template %s: %v", tmplFile, err)
	}
}
