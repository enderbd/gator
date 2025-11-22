package main

import (
	"fmt"
	"log"

	"github.com/enderbd/gator/internal/config"
)
func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file %v", err)
	}
	if err := cfg.SetUser("nemo"); err != nil {
		log.Fatalf("Error setting the user name %v", err)
	}
	
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading config file %v", err)
	}

	fmt.Printf("Contenf of config struct: %+v\n", cfg)

}
