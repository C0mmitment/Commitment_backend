package dto

type ImageRequest struct {
	Base64Image string `json:"image_data_base64"`
	MimeType    string `json:"mime_type"`
}
