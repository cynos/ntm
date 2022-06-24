package tag

import (
	"context"
	"fmt"
)

type UseCase interface {
	FindAll(context context.Context) ([]Tag, error)
	FindByID(context context.Context, id int) (Tag, error)
	Add(context context.Context, model Tag) (Tag, error)
	Update(context context.Context, model Tag, id int) (Tag, error)
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

func (us *useCase) FindAll(context context.Context) (res []Tag, err error) {
	res, err = us.repo.GetAll(context)
	return res, err
}

func (us *useCase) FindByID(context context.Context, id int) (res Tag, err error) {
	res, err = us.repo.GetByID(context, id)
	return res, err
}

func (us *useCase) Add(context context.Context, model Tag) (res Tag, err error) {
	if model.Tag == "" {
		return res, fmt.Errorf("invalid parameters")
	}

	res, err = us.repo.Upsert(context, model)
	return res, err
}

func (us *useCase) Update(context context.Context, model Tag, id int) (res Tag, err error) {
	if model.Tag == "" {
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
	res.Tag = model.Tag
	res, err = us.repo.Upsert(context, res)

	return res, err
}

func (us *useCase) Delete(context context.Context, id int) (err error) {
	err = us.repo.DeleteByID(context, id)
	return err
}
