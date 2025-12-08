package model

import (
	"github.com/86shin/commit_goback/domain/utils"
	"github.com/google/uuid"
)

type AddLocation struct {
	LocationId uuid.UUID
	UserId     uuid.UUID
	Lat        float64
	Lng        float64
	Geo        string
}

func NewAddLocation(user_id uuid.UUID, lat, lng float64, geohash string) (AddLocation, error) {
	id, err := utils.NewGegerateUuid()
	if err != nil {
		return AddLocation{}, err
	}

	return AddLocation{
		LocationId: id,
		UserId:     user_id,
		Lat:        lat,
		Lng:        lng,
		Geo:        geohash,
	}, nil
}
