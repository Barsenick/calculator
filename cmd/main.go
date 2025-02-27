package main

import (
	"fmt"

	"github.com/Barsenick/calculator/internal/application"
)

func main() {
	app := application.New()
	err := app.RunServer()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
