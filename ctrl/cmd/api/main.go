package main

import (
	"github.com/gradusp/crispy/ctrl/server"
	"log"
)

func main() {
	app := server.NewApp()

	if err := app.Run("8080"); err != nil {
		log.Fatal("%s", err.Error())
	}
}
