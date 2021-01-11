package main

import (
	"log"
	"os"

	"github.com/gradusp/crispy/ctrl/server"
)

func main() {
	app := server.NewApp()

	if err := app.Run("8080"); err != nil {
		log.Printf("%s", err.Error())
		os.Exit(1)
	}
}
