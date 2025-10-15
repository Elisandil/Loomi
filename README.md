# Movie Recommendation API

A Go-based REST API for movie management with personalized recommendations and AI-powered sentiment analysis.

## Features

- **JWT Authentication**: Secure user registration and login
- **Movie CRUD**: Complete movie management system
- **Personalized Recommendations**: Movie suggestions based on user's favorite genres
- **AI Sentiment Analysis**: Automatic review classification using HuggingFace or OpenAI
- **Role-Based Access**: User and Admin roles
- **MongoDB**: NoSQL database for data persistence

## Quick Start

### Prerequisites

- Go 1.24.0+
- MongoDB
- HuggingFace token or OpenAI API key

### Environment Variables

Create a `.env` file:

```env
MONGODB_URI=mongodb://localhost:27017
DATABASE_NAME=movie_db
SECRET_KEY=your_jwt_secret
SECRET_REFRESH_KEY=your_refresh_secret
HUGGING_FACE_HUB_TOKEN=your_hf_token
HF_MODEL=meta-llama/Llama-3.2-1B-Instruct
BASE_PROMPT_TEMPLATE=Classify this movie review sentiment into: {rankings}. Review: 
RECOMMENDED_MOVIE_LIMIT=5
```

### Run

```bash
go mod download
go run main.go
```

Server starts on `http://localhost:8080`

## API Endpoints

### Public Routes

- `POST /register` - Create new user account
- `POST /login` - Authenticate and get JWT token

### Protected Routes
*Requires `Authorization: Bearer <token>` header*

- `GET /movies` - Get all movies
- `GET /movie/:imdb_id` - Get single movie
- `POST /add_movie` - Add new movie
- `PUT /update_movie/:imdb_id` - Update movie
- `DELETE /delete_movie/:imdb_id` - Delete movie
- `GET /recommended_movies` - Get personalized recommendations
- `PATCH /update_review/:imdb_id` - Update admin review (Admin only)

## Tech Stack

- **Framework**: Gin
- **Database**: MongoDB
- **Authentication**: JWT
- **AI Integration**: HuggingFace Inference API / OpenAI
- **Password Hashing**: bcrypt

## Project Structure

```
server/
├── controllers/    # Business logic
├── database/       # DB connection
├── middleware/     # Auth middleware
├── models/         # Data models
├── routes/         # Route definitions
├── utils/          # Helper functions
└── main.go         # Entry point
```

## How It Works

1. Users register with their favorite genres
2. Movies are stored with genres and rankings
3. Admins can add reviews that are automatically analyzed by AI
4. System recommends movies matching user preferences, sorted by ranking
5. JWT tokens secure all protected endpoints

## Security

- Passwords hashed with bcrypt
- JWT tokens (24h access, 7d refresh)
- Role-based authorization
- Protected routes with middleware validation

---

*Work in progress*