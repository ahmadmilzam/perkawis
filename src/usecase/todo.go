package usecase

import (
	"context"
	"meliodas/constant"
	"meliodas/src/model"
	"meliodas/src/request"
	"time"
)

type todo struct {
	todoRepo       model.TodoRepository
	contextTimeout time.Duration
}

func (t *todo) GetActivityGroup(ctx context.Context) ([]*model.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	group, err := t.todoRepo.GetActivityGroup(ctx)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func NewTodoUsecase(ar model.TodoRepository, contextTimeout time.Duration) model.TodoUsecase {
	return &todo{
		todoRepo:       ar,
		contextTimeout: contextTimeout,
	}
}

func (t *todo) Create(ctx context.Context, req request.CreateTodo) (*model.Todo, int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	req.IsActive = true

	if req.Priority == "" {
		req.Priority = "very-high"
	}

	record := &model.Todo{
		ActivityGroupID: req.ActivityGroupID,
		Title:           req.Title,
		IsActive:        &req.IsActive,
		Priority:        req.Priority,
	}

	_, err := t.todoRepo.InsertTodos(ctx, record)
	if err != nil {
		return nil, constant.CodeInternalServerError, err
	}

	return record, constant.CodeSuccess, nil
}

func (t *todo) Update(ctx context.Context, id uint64, req request.UpdateTodo) (*model.Todo, int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	record := &model.Todo{
		ID:              id,
		Title:           req.Title,
		Priority:        req.Priority,
		ActivityGroupID: req.ActivityGroupID,
		IsActive:        &req.IsActive,
	}

	updateActivity, err := t.todoRepo.UpdateTodo(ctx, record)

	if err != nil {
		return nil, constant.CodeInternalServerError, err
	}

	return updateActivity, constant.CodeSuccess, nil
}

func (t *todo) Remove(ctx context.Context, id uint64) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	err := t.todoRepo.RemoveTodo(ctx, id)
	if err != nil {
		return constant.CodeInternalServerError, err
	}

	return constant.CodeSuccess, nil
}

func (t *todo) Fetch(ctx context.Context, id uint64) ([]*model.Todo, int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	fetchActivity, err := t.todoRepo.FetchTodos(ctx, id)
	if err != nil {
		return nil, constant.CodeInternalServerError, err
	}

	return fetchActivity, constant.CodeSuccess, nil
}

func (t *todo) GetById(ctx context.Context, id uint64) (*model.Todo, int, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	activity, err := t.todoRepo.GetTodoByID(ctx, id)
	if err != nil {
		return nil, constant.CodeInternalServerError, err
	}

	return activity, constant.CodeSuccess, nil
}
