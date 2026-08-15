package utils

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type Config interface {
	GetModuleName() string
}

func RenderTemplate(templateFS embed.FS, baseDir, destFile, tmplFile string, config Config) {
	// Đọc trực tiếp từ embed.FS
	tmpl, err := template.ParseFS(templateFS, tmplFile)
	if err != nil {
		log.Fatalf("Error: parse template %s: %v", tmplFile, err)
	}

	fullPath := filepath.Join(baseDir, destFile)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		log.Fatalf("Error: create folder for %s: %v", fullPath, err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		log.Fatalf("Error: create file %s: %v", fullPath, err)
	}
	defer file.Close()

	// Execute template trực tiếp
	if err := tmpl.Execute(file, config); err != nil {
		log.Fatalf("Error: execute template %s: %v", tmplFile, err)
	}
}
