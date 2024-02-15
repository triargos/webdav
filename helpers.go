package main

import (
	"log"
	"os"
)

func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Printf("Error closing file: %s \n Cause: %s \n", file.Name(), err.Error())
	}
}

func handleConfigurationReadError(err error) bool {
	log.Printf("Error reading configuration file: %s \n", err.Error())
	return false
}

// initStorage creates the data directory if it does not exist
func initStorage(path string) {
	// Create the data directory if it does not exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			log.Fatalf("Could not create data directory: %s", err.Error())
		}
	}
}
