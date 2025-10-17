package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"server/database"
	"server/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var tvShowCollection *mongo.Collection = database.OpenCollection("tv_shows")
var tvShowValidator = validator.New()

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

func GetTVShowSeason() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := getDBContext()
		defer cancel()

		imdbID := c.Param("imdb_id")
		seasonNumber := c.Param("season_number")

		if imdbID == "" || seasonNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "IMDB ID and season number are required"})
			return
		}

		var tvShow models.TVShow
		err := tvShowCollection.FindOne(ctx, bson.M{"imdb_id": imdbID}).Decode(&tvShow)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"Error": "TV show not found"})
			return
		}

		var foundSeason *models.Season
		for _, season := range tvShow.Seasons {
			if season.SeasonNumber == parseSeasonNumber(seasonNumber) {
				foundSeason = &season
				break
			}
		}

		if foundSeason == nil {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Season not found"})
			return
		}

		c.JSON(http.StatusOK, foundSeason)
	}
}

func AddTVShow() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := getDBContext()
		defer cancel()

		var tvShow models.TVShow
		if err := c.ShouldBindJSON(&tvShow); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input", "details": err.Error()})
			return
		}
		if len(tvShow.Seasons) != tvShow.TotalSeasons {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Total seasons doesn't match the number of seasons provided"})
			return
		}

		if err := tvShowValidator.Struct(tvShow); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Validation failed", "details": err.Error()})
			return
		}

		var existingShow models.TVShow

		err := tvShowCollection.FindOne(ctx, bson.M{"title": tvShow.Title}).Decode(&existingShow)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"Error": "A TV show with this title already exists"})
			return
		}

		result, err := tvShowCollection.InsertOne(ctx, tvShow)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to add TV show"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func UpdateTVShow() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := getDBContext()
		defer cancel()

		imdbID := c.Param("imdb_id")
		if imdbID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "IMDB ID is required"})
			return
		}

		var tvShow models.TVShow
		if err := c.ShouldBindJSON(&tvShow); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input"})
			return
		}

		if err := tvShowValidator.Struct(tvShow); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Validation failed", "details": err.Error()})
			return
		}

		result, err := tvShowCollection.ReplaceOne(
			ctx,
			bson.M{"imdb_id": imdbID},
			tvShow,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update TV show"})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"Error": "TV show not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Message":        "TV show updated successfully",
			"modified_count": result.ModifiedCount,
		})
	}
}

func AddSeason() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := getDBContext()
		defer cancel()

		imdbID := c.Param("imdb_id")
		if imdbID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "ID is required"})
			return
		}

		var season models.Season
		if err := c.ShouldBindJSON(&season); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input"})
			return
		}

		if err := tvShowValidator.Struct(season); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Validation failed", "details": err.Error()})
			return
		}

		var tvShow models.TVShow
		err := tvShowCollection.FindOne(ctx, bson.M{"imdb_id": imdbID}).Decode(&tvShow)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Season already exists"})
			return
		}
		for _, s := range tvShow.Seasons {
			if s.SeasonNumber == season.SeasonNumber {
				c.JSON(http.StatusConflict, gin.H{"Error": "Season already exists"})
				return
			}
		}

		update := bson.M{
			"$push": bson.M{"seasons": season},
			"$inc":  bson.M{"total_seasons": 1},
		}

		result, err := tvShowCollection.UpdateOne(ctx, bson.M{"imdb_id": imdbID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to add season"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"Message":        "Season added successfully",
			"modified_count": result.ModifiedCount,
		})
	}
}

// Utility functions
// ---------------------------------------------------------------------------------------

func parseSeasonNumber(s string) int {
	var num int
	if _, err := fmt.Sscanf(s, "%d", &num); err != nil {
		return 0
	}
	return num
}
