package models

type Order struct {
	TimestampedModel

	ID        int64 `json:"id"`
	UserID    int64 `json:"user_id"`
	ItemID    int64 `json:"item_id"`
	StorageID int64 `json:"storage_id"`
	Active    bool  `json:"active"`
	Code      int64 `json:"code"`
}
