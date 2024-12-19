package main

import "github.com/Barsenick/calculator/internal/application"

func main() {
	app := application.New()
	app.RunServer()
}
