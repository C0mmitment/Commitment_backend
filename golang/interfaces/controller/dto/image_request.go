package dto

import "github.com/google/uuid"

type ImageRequest struct {
	UserId      uuid.UUID `json:"user_uuid"`
	Category    string    `json:"category"`
	Base64Image string    `json:"image_data_base64"`
	MimeType    string    `json:"mime_type"`
	Lat         float64  `json:"latitude"`
	Lng         float64  `json:"longitude"`
	Geo         string    `json:"geohash"`
	SaveLoc     bool      `json:"save_loc"`
}
