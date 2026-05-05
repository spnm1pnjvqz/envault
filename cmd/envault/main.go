package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envault",
	Short: "Lightweight secrets manager using age encryption",
	Long:  `envault encrypts and decrypts .env files using age encryption for local dev workflows.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(lockCmd)
	rootCmd.AddCommand(unlockCmd)
	rootCmd.AddCommand(viewCmd)
}
