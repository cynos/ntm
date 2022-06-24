package news

import (
	"time"

	"github.com/ntm/internal/domain/tag"
	"github.com/ntm/internal/domain/topic"
)

type News struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"not null"`
	Writer    string    `gorm:"not null;type:varchar(100)"`
	Content   string    `gorm:"not null"`
	Status    string    `gorm:"not null;type:varchar(20)"`
	Tags      []tag.Tag `gorm:"many2many:news_tags;"`
	TopicID   uint
	Topic     topic.Topic
	PublishAt time.Time `gorm:"default:current_timestamp;"`
	CreatedAt time.Time `gorm:"default:current_timestamp;index"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
	DeletedAt time.Time `gorm:"default:null"`
}

type NewsDTO struct {
	ID      uint   `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Writer  string `json:"writer,omitempty"`
	Content string `json:"content,omitempty"`
	Status  string `json:"status,omitempty"`
	Tags    []uint `json:"tags,omitempty"`
	TopicID uint   `json:"topic_id,omitempty"`
}

type Status string

const (
	StatusDraft   Status = "draft"
	StatusPublish Status = "publish"
	StatusDelete  Status = "deleted"
)

type Filter struct {
	Status       string `form:"status"`
	Topic        uint   `form:"topic"`
	CreatedStart string `form:"created_start"`
	CreatedEnd   string `form:"created_end"`
}

type ContextKey string
