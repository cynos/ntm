package topic

import "time"

type Topic struct {
	ID        uint      `gorm:"primaryKey" json:"id,omitempty"`
	Topic     string    `gorm:"index" json:"topic,omitempty"`
	CreatedAt time.Time `gorm:"default:current_timestamp;index" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updated_at,omitempty"`
}

type Filter struct {
	Topic        string `form:"topic"`
	CreatedStart string `form:"created_start"`
	CreatedEnd   string `form:"created_end"`
}

type ContextKey string
