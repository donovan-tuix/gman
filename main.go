package main

import (
	"fmt"
	"os"

	"github.com/donovan-tuix/gman/config_manager"
)

func main() {
	if len(os.Args) != 3 || os.Args[1] != "use" {
		fmt.Println("Usage: gam use <new_account_name>")
		os.Exit(1)
	}

	newAccountName := os.Args[2]

	manager, err := config_manager.NewConfigManager(newAccountName)
	if err != nil {
		fmt.Printf("Error setting up configuration manager: %v\n", err)
		os.Exit(1)
	}
	defer manager.Close()

	newConfig, err := manager.Update()
	if err != nil {
		fmt.Printf("Error updating the configuration file: %v\n", err)
		os.Exit(1)
	}

	err = manager.WriteToFile(newConfig)
	if err != nil {
		fmt.Printf("Error writing to the configuration file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Now using account %s \n", newAccountName)
}
