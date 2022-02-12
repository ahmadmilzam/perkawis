package model

import (
	"context"
	"gorm.io/gorm"
	"perkawis/src/request"
	"time"
)

type (
	Todo struct {
		ID              uint64         `json:"id"`
		ActivityGroupID uint64         `json:"activity_group_id" gorm:"index"`
		Title           string         `json:"title"`
		IsActive        *bool          `gorm:"default:1" json:"is_active"`
		Priority        string         `gorm:"default:'very-high'" json:"priority"`
		CreatedAt       time.Time      `json:"created_at"`
		UpdatedAt       time.Time      `json:"updated_at"`
		DeletedAt       gorm.DeletedAt `json:"deleted_at"`
	}

	TodoRepository interface {
		InsertTodos(ctx context.Context, a *Todo) (*Todo, error)
		GetTodoByID(ctx context.Context, id uint64) (*Todo, error)
		UpdateTodo(ctx context.Context, a *Todo) (*Todo, error)
		RemoveTodo(ctx context.Context, id uint64) error
		FetchTodos(ctx context.Context, ActGroupID uint64) ([]*Todo, error)
		GetActivityGroup(ctx context.Context) ([]*Todo, error)
	}

	TodoUsecase interface {
		Create(ctx context.Context, req request.CreateTodo) (*Todo, int, error)
		Update(ctx context.Context, id uint64, req request.UpdateTodo) (*Todo, int, error)
		Remove(ctx context.Context, id uint64) (int, error)
		Fetch(ctx context.Context, id uint64) ([]*Todo, int, error)
		GetById(ctx context.Context, id uint64) (*Todo, int, error)
		GetActivityGroup(ctx context.Context) ([]*Todo, error)
	}
)
