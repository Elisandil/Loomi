package main

import (
	"fmt"
	"server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// entry point of the application
	router := gin.Default()

	routes.SetupUnprotectedRoutes(router)
	routes.SetupProtectedRoutes(router)
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start the server", err)
	}
}
