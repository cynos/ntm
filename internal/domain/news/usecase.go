package news

import (
	"context"
	"fmt"
	"time"

	"github.com/ntm/internal/domain/tag"
	"github.com/ntm/internal/domain/topic"
)

type UseCase interface {
	FindAll(context context.Context) ([]News, error)
	FindByID(context context.Context, id int) (News, error)
	Save(context context.Context, model NewsDTO) (News, error)
	Delete(context context.Context, id int) error
}

type useCase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &useCase{
		repo: repo,
	}
}

func (us *useCase) FindAll(context context.Context) (res []News, err error) {
	res, err = us.repo.GetAll(context)
	return res, err
}

func (us *useCase) FindByID(context context.Context, id int) (res News, err error) {
	res, err = us.repo.GetByID(context, id)
	return res, err
}

func (us *useCase) Save(context context.Context, dto NewsDTO) (res News, err error) {
	tagRepo := tag.NewRepository(us.repo.GetDB())
	topicRepo := topic.NewRepository(us.repo.GetDB())

	// status validation
	if dto.Status != string(StatusDraft) && dto.Status != string(StatusPublish) {
		return res, fmt.Errorf("invalid status parameter")
	}

	// get tag list
	var tags []tag.Tag
	for _, v := range dto.Tags {
		if tag, err := tagRepo.GetByID(context, int(v)); err == nil {
			tags = append(tags, tag)
		}
	}

	// get topic
	topic, errTopic := topicRepo.GetByID(context, int(dto.TopicID))
	if err != nil {
		return res, errTopic
	}

	news := News{
		ID:      dto.ID,
		Title:   dto.Title,
		Writer:  dto.Writer,
		Content: dto.Content,
		Status:  dto.Status,
		Tags:    tags,
		Topic:   topic,
	}

	if dto.Status == string(StatusPublish) {
		news.PublishAt = time.Now()
	}

	news, err = us.repo.Upsert(context, news)
	if err != nil {
		return res, err
	}

	return news, nil
}

func (us *useCase) Delete(context context.Context, id int) (err error) {
	// check news is exist or not
	var news News
	news, err = us.repo.GetByID(context, id)
	if err != nil {
		return err
	}

	news.Status = string(StatusDelete)
	news.DeletedAt = time.Now()
	news, err = us.repo.Upsert(context, news)
	if err != nil {
		return err
	}

	return nil
}
