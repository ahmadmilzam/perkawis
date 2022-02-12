package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"perkawis/config"
	"perkawis/config/database"
	"perkawis/constant"
	"perkawis/helper"
	"perkawis/src/model"
)

type todo struct {
	cfg config.Config
	DB  database.GormDatabase
}

func NewTodoRepository(cfg config.Config) model.TodoRepository {
	return &todo{cfg: cfg, DB: cfg.DB()}
}

func (t *todo) GetActivityGroup(ctx context.Context) ([]*model.Todo, error) {
	var data []*model.Todo

	if err := t.DB.Master().
		Select("activity_group_id").
		Group("activity_group_id").
		Find(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewErrorRecordNotFound(constant.CodeErrDataNotFound, err)
		}
		return nil, helper.NewErrorMsg(constant.CodeErrQueryDB, err)
	}

	return data, nil
}

func (t *todo) InsertTodos(ctx context.Context, a *model.Todo) (*model.Todo, error) {
	if err := t.DB.Master().WithContext(ctx).Create(&a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func (t *todo) GetTodoByID(ctx context.Context, id uint64) (*model.Todo, error) {
	todo := new(model.Todo)
	if err := t.DB.Master().WithContext(ctx).
		Select("id", "activity_group_id", "title", "priority", "created_at", "updated_at", "deleted_at").
		Where("id = ?", id).
		First(todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewErrorRecordNotFound(constant.CodeErrDataNotFound, err)
		}
		return nil, helper.NewErrorMsg(constant.CodeErrQueryDB, err)
	}

	return todo, nil
}

func (t *todo) UpdateTodo(ctx context.Context, a *model.Todo) (*model.Todo, error) {
	res := t.DB.Master().
		Model(&model.Todo{ID: a.ID}).Updates(a).Find(a)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, helper.NewErrorRecordNotFound(constant.CodeErrDataNotFound, fmt.Errorf("0 rows effected"))
	}

	return a, nil
}

func (t *todo) RemoveTodo(ctx context.Context, id uint64) error {
	res := t.DB.Master().WithContext(ctx).Delete(&model.Todo{ID: id})
	if res.Error != nil {
		return helper.NewErrorMsg(constant.CodeErrQueryDB, res.Error)
	}
	if res.RowsAffected == 0 {
		return helper.NewErrorRecordNotFound(constant.CodeErrDataNotFound, fmt.Errorf("0 rows effected"))
	}
	return nil
}

func (t *todo) FetchTodos(ctx context.Context, ActGroupID uint64) ([]*model.Todo, error) {
	var data []*model.Todo

	query := t.DB.Master().
		Select("id", "activity_group_id", "title", "priority", "created_at", "updated_at", "deleted_at")
	if ActGroupID != 0 {
		query = query.Where("activity_group_id = ? ", ActGroupID)
	}
	result := query.Find(&data)

	if result.Error != nil {
		return nil, helper.NewErrorMsg(constant.CodeErrQueryDB, result.Error)
	}

	return data, nil
}
