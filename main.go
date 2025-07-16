package main

import (
	"fmt"
	"log"

	"github.com/NikolaTosic-sudo/gator/internal/config"
)

func main() {
	cfg, err := config.Read()

	if err != nil {
		log.Fatalf("error reading config %v", err)
	}

	fmt.Printf("Read config: %v\n", cfg)

	errName := cfg.SetUser("nikola")

	if errName != nil {
		log.Fatalf("error reading config %v", errName)
	}

	cfg, err = config.Read()

	if err != nil {
		log.Fatalf("error reading config %v", errName)
	}

	fmt.Printf("Read config again: %v\n", cfg)
}
