# Loomi - Movie & TV Show Recommendation API

A Go-based REST API for movies and TV shows management with personalized recommendations and AI-powered sentiment analysis.

## Tech Stack

- **Framework**: Gin Web Framework
- **Database**: MongoDB (with mongo-driver v2)
- **Authentication**: JWT (golang-jwt/jwt v5)
- **Validation**: go-playground/validator v10
- **Password Hashing**: bcrypt
- **AI Integration**: HuggingFace Inference API (Groq) / OpenAI
- **CORS**: gin-contrib/cors
- **Environment**: godotenv

## Features

- **JWT Authentication**: Secure user registration and login with access/refresh tokens
- **Movie & TV Show CRUD**: Complete management system for movies and TV series
- **Personalized Recommendations**: Content suggestions based on user's favorite genres
- **AI Sentiment Analysis**: Automatic review classification using HuggingFace (Groq) or OpenAI
- **Role-Based Access**: User and Admin roles with different permissions
- **Season & Episode Management**: Complete TV show tracking with seasons and episodes
- **MongoDB**: NoSQL database for data persistence
- **CORS Enabled**: Ready for frontend integration

## Project Structure

```
server/
├── controllers/           # Business logic
│   ├── movie_controller.go
│   ├── tv_show_controller.go
│   ├── user_controller.go
│   └── user_controller_test.go
├── database/             # MongoDB connection
│   └── db_conn.go
├── middleware/           # Auth middleware
│   └── auth_middleware.go
├── models/               # Data models
│   ├── movie_model.go
│   ├── tv_show_model.go
│   ├── season_model.go
│   ├── episode_model.go
│   ├── user_model.go
│   ├── genre_model.go
│   └── ranking_model.go
├── routes/               # Route definitions
│   ├── protected_routes.go
│   └── unprotected_routes.go
├── utils/                # Helper functions
│   └── token_util.go
├── tests/                # HTTP test files
│   └── endpoints/
├── .env                  # Environment variables
└── main.go              # Entry point
```

## Quick Start

### Prerequisites

- Go 1.24.0+
- MongoDB 4.0+
- HuggingFace token (for Groq) or OpenAI API key

### Environment Variables

Create a `.env` file in the `server/` directory:

```env
DATABASE_NAME=Loomi-movies
MONGODB_URI=mongodb://localhost:27017/
SECRET_KEY=your_secret_key
SECRET_REFRESH_KEY=your_refresh_secret_key
BASE_PROMPT_TEMPLATE='Return a response using one of the words: {rankings}.The response should be a single word and should not contain any other text.The response should be based on the following review:'
OPENAI_API_KEY=your_open_ai_key
USE_HUGGING_FAME=true
HUGGING_FACE_HUB_TOKEN=your_hf_token
HF_MODEL=openai/gpt-oss-20b
HF_INFERENCE_PROVIDER=groq
RECOMMENDED_MOVIE_LIMIT=5
```

### Installation & Run

```bash
cd server
go mod download
go run main.go
```

Server starts on `http://localhost:8080`

## API Endpoints

### Public Routes

- `POST /register` - Create new user account
- `POST /login` - Authenticate and get JWT tokens
- `GET /genres` - Get all available genres
- `GET /movies` - Get all movies
- `GET /tv_shows` - Get all TV shows

### Protected Routes
*Requires `Authorization: Bearer <token>` header*

#### Movies

- `GET /movie/:imdb_id` - Get single movie details
- `POST /add_movie` - Add new movie (Admin only)
- `PUT /update_movie/:imdb_id` - Update movie (Admin only)
- `DELETE /delete_movie/:imdb_id` - Delete movie (Admin only)
- `GET /recommended_movies` - Get personalized movie recommendations
- `PATCH /update_review/:imdb_id` - Update movie review with AI analysis (Admin only)

#### TV Shows

- `GET /tv_show/:imdb_id` - Get single TV show details
- `GET /tv_show/:imdb_id/season/:season_number` - Get a TV show season
- `POST /add_tv_show` - Add new TV show (Admin only)
- `PUT /update_tv_show/:imdb_id` - Update TV show (Admin only)
- `POST /tv_show/:imdb_id/add_season` - Add season to TV show (Admin only)
- `DELETE /delete_tv_show/:imdb_id` - Delete TV show (Admin only)
- `PATCH /update_tv_show_review/:imdb_id` - Update TV show review with AI analysis (Admin only)
- `GET /recommended_tv_shows` - Get personalized TV show recommendations

## Data Models

### Movie
```go
{
  "imdb_id": "tt1234567",
  "title": "Movie Title",
  "release_date": "2024-01-01T00:00:00Z",
  "genres": ["Action", "Thriller"],
  "duration": 120,
  "ranking": "Must Watch",
  "admin_review": "Amazing movie!",
  "sentiment": "positive"
}
```

### TV Show
```go
{
  "imdb_id": "tt7654321",
  "title": "Show Title",
  "release_date": "2024-01-01T00:00:00Z",
  "genres": ["Drama", "Sci-Fi"],
  "ranking": "Must Watch",
  "status": "Ongoing", // Ongoing, Finished, Cancelled
  "admin_review": "Great series!",
  "sentiment": "positive",
  "seasons": [
    {
      "season_number": 1,
      "episodes": [
        {
          "episode_number": 1,
          "title": "Pilot",
          "duration": 45,
          "air_date": "2024-01-01T00:00:00Z"
        }
      ]
    }
  ]
}
```

### User
```go
{
  "email": "user@example.com",
  "password": "hashed_password",
  "role": "user", // admin by default
  "favorite_genres": ["Action", "Sci-Fi"]
}
```

## How It Works

1. **User Registration**: Users register with email, password, and favorite genres
2. **Content Management**: Admins can add movies/TV shows with genres and rankings
3. **AI Review Analysis**: When admins add reviews, AI automatically classifies sentiment
4. **Personalized Recommendations**: System suggests content matching user preferences, sorted by ranking
5. **Season Tracking**: TV shows include complete season and episode information
6. **Secure Access**: JWT tokens protect all user-specific and admin endpoints

## AI Sentiment Analysis

The system uses AI to automatically classify admin reviews into predefined rankings:
- Reviews are sent to HuggingFace (Groq) or OpenAI
- AI returns one of the configured rankings
- Sentiment is stored with the content for recommendation algorithms

## Security

- **Password Hashing**: bcrypt
- **JWT Tokens**: 
  - Access tokens
  - Refresh tokens
- **Role-Based Authorization**: Separate permissions for users and admins
- **Protected Routes**: Middleware validates tokens on all protected endpoints
- **CORS**: Configured for secure cross-origin requests

## Testing

HTTP test files are available in `tests/endpoints/` for manual API testing with REST clients.

---

*Built using Go and MongoDB*