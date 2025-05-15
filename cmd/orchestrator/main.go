package main

import (
	"log"

	orchestrator "github.com/Barsenick/calculator/internal/application/orchestrator"
)

func main() {
	db, err := orchestrator.OpenDB()
	if err != nil {
		log.Fatal("Error opening database:", err)
		return
	}

	orchestrator.DB = db

	app := orchestrator.New()

	err = app.RunServer()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
