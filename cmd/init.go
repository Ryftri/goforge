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

	"github.com/AlecAivazis/survey/v2" // Jangan lupa import survey
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
		// Pastikan untuk menangani error jika pengguna membatalkan (Ctrl+C)
		err := survey.AskOne(prompt, &dbChoice)
		if err != nil {
			fmt.Println("\nOperasi dibatalkan.")
			return
		}

		fmt.Printf("🚀 Akan membuat proyek baru bernama: %s dengan database %s\n", projectName, dbChoice)

		// --- BAGIAN YANG HILANG ---
		// Sekarang kita panggil fungsi untuk membuat direktori dan file
		createDirectories(projectName)
		createFiles(projectName, dbChoice) // Kirim pilihan database ke fungsi createFiles
		// -------------------------

		fmt.Println("\n✅ Proyek berhasil dibuat!")
		fmt.Printf("Langkah selanjutnya:\n  cd %s\n  go mod tidy\n  go run cmd/api/main.go\n", projectName)
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

// TemplateData menampung data yang akan dimasukkan ke dalam template
type TemplateData struct {
	ProjectName string
	ModuleName  string // Contoh: github.com/user/my-project
	DBDriver    string // Untuk menyimpan driver DB pilihan
	DBDSN       string // Contoh DSN untuk database
}

// Fungsi untuk membuat file dari template
func createFiles(projectName string, dbChoice string) {
	fmt.Println("📄 Membuat file boilerplate...")

	// Ganti "your-user" dengan username GitHub Anda
	moduleName := fmt.Sprintf("github.com/your-user/%s", projectName)

	data := TemplateData{
		ProjectName: projectName,
		ModuleName:  moduleName,
	}

	// Logika untuk menentukan driver dan DSN berdasarkan pilihan
	switch dbChoice {
	case "PostgreSQL":
		data.DBDriver = "gorm.io/driver/postgres"
		data.DBDSN = `"host=localhost user=your_user password=your_password dbname=your_dbname port=5432 sslmode=disable"`
	case "MySQL":
		data.DBDriver = "gorm.io/driver/mysql"
		data.DBDSN = `"your_user:your_password@tcp(127.0.0.1:3306)/your_dbname?charset=utf8mb4&parseTime=True&loc=Local"`
	default:
		// Kosongkan jika pengguna memilih manual
		data.DBDriver = `// "gorm.io/driver/postgres" // Silakan uncomment dan install driver pilihan Anda`
		data.DBDSN = `"host=localhost user=your_user password=your_password dbname=your_dbname port=5432 sslmode=disable"` // Beri contoh
	}

	// Definisikan file mana yang mau dibuat dan dari template mana
	files := map[string]string{
		"go.mod":          "templates/go.mod.tmpl",
		"cmd/api/main.go": "templates/main.go.tmpl",
		// TODO: Tambahkan template untuk config.yaml, pkg/database/database.go, dll.
	}

	for dest, srcTmpl := range files {
		// Pastikan folder templates ada
		if _, err := os.Stat(srcTmpl); os.IsNotExist(err) {
			log.Printf("Peringatan: File template %s tidak ditemukan, file %s dilewati.", srcTmpl, dest)
			continue // Lanjutkan ke file berikutnya
		}

		tmplContent, err := os.ReadFile(srcTmpl)
		if err != nil {
			log.Fatalf("Gagal membaca template %s: %v", srcTmpl, err)
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
