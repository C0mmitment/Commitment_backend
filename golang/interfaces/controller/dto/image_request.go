package dto

type ImageRequest struct {
	Category    string `json:"category"`
	Base64Image string `json:"image_data_base64"`
	MimeType    string `json:"mime_type"`
}
