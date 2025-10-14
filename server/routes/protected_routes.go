package routes

import (
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware())

	router.GET("/movies", controllers.GetMovies())
	router.GET("/movie/:imdb_id", controllers.GetMovie())
	router.POST("/add_movie", controllers.AddMovie())
	router.PUT("/update_movie/:imdb_id", controllers.UpdateMovie())
	router.DELETE("/delete_movie/:imdb_id", controllers.DeleteMovie())
}
