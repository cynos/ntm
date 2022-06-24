package tag

import "time"

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id,omitempty"`
	Tag       string    `gorm:"index" json:"tag,omitempty"`
	CreatedAt time.Time `gorm:"default:current_timestamp;index" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updated_at,omitempty"`
}

type Filter struct {
	Tag          string `form:"tag"`
	CreatedStart string `form:"created_start"`
	CreatedEnd   string `form:"created_end"`
}

type ContextKey string
