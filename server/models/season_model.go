package models

type Season struct {
	SeasonNumber int       `bson:"season_number" json:"season_number" validate:"required,min=1"`
	Episodes     []Episode `bson:"episodes" json:"episodes" validate:"required,dive"`
}
