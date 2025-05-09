package main

import (
	"log"
	"os"
	"strconv"

	agent "github.com/Barsenick/calculator/internal/application/agent"
	orchestrator "github.com/Barsenick/calculator/internal/application/orchestrator"
)

func StartAgents() error {
	comp_power_str := os.Getenv("COMPUTING_POWER")
	comp_power := 5
	var err1 error
	if comp_power_str != "" {
		comp_power, err1 = strconv.Atoi(comp_power_str)
		if err1 != nil {
			return err1
		}
	}

	for range comp_power {
		go agent.StartAgent()
	}

	return nil
}

func main() {
	go func() {
		err := StartAgents()
		if err != nil {
			log.Fatal("Error starting agents:", err)
		}
	}()

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
