package main

import (
	"log"
	"uas-go/config"
)

func main() {
	app := config.SetupApp()

	port := "3000"
	log.Println("Running on port", port)

	app.Listen(":" + port)
}
