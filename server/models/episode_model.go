package models

type Episode struct {
	EpisodeNumber int    `bson:"episode_number" json:"episode_number" validate:"required,min=1"`
	EpisodeTitle  string `bson:"episode_title" json:"episode_title" validate:"required,min=1,max=200"`
	Duration      int    `bson:"duration" json:"duration" validate:"required,min=1"`
	AirDate       string `bson:"air_date" json:"air_date" validate:"required"`
	Synopsis      string `bson:"synopsis" json:"synopsis" validate:"max=1000"`
}
