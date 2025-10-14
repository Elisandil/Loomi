package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func SetupUnprotectedRoutes(router *gin.Engine) {
	router.POST("/register", controllers.RegisterUser())
	router.POST("/login", controllers.LoginUser())
	router.PATCH("/update_review/:imdb_id", controllers.AdminReviewUpdate())
}
