/*
Copyright © 2025 Rayhan Zulfitri <rayhanzulfitri@gmail.com>
*/
package cmd

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

//go:embed templates/*
var templateFiles embed.FS

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new Go project with a layered architecture",
	Long:  `Creates a new project with a predefined structure, including a "Hello World" example, ready for you to build upon.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		moduleName := ""
		err := survey.AskOne(&survey.Input{
			Message: "Enter the Go module name:",
			Default: fmt.Sprintf("github.com/your-username/%s", projectName),
		}, &moduleName)
		if err != nil {
			fmt.Println("\nOperation cancelled.")
			return
		}

		dbChoice := ""
		dbPrompt := &survey.Select{
			Message: "Select the database you will use:",
			Options: []string{"PostgreSQL", "MySQL"},
			Default: "PostgreSQL",
		}
		err = survey.AskOne(dbPrompt, &dbChoice)
		if err != nil {
			fmt.Println("\nOperation cancelled.")
			return
		}

		fmt.Printf("🚀 Creating a new project named: %s\n", projectName)
		fmt.Printf("   - Module: %s\n", moduleName)
		fmt.Printf("   - Database: %s\n", dbChoice)

		createDirectories(projectName)
		createFiles(projectName, dbChoice, moduleName)
		installDependencies(projectName, dbChoice)
		runGoGenerate(projectName) // Diubah ke metode yang lebih baik

		fmt.Println("\n✅ Project successfully created and ready to run!")
		fmt.Println("\nNext steps:")
		fmt.Printf("  1. cd %s\n", projectName)
		fmt.Printf("  2. (Optional) Edit configs/config.yaml with your database details.\n")
		fmt.Println("  3. go run ./cmd/api")
		fmt.Printf("\nTest your endpoint:\n  curl http://localhost:3000/api/v1/hello\n")
	},
}

// ... fungsi createDirectories tidak berubah ...
func createDirectories(projectName string) {
	fmt.Println("📂 Creating directory structure...")
	dirs := []string{
		"cmd/api", "configs", "internal/config", "internal/database", "internal/handler/v1",
		"internal/logger", "internal/model/domain", "internal/model/web", "internal/repository",
		"internal/service", "internal/router", "internal/server", "internal/util", "logs",
		"migrations", "api/v1",
	}
	for _, dir := range dirs {
		path := filepath.Join(projectName, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalf("Error creating directory %s: %v\n", path, err)
		}
	}
}

// ... fungsi createFiles tidak berubah ...
type TemplateData struct {
	ModuleName   string
	DBDriver     string
	DBImportPath string
}

func createFiles(projectName string, dbChoice string, moduleName string) {
	fmt.Println("📄 Creating boilerplate files...")
	data := TemplateData{ModuleName: moduleName}

	switch dbChoice {
	case "PostgreSQL":
		data.DBDriver, data.DBImportPath = "postgres", "gorm.io/driver/postgres"
	case "MySQL":
		data.DBDriver, data.DBImportPath = "mysql", "gorm.io/driver/mysql"
	}

	files := map[string]string{
		".gitignore":                              "templates/gitignore.tmpl",
		".mockery.yaml":                           "templates/mockery.yaml.tmpl",
		"go.mod":                                  "templates/go.mod.tmpl",
		"cmd/api/main.go":                         "templates/main.go.tmpl",
		"configs/config.yaml":                     "templates/config.yaml.tmpl",
		"internal/config/config.go":               "templates/config.go.tmpl",
		"internal/database/database.go":           "templates/database.go.tmpl",
		"internal/handler/v1/hello_handler.go":    "templates/hello_handler.go.tmpl",
		"internal/logger/logger.go":               "templates/logger.go.tmpl",
		"internal/model/web/standard_response.go": "templates/standard_response.go.tmpl",
		"internal/repository/hello_repository.go": "templates/hello_repository.go.tmpl",
		"internal/service/hello_service.go":       "templates/hello_service.go.tmpl",
		"internal/router/router.go":               "templates/router.go.tmpl",
		"internal/server/wire.go":                 "templates/wire.go.tmpl",
		"internal/util/response.go":               "templates/response.go.tmpl",
		"api/v1/openapi.yaml":                     "templates/openapi.yaml.tmpl",
	}

	for dest, srcTmpl := range files {
		tmplContent, _ := templateFiles.ReadFile(srcTmpl)
		destPath := filepath.Join(projectName, dest)
		file, _ := os.Create(destPath)
		defer file.Close()
		tmpl, _ := template.New(dest).Parse(string(tmplContent))
		tmpl.Execute(file, data)
	}

	placeholders := []string{"logs/.gitkeep", "migrations/.gitkeep", "internal/model/domain/.gitkeep"}
	for _, p := range placeholders {
		os.Create(filepath.Join(projectName, p))
	}
}

// ... fungsi installDependencies tidak berubah ...
func installDependencies(projectName string, dbChoice string) {
	fmt.Println("📦 Installing dependencies...")
	packages := []string{
		"github.com/gin-gonic/gin", "github.com/spf13/viper", "gorm.io/gorm",
		"github.com/google/wire/cmd/wire", "github.com/go-playground/validator/v10", "github.com/vektra/mockery/v3",
	}

	switch dbChoice {
	case "PostgreSQL":
		packages = append(packages, "gorm.io/driver/postgres")
	case "MySQL":
		packages = append(packages, "gorm.io/driver/mysql")
	}

	for _, pkg := range packages {
		fmt.Printf("   - Getting %s...\n", pkg)
		cmd := exec.Command("go", "get", pkg)
		cmd.Dir = projectName
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Fatalf("Failed to run 'go get %s': %v\nOutput:\n%s", pkg, err, string(output))
		}
	}

	fmt.Println("🧹 Tidying modules (go mod tidy)...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectName
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to run 'go mod tidy': %v\nOutput:\n%s", err, string(output))
	}
	fmt.Println("Dependencies installed successfully.")
}

// --- FUNGSI YANG DIPERBARUI ---
func runGoGenerate(projectName string) {
	fmt.Println("⚙️  Generating dependency injection code (go generate)...")
	// Menjalankan `go generate` pada direktori spesifik untuk menghindari masalah kompilasi
	cmd := exec.Command("go", "generate", "./internal/server")
	cmd.Dir = projectName // Pastikan perintah dijalankan di dalam direktori proyek baru
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to run 'go generate': %v\nOutput:\n%s", err, string(output))
	}
	fmt.Println("Code generation complete.")
}

func init() {
	rootCmd.AddCommand(initCmd)
}
