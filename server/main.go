package main

import (
	"fmt"
	"time"

	"server/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// entry point of the application
	router := gin.Default()

	config := cors.Config{}
	config.AllowAllOrigins = true
	config.AllowMethods = []string{
		"GET", "POST", "PUT", "PATCH", "DELETE",
	}
	config.AllowHeaders = []string{
		"Origin", "Content-Type", "Authorization",
	}
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	router.Use(gin.Logger())

	routes.SetupUnprotectedRoutes(router)
	routes.SetupProtectedRoutes(router)
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start the server", err)
	}
}
