package dto

import "github.com/google/uuid"

type AddLocationRequest struct {
	UserId uuid.UUID `json:"user_uuid"`
	Lat    float64   `json:"latitude"`
	Lng    float64   `json:"longnitude"`
	Geo    string    `json:"geohash"`
}
