/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "DIAGNOSTIC: Mencoba membuat satu direktori dan mencetak error detail.",
	Long:  `Ini adalah versi khusus untuk debugging. Program ini hanya akan mencoba membuat direktori utama proyek dan akan mencetak informasi error selengkap mungkin jika gagal.`,
	Args:  cobra.ExactArgs(1), // Memastikan nama proyek harus ada
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		// Dapatkan direktori kerja saat ini untuk konteks
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Gagal mendapatkan direktori kerja saat ini: %v\n", err)
		}
		fmt.Printf("-> Direktori kerja saat ini: %s\n", wd)
		fmt.Printf("-> Nama proyek yang diberikan: %s\n", projectName)

		// Ubah path proyek menjadi path absolut untuk kejelasan
		absPath, err := filepath.Abs(projectName)
		if err != nil {
			fmt.Printf("Gagal mengubah '%s' menjadi path absolut: %v\n", projectName, err)
			return
		}
		fmt.Printf("-> Path absolut yang akan dibuat: %s\n", absPath)

		// Coba buat HANYA direktori utama
		fmt.Println("\n--- MEMULAI OPERASI os.MkdirAll ---")
		err = os.MkdirAll(absPath, 0755)
		if err != nil {
			// Ini bagian paling penting. Kita cetak errornya dengan detail.
			fmt.Printf("\n--- GAGAL ---\n")
			fmt.Printf("Operasi os.MkdirAll gagal dengan error:\n%v\n", err)
			fmt.Println("\n--- SELESAI ---")
			os.Exit(1)
		}

		// Jika berhasil
		fmt.Println("\n--- SUKSES ---")
		fmt.Printf("Direktori '%s' berhasil dibuat.\n", absPath)
		fmt.Println("--- SELESAI ---")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
