package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"server/database"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserID    string
	jwt.RegisteredClaims
}

var SecretKey string = os.Getenv("SECRET_KEY")
var SecretRefreshKey string = os.Getenv("SECRET_REFRESH_KEY")
var userCollection *mongo.Collection = database.OpenCollection("users")

func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error) {
	signedToken, err := generateToken(email, firstName, lastName, role, userId, SecretKey,
		24*time.Hour)
	if err != nil {
		return "", "", err // returns an empty token, an empty refresh token and an error
	}

	signedRefreshToken, err := generateToken(email, firstName, lastName, role, userId, SecretRefreshKey,
		168*time.Hour)
	if err != nil {
		return "", "", err // returns an empty token, an empty refresh token and an error
	}

	return signedToken, signedRefreshToken, nil
}

func UpdateAllTokens(userId, token, refreshToken string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"update_at":     updateAt,
		},
	}

	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)
	if err != nil {
		return err
	}

	return nil
}

func GetAccessToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	tokenString := authHeader[len("Bearer "):]
	if tokenString == "" {
		return "", errors.New("bearer token is required")
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

func GetUserIdFromContext(c *gin.Context) (string, error) {
	return getContextValue(c, "userId")
}

func GetRoleFromContext(c *gin.Context) (string, error) {
	return getContextValue(c, "role")
}

// ------------------------------------------------------------------------------------
func generateToken(email, firstName, lastName, role, userId, secret string,
	expirationTime time.Duration) (string, error) {

	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserID:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Loomi",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func getContextValue(c *gin.Context, key string) (string, error) {
	value, exists := c.Get(key)
	if !exists {
		return "", fmt.Errorf("%s does not exist in context", key)
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("unable to retrieve %s: invalid type", key)
	}

	return str, nil
}
