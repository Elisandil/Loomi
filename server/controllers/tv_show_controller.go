package controllers

import (
	"errors"
	"net/http"

	"server/database"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var tvShowCollection *mongo.Collection = database.OpenCollection("tv_shows")

func GetTVShows() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := getDBContext()
		defer cancel()

		var tvShows []models.TVShow

		cursor, err := tvShowCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch TV shows"})
			return
		}
		defer cursor.Close(ctx)
		if err = cursor.All(ctx, &tvShows); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to decode TV shows"})
			return
		}

		c.JSON(http.StatusOK, tvShows)
	}
}

func GetTVShow() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := getDBContext()
		defer cancel()

		imdbID := c.Param("imdb_id")
		if imdbID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "IMDB ID is required"})
			return
		}

		var tvShow models.TVShow
		err := tvShowCollection.FindOne(ctx, bson.M{"imdb_id": imdbID}).Decode(&tvShow)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusNotFound, gin.H{"Error": "TV show not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch TV show"})
			return
		}

		c.JSON(http.StatusOK, tvShow)
	}
}
