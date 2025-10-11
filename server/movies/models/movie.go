package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Movie struct {
	ID          bson.ObjectID `bson:"_id" json:"_id"`
	ImdbID      string        `bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title       string        `bson:"title" json:"title" validate:"required,min=2,max=500"`
	PosterPath  string        `bson:"poster_path" json:"poster_path" validate:"required,url"`
	YoutubeID   string        `bson:"youtube_id" json:"youtube_id" validate:"required"`
	Genre       []Genre       `bson:"genre" json:"genre" validate:"required,dive"`
	AdminReview string        `bson:"admin_review" json:"admin_review" validate:"required"`
	Ranking     Ranking       `bson:"ranking" json:"ranking" validate:"required"`
}
