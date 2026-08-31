package template

import (
	"embed"
	"fmt"
	"github.com/DVV-15324/witches/pkg/core/templates/utils"
	"os"
	"path/filepath"
)

//go:embed member/pkg/redis/*.tmpl
//go:embed member/*.tmpl
//go:embed member/cmd/*.tmpl
//go:embed member/cmd/server/config/*.tmpl
//go:embed member/cmd/server/routers/*.tmpl
//go:embed member/cmd/server/core/core.go.tmpl
//go:embed member/internal/shared/utils/*.tmpl
//go:embed member/internal/shared/middleware/*.tmpl
var templateMemberFS embed.FS

type MemberConfig struct {
	ModuleName string
}

func (p MemberConfig) GetModuleName() string {
	return p.ModuleName
}

func CreateMemberGoArc(projectName string, typeDb string) {
	config := CaptainConfig{
		ModuleName: projectName,
	}
	fmt.Printf("Generating project: %s\n", projectName)
	fmt.Println("Creating structure...")

	if err := createMemberProjectStructure(config, typeDb); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Project created successfully!")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  witches install\n")
	fmt.Printf("  witches run\n")
}

func createMemberProjectStructure(config CaptainConfig, typeDb string) error {
	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	baseDir = filepath.Join(baseDir)

	// Tạo thư mục project
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	files := map[string]string{
		// ROOT
		"member/main.go.tmpl":   "main.go",
		"member/README.md.tmpl": "README.md",
		"member/go.mod.tmpl":    "go.mod",

		// CMD
		"member/cmd/root.go.tmpl":                    filepath.Join("cmd", "root.go"),
		"member/cmd/server/config/config.go.tmpl":    filepath.Join("cmd", "server", "config", "config.go"),
		"member/cmd/server/routers/composer.go.tmpl": filepath.Join("cmd", "server", "routers", "composer.go"),
		"member/cmd/server/routers/routers.go.tmpl":  filepath.Join("cmd", "server", "routers", "routers.go"),
		"member/cmd/server/routers/modules.go.tmpl":  filepath.Join("cmd", "server", "routers", "modules.go"),
		"member/cmd/server/core/core.go.tmpl":        filepath.Join("cmd", "server", "core", "core.go"),

		// SHARED - MIDDLEWARE
		"member/internal/shared/middleware/limit.go.tmpl":  filepath.Join("internal", "shared", "middleware", "limit.go"),
		"member/internal/shared/middleware/timing.go.tmpl": filepath.Join("internal", "shared", "middleware", "timing.go"),

		// SHARED - UTILS
		"member/internal/shared/utils/dummy.go.tmpl":      filepath.Join("internal", "shared", "utils", "dummy.go"),
		"member/internal/shared/utils/key_object.go.tmpl": filepath.Join("internal", "shared", "utils", "key_object.go"),
		"member/internal/shared/utils/uid.go.tmpl":        filepath.Join("internal", "shared", "utils", "uid.go"),

		// PKG REDIS
		"member/pkg/redis/client.go.tmpl": filepath.Join("pkg", "redis", "client.go"),
	}
	directories := []string{
		filepath.Join("internal", "shared", "domain"),
	}
	for _, dir := range directories {
		if err := os.MkdirAll(filepath.Join(baseDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	for tmpl, dest := range files {
		utils.RenderTemplate(templateMemberFS, baseDir, dest, tmpl, config)
	}

	return nil
}
