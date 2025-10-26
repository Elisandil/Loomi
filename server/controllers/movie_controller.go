package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"server/database"
	"server/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var rankingCollection *mongo.Collection = database.OpenCollection("rankings")
var genreCollection *mongo.Collection = database.OpenCollection("genres")
var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := getDBContext()
		defer cancel()
		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch movies"})
		}
		defer cursor.Close(ctx)
		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to decode movies"})
		}

		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := getDBContext()
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Movie ID is required"})
			return
		}

		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := getDBContext()
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input"})
			return
		}
		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Validation failed", "details": err.Error()})
			return
		}

		var existingMovie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"title": movie.Title}).Decode(&existingMovie)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"Error": "A movie with this title already exists"})
			return
		}

		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Operation to add a movie failed"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func UpdateMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := getDBContext()
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Movie ID is required"})
			return
		}

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input"})
			return
		}

		result, err := movieCollection.ReplaceOne(
			ctx,
			bson.M{"imdb_id": movieID},
			movie,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update movie"})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Message":        "Movie updated successfully",
			"modified_count": result.ModifiedCount,
		})
	}
}

func DeleteMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := getDBContext()
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Movie ID is required"})
			return
		}

		result, err := movieCollection.DeleteOne(ctx, bson.M{"imdb_id": movieID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete the movie"})
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Message": "Movie deleted successfully"})
	}
}

func AdminReviewUpdate() gin.HandlerFunc {
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

		movieId := c.Param("imdb_id")
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Movie Id required"})
			return
		}

		var req struct {
			AdminReview string `json:"admin_review"`
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
			return
		}

		sentiment, rankVal, err := getReviewRankingWithHugging(req.AdminReview)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error getting review ranking"})
			return
		}

		filter := bson.M{"imdb_id": movieId}
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

		result, err := movieCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error updating movie"})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Movie not found"})
			return
		}

		var response struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}
		response.RankingName = sentiment
		response.AdminReview = req.AdminReview

		c.JSON(http.StatusOK, response)
	}
}

func GetUsersFavouriteGenres(userId string) ([]string, error) {
	ctx, cancel := getDBContext()
	defer cancel()

	filter := bson.M{"user_id": userId}

	projection := bson.M{
		"favourite_genres.genre_name": 1,
		"_id":                         0,
	}

	opts := options.FindOne().SetProjection(projection)

	var result bson.M
	err := usersCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []string{}, nil
		}
	}

	favGenres, ok := result["favourite_genres"].(bson.A)
	if !ok {
		return []string{}, errors.New("unable to retrieve favourite genres for user")
	}
	var genreNames []string
	for _, item := range favGenres {
		if genreMap, ok := item.(bson.D); ok {
			for _, elems := range genreMap {
				if elems.Key == "genre_name" {
					if name, ok := elems.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}

	return genreNames, nil
}

func GetRecommendedMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := utils.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "User Id not found"})
			return
		}

		favouriteGenres, err := GetUsersFavouriteGenres(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		// for development
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("Warning: .env not found")
		}

		var recommendedMovieLimitVal int64 = 5
		recommendedMovieLimitStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT")
		if recommendedMovieLimitStr != "" {
			recommendedMovieLimitVal, _ = strconv.ParseInt(recommendedMovieLimitStr, 10, 64)
		}

		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})
		findOptions.SetLimit(recommendedMovieLimitVal)

		filter := bson.M{"genre.genre_name": bson.M{"$in": favouriteGenres}}

		ctx, cancel := getDBContext()
		defer cancel()

		cursor, err := movieCollection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error fetching recommended movies"})
			return
		}
		defer cursor.Close(ctx)

		var recommendedMovies []models.Movie
		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, recommendedMovies)
	}
}

func GetGenres() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := getDBContext()
		defer cancel()

		var genres []models.Genre

		cursor, err := genreCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error fetching movie genres"})
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &genres); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, genres)
	}
}

// Utility functions
// ---------------------------------------------------------------------------------------

const dbTimeout = 100 * time.Second

func getDBContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), dbTimeout)
}

func getReviewRankingWithHugging(adminReview string) (string, int, error) {
	rankings, err := getRankings()
	if err != nil {
		return "", 0, err
	}

	sentiment := ""
	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentiment = sentiment + ranking.RankingName + ","
		}
	}
	sentiment = strings.Trim(sentiment, ",")

	// test development
	err = godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	basePromptTemplate := os.Getenv("BASE_PROMPT_TEMPLATE")
	basePrompt := strings.Replace(basePromptTemplate, "{rankings}", sentiment, 1)
	fullPrompt := basePrompt + adminReview

	var response string

	hfToken := os.Getenv("HUGGING_FACE_HUB_TOKEN")
	hfModel := os.Getenv("HF_MODEL")

	if hfToken == "" {
		return "", 0, errors.New("HUGGING_FACE_HUB_TOKEN not set")
	}
	if hfModel == "" {
		return "", 0, errors.New("HF_MODEL not set")
	}

	response, err = callHuggingFaceAPI(hfToken, hfModel, fullPrompt)
	if err != nil {
		return "", 0, err
	}

	rankVal := 0
	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}

	return response, rankVal, nil
}

func callHuggingFaceAPI(token, model, prompt string) (string, error) {
	url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", model)

	payload := map[string]interface{}{
		"inputs": prompt,
		"parameters": map[string]interface{}{
			"max_new_tokens": 50,
			"temperature":    0.1,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HuggingFace API error: %s - %s", resp.Status, string(body))
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result) == 0 {
		return "", errors.New("empty response from HuggingFace")
	}

	generatedText, ok := result[0]["generated_text"].(string)
	if !ok {
		return "", errors.New("invalid response format from HuggingFace")
	}

	generatedText = strings.TrimPrefix(generatedText, prompt)
	generatedText = strings.TrimSpace(generatedText)

	return generatedText, nil
}

func getReviewRankingOpenAi(adminReview string) (string, int, error) {
	rankings, err := getRankings()
	if err != nil {
		return "", 0, err
	}

	sentiment := ""
	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentiment = sentiment + ranking.RankingName + ","
		}
	}
	sentiment = strings.Trim(sentiment, ",")

	// Only for testing development
	err = godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")
	if OpenAiApiKey == "" {
		return "", 0, errors.New("could not read OpenAI Key")
	}

	llm, err := openai.New(openai.WithToken(OpenAiApiKey))
	if err != nil {
		return "", 0, err
	}

	basePromptTemplate := os.Getenv("BASE_PROMPT_TEMPLATE")
	basePrompt := strings.Replace(basePromptTemplate, "{rankings}", sentiment, 1)

	response, err := llm.Call(context.Background(), basePrompt+adminReview)
	if err != nil {
		return "", 0, err
	}

	rankVal := 0
	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}

	return response, rankVal, nil
}

func getRankings() ([]models.Ranking, error) {
	ctx, cancel := getDBContext()
	defer cancel()

	cursor, err := rankingCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rankings []models.Ranking
	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil
}
