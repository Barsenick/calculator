package main

import (
	"log"
	"os"
	"strconv"

	agent "github.com/Barsenick/calculator/internal/application/agent"
)

func main() {
	comp_power_str := os.Getenv("COMPUTING_POWER")
	comp_power := 5
	var err1 error
	if comp_power_str != "" {
		comp_power, err1 = strconv.Atoi(comp_power_str)
		if err1 != nil {
			log.Fatal(err1.Error())
		}
	}

	for range comp_power {
		go agent.StartAgent()
	}

	<-make(chan struct{})
}
