package usecase

import (
	"context"
	"perkawis/constant"
	"perkawis/src/model"
	"perkawis/src/request"
	"time"
)

type activity struct {
	activityRepo   model.ActivityRepository
	contextTimeout time.Duration
}

func (a *activity) Create(ctx context.Context, req request.CreateActivity) (*model.Activity, int, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	gatra, err := a.activityRepo.InsertActivity(ctx, &model.Activity{
		Email: req.Email,
		Title: req.Title,
	})

	if err != nil {
		return nil, constant.CodeInternalServerError, err
	}

	return gatra, constant.CodeSuccess, nil
}

func (a *activity) Update(ctx context.Context, id int, req request.UpdateActivity) (*model.Activity, int, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	updateActivity, err := a.activityRepo.UpdateActivity(ctx, &model.Activity{
		ID:    id,
		Title: req.Title,
		Email: req.Email,
	})

	if err != nil {
		return nil, constant.CodeInternalServerError, err
	}

	return updateActivity, constant.CodeSuccess, nil

}

func (a *activity) Remove(ctx context.Context, id int) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	err := a.activityRepo.RemoveActivity(ctx, id)
	if err != nil {
		return constant.CodeInternalServerError, err
	}

	return constant.CodeSuccess, nil
}

func (a *activity) Fetch(ctx context.Context) ([]*model.Activity, int, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	fetchActivity, err := a.activityRepo.FetchActivity(ctx)
	if err != nil {
		return nil, constant.CodeInternalServerError, err
	}

	return fetchActivity, constant.CodeSuccess, nil
}

func (a *activity) GetById(ctx context.Context, id int) (*model.Activity, int, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	activity, err := a.activityRepo.GetActivityByID(ctx, id)
	if err != nil {
		return nil, constant.CodeInternalServerError, err
	}

	return activity, constant.CodeSuccess, nil
}

func NewActivityUsecase(ar model.ActivityRepository, contextTimeout time.Duration) model.ActivityUsecase {
	return &activity{activityRepo: ar,
		contextTimeout: contextTimeout,
	}
}
