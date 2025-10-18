package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"server/utils"

	"server/database"
	"server/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func DeleteTVShow() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := getDBContext()
		defer cancel()

		imdbID := c.Param("imdb_id")
		if imdbID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "IMDB ID is required"})
			return
		}

		result, err := tvShowCollection.DeleteOne(ctx, bson.M{"imdb_id": imdbID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete TV show"})
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"Error": "TV show not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Message": "TV show deleted successfully"})
	}
}

func AdminTVShowReviewUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Role not found in context"})
			return
		}
		if role != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "User is not an ADMIN"})
			return
		}

		imdbID := c.Param("imdb_id")
		if imdbID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "IMDB ID required"})
			return
		}

		var req struct {
			AdminReview string `json:"admin_review"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
			return
		}

		sentiment, rankVal, err := getReviewRankingWithHugging(req.AdminReview)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error getting review ranking"})
			return
		}

		filter := bson.M{"imdb_id": imdbID}
		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}

		ctx, cancel := getDBContext()
		defer cancel()

		result, err := tvShowCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error updating TV show"})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"Error": "TV show not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ranking_name": sentiment,
			"admin_review": req.AdminReview,
		})
	}
}

func GetRecommendedTVShows() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := utils.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "User ID not found"})
			return
		}

		favouriteGenres, err := GetUsersFavouriteGenres(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		var recommendedLimit int64 = 5
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})
		findOptions.SetLimit(recommendedLimit)

		filter := bson.M{"genre.genre_name": bson.M{"$in": favouriteGenres}}

		ctx, cancel := getDBContext()
		defer cancel()

		cursor, err := tvShowCollection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error fetching recommended TV shows"})
			return
		}
		defer cursor.Close(ctx)

		var recommendedShows []models.TVShow
		if err := cursor.All(ctx, &recommendedShows); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, recommendedShows)
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
