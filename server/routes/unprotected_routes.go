package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func SetupUnprotectedRoutes(router *gin.Engine) {
	router.POST("/register", controllers.RegisterUser())
	router.POST("/login", controllers.LoginUser())
}
