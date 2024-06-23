package main

import (
	"log"

	"github.com/erizkiatama/gotu-assignment/internal/app"
	"github.com/erizkiatama/gotu-assignment/internal/config"
)

func main() {
	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	log.Println("starting application...")
	if err := app.Initialize(cfg); err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
}
