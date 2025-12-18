package main

import (
	"fmt"
	"log"
	"os"

	"uas-go/config"
	"uas-go/database"
	"uas-go/route"

	_ "uas-go/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Sistem Pelaporan Prestasi Mahasiswa API
// @version 1.0
// @description Backend API Sistem Pelaporan Prestasi Mahasiswa
// @termsOfService http://swagger.io/terms/

// @contact.name Aisha Purwanto
// @contact.email aishapurwanto249@gmail.com

// @host localhost:3000
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadEnv()
	

	if err := database.ConnectPostgres(); err != nil {
		log.Fatal(err)
	}
	if err := database.ConnectMongo(); err != nil {
		log.Fatal(err)
	}

	app := config.NewFiberApp()

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	route.RegisterRoutes(app)

	port := os.Getenv("APP_PORT")
	fmt.Println("Server running on port", port)

	log.Fatal(app.Listen(":" + port))
}
