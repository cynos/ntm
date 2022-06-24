package topic

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context) ([]Topic, error)
	GetByID(ctx context.Context, id int) (Topic, error)
	Upsert(ctx context.Context, model Topic) (Topic, error)
	DeleteByID(ctx context.Context, id int) error
	GetDB() *gorm.DB
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetAll(ctx context.Context) (res []Topic, err error) {
	filter := ctx.Value(ContextKey("topics_filter")).(Filter)

	var dateValid bool
	_, a := time.Parse("2006-01-02", filter.CreatedStart)
	_, b := time.Parse("2006-01-02", filter.CreatedStart)
	if a == nil && b == nil {
		dateValid = true
	}

	var result *gorm.DB

	if filter.Topic != "" {
		if dateValid {
			result = r.db.Find(&res, "topic = ? and to_char(created_at, 'YYYY-MM-DD') between ? and ?", filter.Topic, filter.CreatedStart, filter.CreatedEnd)
			return res, result.Error
		}

		result = r.db.Find(&res, "topic = ?", filter.Topic)
		return res, result.Error
	}

	if dateValid {
		result = r.db.Find(&res, "to_char(created_at, 'YYYY-MM-DD') between ? and ?", filter.CreatedStart, filter.CreatedEnd)
		return res, result.Error
	}

	result = r.db.Find(&res)
	return res, result.Error
}

func (r *repository) GetByID(ctx context.Context, id int) (res Topic, err error) {
	result := r.db.First(&res, id)
	return res, result.Error
}

func (r *repository) Upsert(ctx context.Context, model Topic) (res Topic, err error) {
	result := r.db.Save(&model)
	return model, result.Error
}

func (r *repository) DeleteByID(ctx context.Context, id int) error {
	result := r.db.Unscoped().Delete(&Topic{ID: uint(id)})
	return result.Error
}

func (r *repository) GetDB() *gorm.DB {
	return r.db
}
