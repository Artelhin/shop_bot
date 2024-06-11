package models

type ItemToStorage struct {
	ItemID    int64 `json:"item_id"`
	StorageID int64 `json:"storage_id"`
	Count     int64 `json:"count"`
}
