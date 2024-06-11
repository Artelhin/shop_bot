package models

type Category struct {
	TimestampedModel

	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ParentID *int64 `json:"parent_id"`
}
