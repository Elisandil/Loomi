package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"server/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func cleanupTestUser(email string) {
	ctx, cancel := getDBContext()
	defer cancel()
	usersCollection.DeleteOne(ctx, bson.M{"email": email})
}

func TestRegisterUser_Success(t *testing.T) {
	router := setupTestRouter()
	router.POST("/register", RegisterUser())

	testEmail := "test@example.com"
	defer cleanupTestUser(testEmail)

	user := models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     testEmail,
		Password:  "SecurePass123!",
		Role:      "USER",
	}

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	ctx, cancel := getDBContext()
	defer cancel()

	var foundUser models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": testEmail}).Decode(&foundUser)
	assert.NoError(t, err)
	assert.Equal(t, testEmail, foundUser.Email)
	assert.NotEmpty(t, foundUser.UserID)
}

func TestRegisterUser_InvalidInput(t *testing.T) {
	router := setupTestRouter()
	router.POST("/register", RegisterUser())

	invalidJSON := []byte(`{"invalid": json}`)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Contains(t, response, "Error")
}

func TestRegisterUser_ValidationFailed(t *testing.T) {
	router := setupTestRouter()
	router.POST("/register", RegisterUser())

	user := models.User{
		FirstName: "",
		LastName:  "Doe",
		Email:     "invalid-email",
		Password:  "123",
	}

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Contains(t, response, "Error")
	assert.Equal(t, "Validation failed", response["Error"])
}

func TestRegisterUser_DuplicateEmail(t *testing.T) {
	router := setupTestRouter()
	router.POST("/register", RegisterUser())

	testEmail := "duplicate@example.com"
	defer cleanupTestUser(testEmail)

	user := models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     testEmail,
		Password:  "SecurePass123!",
		Role:      "USER",
	}

	jsonData, _ := json.Marshal(user)

	req1, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	req2, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusConflict, w2.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w2.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, "User already exists", response["Error"])
}

func TestLoginUser_Success(t *testing.T) {
	router := setupTestRouter()
	router.POST("/register", RegisterUser())
	router.POST("/login", LoginUser())

	testEmail := "login@example.com"
	testPassword := "SecurePass123!"
	defer cleanupTestUser(testEmail)

	user := models.User{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     testEmail,
		Password:  testPassword,
		Role:      "USER",
		FavouriteGenres: []models.Genre{
			{GenreID: 1, GenreName: "Action"},
			{GenreID: 2, GenreName: "Comedy"},
		},
	}

	jsonData, _ := json.Marshal(user)
	regReq, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	router.ServeHTTP(regW, regReq)

	time.Sleep(100 * time.Millisecond)

	loginData := models.UserLogin{
		Email:    testEmail,
		Password: testPassword,
	}
	loginJSON, _ := json.Marshal(loginData)
	loginReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)

	assert.Equal(t, http.StatusOK, loginW.Code)

	var response models.UserResponse
	err := json.Unmarshal(loginW.Body.Bytes(), &response)
	if err != nil {
		return
	}

	assert.Equal(t, testEmail, response.Email)
	assert.Equal(t, "Jane", response.FirstName)
	assert.Equal(t, "Smith", response.LastName)
	assert.NotEmpty(t, response.Token)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Len(t, response.FavouriteGenres, 2)
	assert.Equal(t, 1, response.FavouriteGenres[0].GenreID)
	assert.Equal(t, "Action", response.FavouriteGenres[0].GenreName)
	assert.Equal(t, 2, response.FavouriteGenres[1].GenreID)
	assert.Equal(t, "Comedy", response.FavouriteGenres[1].GenreName)
}

func TestLoginUser_InvalidEmail(t *testing.T) {
	router := setupTestRouter()
	router.POST("/login", LoginUser())

	loginData := models.UserLogin{
		Email:    "nonexistent@example.com",
		Password: "anyPassword123",
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, "Invalid email/password", response["Error"])
}

func TestLoginUser_IncorrectPassword(t *testing.T) {
	router := setupTestRouter()
	router.POST("/register", RegisterUser())
	router.POST("/login", LoginUser())

	testEmail := "wrongpass@example.com"
	testPassword := "CorrectPass123!"
	defer cleanupTestUser(testEmail)

	user := models.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     testEmail,
		Password:  testPassword,
		Role:      "USER",
	}
	jsonData, _ := json.Marshal(user)
	regReq, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	router.ServeHTTP(regW, regReq)

	time.Sleep(100 * time.Millisecond)

	loginData := models.UserLogin{
		Email:    testEmail,
		Password: "WrongPassword123!",
	}

	loginJSON, _ := json.Marshal(loginData)
	loginReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)

	assert.Equal(t, http.StatusUnauthorized, loginW.Code)

	var response map[string]interface{}
	err := json.Unmarshal(loginW.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, "Incorrect email/password", response["Error"])
}

func TestLoginUser_InvalidInput(t *testing.T) {
	router := setupTestRouter()
	router.POST("/login", LoginUser())

	invalidJSON := []byte(`{"invalid": json}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Contains(t, response, "Error")
}

func TestHashPassword_Success(t *testing.T) {
	password := "MySecurePassword123!"

	hashedPassword, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err)
}

func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "SamePassword123!"

	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2, "hashes should be different to salt")
}
