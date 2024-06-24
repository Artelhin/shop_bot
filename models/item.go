package models

import "database/sql"

type Item struct {
	TimestampedModel

	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	CategoryID  int64          `json:"category_id"`
	Image       []byte         `json:"image"` // marshalled models.Image
}

type Image struct {
	MimeType string `json:"mime_type"`
	Bytes    []byte `json:"bytes"`
}
