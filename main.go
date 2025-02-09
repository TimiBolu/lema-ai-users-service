package main

import (
	"log"

	"github.com/TimiBolu/lema-ai-users-service/config"
	"github.com/TimiBolu/lema-ai-users-service/database"
	"github.com/TimiBolu/lema-ai-users-service/router"
)

func main() {
	// Initialise the Server environment
	config.InitEnvSchema()

	// Connect to the database
	database.Connect()
	log.Println("✅ Database connection established")

	// start the router
	router.Setup()
}
