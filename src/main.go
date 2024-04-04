package main

import "github.com/charmbracelet/log"

func main() {
	rootCmd.AddCommand(setupCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing rootCmd: %v", err)
	}
}
