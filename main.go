package main

import (
	"fmt"
	"log"
	"os"

	"uas-go/config"
	"uas-go/database"
	"uas-go/route"
)

func main() {
	config.LoadEnv()
	

	if err := database.ConnectPostgres(); err != nil {
		log.Fatal(err)
	}
	if err := database.ConnectMongo(); err != nil {
		log.Fatal(err)
	}

	app := config.NewFiberApp()

	route.RegisterRoutes(app)

	port := os.Getenv("APP_PORT")
	fmt.Println("Server running on port", port)

	log.Fatal(app.Listen(":" + port))
}
