/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
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

// Make sure your folder structure is: goforge/cmd/templates/
//
//go:embed templates/*
var templateFiles embed.FS

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new Go modular monolith project",
	Long:  `Creates a new project with a predefined modular structure, including API versioning and domain-driven design principles.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		// --- Ask for Module Name (New Feature) ---
		moduleName := ""
		modulePrompt := &survey.Input{
			Message: "Enter the Go module name:",
			Default: fmt.Sprintf("github.com/your-username/%s", projectName),
		}
		err := survey.AskOne(modulePrompt, &moduleName)
		if err != nil {
			fmt.Println("\nOperation cancelled.")
			return
		}

		// --- Ask for Database ---
		dbChoice := ""
		dbPrompt := &survey.Select{
			Message: "Select the database you will use:",
			Options: []string{"PostgreSQL", "MySQL", "Other (manual install)"},
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
		// Pass the moduleName to createFiles
		createFiles(projectName, dbChoice, moduleName)
		runGoGet(projectName, dbChoice)
		runGoGenerate(projectName) // Automatically run go generate

		// Update the success message
		fmt.Println("\n✅ Project ready to run! All dependencies installed and code generated.")
		fmt.Printf("Next steps:\n  cd %s\n  go run cmd/api/main.go\n", projectName)
	},
}

// runGoGet executes 'go get' for each necessary package.
func runGoGet(projectName string, dbChoice string) {
	fmt.Println("📦 Installing dependencies (go get)...")

	// List of packages to install
	packages := []string{
		"github.com/gin-gonic/gin",
		"github.com/spf13/viper",
		"gorm.io/gorm",
		"github.com/google/wire/cmd/wire",
		"github.com/stretchr/testify",
		"github.com/vektra/mockery/v2/...@latest", // Added mockery
	}

	// Conditionally add database drivers
	switch dbChoice {
	case "PostgreSQL":
		packages = append(packages, "gorm.io/driver/postgres")
	case "MySQL":
		packages = append(packages, "gorm.io/driver/mysql")
	}

	// Run 'go get' for each package
	for _, pkg := range packages {
		fmt.Printf("   - Installing %s...\n", pkg)
		cmd := exec.Command("go", "get", pkg)
		cmd.Dir = projectName // Run inside the new project folder

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Failed to run 'go get %s': %v\nOutput:\n%s", pkg, err, string(output))
		}
	}

	fmt.Println("Dependencies installed successfully.")
}

// runGoGenerate executes 'go generate ./...' to produce the wire_gen.go file.
func runGoGenerate(projectName string) {
	fmt.Println("⚙️  Generating dependency injection code (go generate)...")
	cmd := exec.Command("go", "generate", "./...")
	cmd.Dir = projectName // Run inside the new project folder

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to run 'go generate': %v\nOutput:\n%s", err, string(output))
	}
	fmt.Println("Code generation complete.")
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func createDirectories(projectName string) {
	fmt.Println("📂 Creating directory structure...")
	dirs := []string{
		"cmd/api",
		"api/v1/handler",
		"api/v1/request",
		"api/v1/response",
		"internal/category",
		"pkg/config",
		"pkg/database",
		"migrations",
	}
	for _, dir := range dirs {
		path := filepath.Join(projectName, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalf("Error creating directory %s: %v\n", path, err)
		}
	}
}

type TemplateData struct {
	ProjectName string
	ModuleName  string
	DBDSN       string
}

// Updated function signature to accept moduleName
func createFiles(projectName string, dbChoice string, moduleName string) {
	fmt.Println("📄 Creating boilerplate files...")
	data := TemplateData{
		ProjectName: projectName,
		ModuleName:  moduleName, // Use the moduleName from user input
	}
	switch dbChoice {
	case "PostgreSQL":
		data.DBDSN = `"host=localhost user=your_user password=your_password dbname=your_dbname port=5432 sslmode=disable"`
	case "MySQL":
		data.DBDSN = `"your_user:your_password@tcp(127.0.0.1:3306)/your_dbname?charset=utf8mb4&parseTime=True&loc=Local"`
	default:
		data.DBDSN = `"host=localhost user=your_user password=your_password dbname=your_dbname port=5432 sslmode=disable"`
	}

	files := map[string]string{
		".mockery.yaml":           "templates/.mockery.yaml.tmpl",
		"go.mod":                  "templates/go.mod.tmpl",
		"cmd/api/main.go":         "templates/main.go.tmpl",
		"cmd/api/wire.go":         "templates/wire.go.tmpl",
		"config.yaml":             "templates/config.yaml.tmpl",
		"pkg/config/config.go":    "templates/config.go.tmpl",
		"api/v1/router.go":        "templates/router.go.tmpl",
		"api/v1/handler/hello.go": "templates/hello_handler.go.tmpl",
	}

	for dest, srcTmpl := range files {
		tmplContent, err := templateFiles.ReadFile(srcTmpl)
		if err != nil {
			log.Fatalf("Failed to read template from embed %s: %v", srcTmpl, err)
		}
		destPath := filepath.Join(projectName, dest)
		file, err := os.Create(destPath)
		if err != nil {
			log.Fatalf("Failed to create file %s: %v", destPath, err)
		}
		defer file.Close()
		tmpl, err := template.New(dest).Parse(string(tmplContent))
		if err != nil {
			log.Fatalf("Failed to parse template %s: %v", srcTmpl, err)
		}
		if err := tmpl.Execute(file, data); err != nil {
			log.Fatalf("Failed to execute template %s: %v", srcTmpl, err)
		}
	}
}
