package tag

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context) ([]Tag, error)
	GetByID(ctx context.Context, id int) (Tag, error)
	Upsert(ctx context.Context, model Tag) (Tag, error)
	DeleteByID(ctx context.Context, id int) error
	GetDB() *gorm.DB
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetAll(ctx context.Context) (res []Tag, err error) {
	filter := ctx.Value(ContextKey("tags_filter")).(Filter)

	var dateValid bool
	_, a := time.Parse("2006-01-02", filter.CreatedStart)
	_, b := time.Parse("2006-01-02", filter.CreatedStart)
	if a == nil && b == nil {
		dateValid = true
	}

	var result *gorm.DB

	if filter.Tag != "" {
		if dateValid {
			result = r.db.Find(&res, "tag = ? and to_char(created_at, 'YYYY-MM-DD') between ? and ?", filter.Tag, filter.CreatedStart, filter.CreatedEnd)
			return res, result.Error
		}

		result = r.db.Find(&res, "tag = ?", filter.Tag)
		return res, result.Error
	}

	if dateValid {
		result = r.db.Find(&res, "to_char(created_at, 'YYYY-MM-DD') between ? and ?", filter.CreatedStart, filter.CreatedEnd)
		return res, result.Error
	}

	result = r.db.Find(&res)
	return res, result.Error
}

func (r *repository) GetByID(ctx context.Context, id int) (res Tag, err error) {
	result := r.db.First(&res, id)
	return res, result.Error
}

func (r *repository) Upsert(ctx context.Context, model Tag) (res Tag, err error) {
	result := r.db.Save(&model)
	return model, result.Error
}

func (r *repository) DeleteByID(ctx context.Context, id int) error {
	result := r.db.Unscoped().Delete(&Tag{ID: uint(id)})
	return result.Error
}

func (r *repository) GetDB() *gorm.DB {
	return r.db
}
