package topic

import (
	"context"
	"fmt"
)

type UseCase interface {
	FindAll(context context.Context) ([]Topic, error)
	FindByID(context context.Context, id int) (Topic, error)
	Add(context context.Context, model Topic) (Topic, error)
	Update(context context.Context, model Topic, id int) (Topic, error)
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

func (us *useCase) FindAll(context context.Context) (res []Topic, err error) {
	res, err = us.repo.GetAll(context)
	return res, err
}

func (us *useCase) FindByID(context context.Context, id int) (res Topic, err error) {
	res, err = us.repo.GetByID(context, id)
	return res, err
}

func (us *useCase) Add(context context.Context, model Topic) (res Topic, err error) {
	if model.Topic == "" {
		return res, fmt.Errorf("invalid parameters")
	}
	res, err = us.repo.Upsert(context, model)
	return res, err
}

func (us *useCase) Update(context context.Context, model Topic, id int) (res Topic, err error) {
	if model.Topic == "" {
		return res, fmt.Errorf("invalid parameters")
	}

	if id == 0 {
		return res, fmt.Errorf("invalid parameters")
	}

	// get id first
	res, err = us.repo.GetByID(context, id)
	if err != nil {
		return res, err
	}

	// update tag
	res.Topic = model.Topic
	res, err = us.repo.Upsert(context, res)

	return res, err
}

func (us *useCase) Delete(context context.Context, id int) (err error) {
	err = us.repo.DeleteByID(context, id)
	return err
}
