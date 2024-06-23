package main

import (
	"log"

	"github.com/erizkiatama/gotu-assignment/internal/app"
)

func main() {
	log.Println("starting application...")
	if err := app.Initialize(); err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
}
