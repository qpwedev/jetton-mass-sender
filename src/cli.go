package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/spf13/cobra"
)

// Configuration stores the user's settings
type Configuration struct {
	SeedPhrase   string `json:"seedPhrase"`
	JettonMasterAddress string `json:"jettonMasterAddress"`
	Commentary   string `json:"commentary"`
	MessageEntryFilename string `json:"messageEntryFilename"`
}

const configFilePath = "config.json"

var rootCmd = &cobra.Command{
	Use:   "toncli",
	Short: "TON CLI application for managing transactions",
	Run:   runMainLogic,
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up or update the application configuration",
	Run:   func(cmd *cobra.Command, args []string) { setupConfiguration() },
}


func runMainLogic(cmd *cobra.Command, args []string) {
	config, err := loadConfiguration(configFilePath)
	if err != nil {
		fmt.Println("Configuration not found or invalid. Please run `setup` command.")
		return
	}
	fmt.Println("Running main logic with current configuration...")

	massSender(config.SeedPhrase, config.JettonMasterAddress, config.Commentary, config.MessageEntryFilename)
}

func setupConfiguration() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter Seed Phrase (From Any Wallet):")
	seedPhrase, _ := reader.ReadString('\n')

	fmt.Println("Enter Token Address:")
	jettonMasterAddress, _ := reader.ReadString('\n')

	fmt.Println("Enter Commentary:")
	commentary, _ := reader.ReadString('\n')

	fmt.Println("Enter Message Entry Filename:")
	messageEntryFilename, _ := reader.ReadString('\n')

	config := Configuration{
		SeedPhrase:           seedPhrase,
		JettonMasterAddress:         jettonMasterAddress,
		Commentary:           commentary,
		MessageEntryFilename: messageEntryFilename,
	}
	if err := saveConfiguration(config, configFilePath); err != nil {
		log.Fatalf("Failed to save configuration: %v", err)
	}

	fmt.Println("Configuration saved successfully.")
}


func saveConfiguration(config Configuration, filePath string) error {
	// Trim whitespace from the configuration values
	config.SeedPhrase = strings.TrimSpace(config.SeedPhrase)
	config.JettonMasterAddress = strings.TrimSpace(config.JettonMasterAddress)
	config.Commentary = strings.TrimSpace(config.Commentary)
	config.MessageEntryFilename = strings.TrimSpace(config.MessageEntryFilename)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(config)
}

func loadConfiguration(filePath string) (Configuration, error) {
	var config Configuration
	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}
