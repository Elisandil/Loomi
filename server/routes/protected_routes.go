package routes

import (
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware())

	// Movies
	router.GET("/movies", controllers.GetMovies())
	router.GET("/movie/:imdb_id", controllers.GetMovie())
	router.POST("/add_movie", controllers.AddMovie())
	router.PUT("/update_movie/:imdb_id", controllers.UpdateMovie())
	router.DELETE("/delete_movie/:imdb_id", controllers.DeleteMovie())
	router.GET("/recommended_movies", controllers.GetRecommendedMovies())
	router.PATCH("/update_review/:imdb_id", controllers.AdminReviewUpdate())

	// TV Shows
	router.GET("/tv_shows", controllers.GetTVShows())
	router.GET("/tv_shows/:imdb_id", controllers.GetTVShow())
	router.GET("/tv_show/:imdb_id/season/:season_number", controllers.GetTVShowSeason())
	router.POST("/add_tv_show", controllers.AddTVShow())
	router.PUT("/update_tv_show/:imdb_id", controllers.UpdateTVShow())
	router.POST("/tv_show/:imdb_id/add_season", controllers.AddSeason())
	router.DELETE("/delete_tv_show/:imdb_id", controllers.DeleteTVShow())
	router.PATCH("/update_tv_show_review/:imdb_id", controllers.AdminTVShowReviewUpdate())
	router.GET("/recommended_tv_shows", controllers.GetRecommendedTVShows())
}
