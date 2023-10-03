package main

import (
	"log"
	"wakeward/yaml-ctf/_app/secret-store/config"
	"wakeward/yaml-ctf/_app/secret-store/controllers"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// run database
	config.ConnectDB()

	// run OIDC
	auth, err := config.ConnectOIDC()
	if err != nil {
		log.Fatalf("Failed to connect the provider: %v", err)
	}

	controllers.UserRoutes(auth, router)

	router.Run(":5050")
}
