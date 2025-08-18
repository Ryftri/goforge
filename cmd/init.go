/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new Go modular monolith project",
	Long:  `Creates a new project with a predefined modular structure, including API versioning and domain-driven design principles.`,
	Args:  cobra.ExactArgs(1), // Memastikan nama proyek harus ada
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		// --- TANYA DATABASE ---
		dbChoice := ""
		prompt := &survey.Select{
			Message: "Pilih database yang akan Anda gunakan:",
			Options: []string{"PostgreSQL", "MySQL", "Lainnya (install manual)"},
		}
		survey.AskOne(prompt, &dbChoice)
		// ----------------------

		fmt.Printf("🚀 Akan membuat proyek baru bernama: %s dengan database %s\n", projectName, dbChoice)

		// ... panggil createDirectories dan createFiles ...
		// Anda bisa passing 'dbChoice' ke fungsi createFiles
		// untuk men-generate go.mod dan kode database yang sesuai.
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// Fungsi untuk membuat direktori
func createDirectories(projectName string) {
	fmt.Println("📂 Membuat struktur direktori...")
	dirs := []string{
		"cmd/api",
		"api/v1/handler",
		"api/v1/request",
		"api/v1/response",
		"internal/category", // Contoh domain pertama
		"pkg/config",
		"pkg/database",
		"migrations",
	}

	for _, dir := range dirs {
		path := filepath.Join(projectName, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			fmt.Printf("Error membuat direktori %s: %v\n", path, err)
			os.Exit(1)
		}
	}
}

// Fungsi untuk membuat file
type TemplateData struct {
	ProjectName string
	ModuleName  string // Contoh: github.com/user/my-project
}

func createFiles(projectName string) {
	fmt.Println("📄 Membuat file boilerplate...")

	// Dapatkan module name, bisa ditanyakan ke user atau generate otomatis
	moduleName := fmt.Sprintf("github.com/your-user/%s", projectName) // Sederhanakan dulu

	data := TemplateData{
		ProjectName: projectName,
		ModuleName:  moduleName,
	}

	// Definisikan file mana yang mau dibuat dan dari template mana
	files := map[string]string{
		"go.mod":          "templates/go.mod.tmpl",
		"cmd/api/main.go": "templates/main.go.tmpl",
		// Tambahkan file lainnya di sini
	}

	for dest, srcTmpl := range files {
		// Baca template
		tmplContent, err := os.ReadFile(srcTmpl) // Nanti ganti dengan embed
		if err != nil {
			log.Fatalf("Gagal membaca template %s: %v", srcTmpl, err)
		}

		// Buat file tujuan
		destPath := filepath.Join(projectName, dest)
		file, err := os.Create(destPath)
		if err != nil {
			log.Fatalf("Gagal membuat file %s: %v", destPath, err)
		}
		defer file.Close()

		// Eksekusi template
		tmpl, err := template.New(dest).Parse(string(tmplContent))
		if err != nil {
			log.Fatalf("Gagal mem-parsing template %s: %v", srcTmpl, err)
		}

		if err := tmpl.Execute(file, data); err != nil {
			log.Fatalf("Gagal mengeksekusi template %s: %v", srcTmpl, err)
		}
	}
}
