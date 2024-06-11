package models

import "database/sql"

type Storage struct {
	TimestampedModel

	ID      int64          `json:"id"`
	Name    string         `json:"name"`
	Address sql.NullString `json:"address"`
}
