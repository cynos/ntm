package news

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context) ([]News, error)
	GetByID(ctx context.Context, id int) (News, error)
	Upsert(ctx context.Context, model News) (News, error)
	GetDB() *gorm.DB
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetAll(ctx context.Context) (res []News, err error) {
	exec := r.db.Preload("Topic").Preload("Tags")
	filter := ctx.Value(ContextKey("news_filter")).(Filter)

	var dateValid bool
	_, a := time.Parse("2006-01-02", filter.CreatedStart)
	_, b := time.Parse("2006-01-02", filter.CreatedStart)
	if a == nil && b == nil {
		dateValid = true
	}

	if dateValid {
		exec.Where("to_char(created_at, 'YYYY-MM-DD') between ? and ?", filter.CreatedStart, filter.CreatedEnd)
	}

	if filter.Status != "" {
		if filter.Topic != 0 {
			err = exec.Find(&res, "status = ? and topic_id = ?", filter.Status, filter.Topic).Error
		} else {
			err = exec.Find(&res, "status = ?", filter.Status).Error
		}
	} else {
		if filter.Topic != 0 {
			err = exec.Find(&res, "topic_id = ?", filter.Topic).Error
		} else {
			err = exec.Find(&res).Error
		}
	}

	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return res, fmt.Errorf("record not found")
		}
		return res, err
	}

	return res, nil
}

func (r *repository) GetByID(ctx context.Context, id int) (res News, err error) {
	result := r.db.Where("id = ?", id).Preload("Topic").Preload("Tags").Find(&res)
	return res, result.Error
}

func (r *repository) Upsert(ctx context.Context, model News) (res News, err error) {
	result := r.db.Save(&model)
	return model, result.Error
}

func (r *repository) GetDB() *gorm.DB {
	return r.db
}
