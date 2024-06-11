package models

type User struct {
	TimestampedModel

	ID         int64  `json:"id"`
	Username   string `json:"username"`
	ChatID     int64  `json:"chat_id"`
	AccessHash *int64 `json:"access_hash"`
}
