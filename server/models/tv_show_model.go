package models

import "go.mongodb.org/mongo-driver/v2/bson"

type TVShow struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ImdbID       string        `bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title        string        `bson:"title" json:"title" validate:"required"`
	PosterPath   string        `bson:"poster_path" json:"poster_path" validate:"required,url"`
	TrailerID    string        `bson:"trailer_id" json:"trailer_id" validate:"required"`
	Genre        []Genre       `bson:"genre" json:"genre" validate:"required,dive"`
	AdminReview  string        `bson:"admin_review" json:"admin_review"`
	Ranking      Ranking       `bson:"ranking" json:"ranking" validate:"required"`
	Seasons      []Season      `bson:"seasons" json:"seasons" validate:"required,dive"`
	TotalSeasons int           `bson:"total_seasons" json:"total_seasons" validate:"required,min=1"`
	Status       string        `bson:"status" json:"status" validate:"required,oneof=Ongoing Finished Cancelled"`
	FirstAired   string        `bson:"first_aired" json:"first_aired" validate:"required"`
}
