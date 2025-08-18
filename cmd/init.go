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

// PENTING: AGAR EMBED BERHASIL, STRUKTUR FOLDER ANDA HARUS SEPERTI INI:
//
// goforge/
// ├── cmd/
// │   ├── init.go       <-- File ini
// │   ├── root.go
// │   └── templates/    <-- Pindahkan folder templates ke sini
// │       ├── config.go.tmpl
// │       ├── config.yaml.tmpl
// │       ├── go.mod.tmpl
// │       └── main.go.tmpl
// ├── go.mod
// └── main.go
//
// Path 'templates/*' di bawah ini bersifat relatif terhadap file init.go ini.
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

		dbChoice := ""
		prompt := &survey.Select{
			Message: "Pilih database yang akan Anda gunakan:",
			Options: []string{"PostgreSQL", "MySQL", "Lainnya (install manual)"},
		}
		err := survey.AskOne(prompt, &dbChoice)
		if err != nil {
			fmt.Println("\nOperasi dibatalkan.")
			return
		}

		fmt.Printf("🚀 Akan membuat proyek baru bernama: %s dengan database %s\n", projectName, dbChoice)

		createDirectories(projectName)
		createFiles(projectName, dbChoice)

		runGoModTidy(projectName)

		fmt.Println("\n✅ Proyek berhasil dibuat dan semua dependensi telah diinstall!")
		fmt.Printf("Langkah selanjutnya:\n  cd %s\n  go run cmd/api/main.go\n", projectName)
	},
}

// Fungsi baru untuk menjalankan 'go mod tidy'
func runGoModTidy(projectName string) {
	fmt.Println("📦 Menginstall dependensi (go mod tidy)...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectName // PENTING: Jalankan perintah di dalam folder proyek baru

	// Tampilkan output dari perintah untuk debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Gagal menjalankan 'go mod tidy': %v\nOutput:\n%s", err, string(output))
	}
	fmt.Println("Dependensi berhasil diinstall.")
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func createDirectories(projectName string) {
	fmt.Println("📂 Membuat struktur direktori...")
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
			log.Fatalf("Error membuat direktori %s: %v\n", path, err)
		}
	}
}

type TemplateData struct {
	ProjectName string
	ModuleName  string
	DBDriver    string
	DBDSN       string
}

func createFiles(projectName string, dbChoice string) {
	fmt.Println("📄 Membuat file boilerplate...")

	moduleName := fmt.Sprintf("github.com/Ryftri/%s", projectName)

	data := TemplateData{
		ProjectName: projectName,
		ModuleName:  moduleName,
	}

	switch dbChoice {
	case "PostgreSQL":
		data.DBDriver = "gorm.io/driver/postgres"
		data.DBDSN = `"host=localhost user=your_user password=your_password dbname=your_dbname port=5432 sslmode=disable"`
	case "MySQL":
		data.DBDriver = "gorm.io/driver/mysql"
		data.DBDSN = `"your_user:your_password@tcp(127.0.0.1:3306)/your_dbname?charset=utf8mb4&parseTime=True&loc=Local"`
	default:
		data.DBDriver = `// "gorm.io/driver/postgres" // Silakan uncomment dan install driver pilihan Anda`
		data.DBDSN = `"host=localhost user=your_user password=your_password dbname=your_dbname port=5432 sslmode=disable"`
	}

	files := map[string]string{
		"go.mod":               "templates/go.mod.tmpl",
		"cmd/api/main.go":      "templates/main.go.tmpl",
		"config.yaml":          "templates/config.yaml.tmpl",
		"pkg/config/config.go": "templates/config.go.tmpl",
	}

	for dest, srcTmpl := range files {
		tmplContent, err := templateFiles.ReadFile(srcTmpl)
		if err != nil {
			log.Fatalf("Gagal membaca template dari embed %s: %v", srcTmpl, err)
		}

		destPath := filepath.Join(projectName, dest)
		file, err := os.Create(destPath)
		if err != nil {
			log.Fatalf("Gagal membuat file %s: %v", destPath, err)
		}
		defer file.Close()

		tmpl, err := template.New(dest).Parse(string(tmplContent))
		if err != nil {
			log.Fatalf("Gagal mem-parsing template %s: %v", srcTmpl, err)
		}

		if err := tmpl.Execute(file, data); err != nil {
			log.Fatalf("Gagal mengeksekusi template %s: %v", srcTmpl, err)
		}
	}
}
