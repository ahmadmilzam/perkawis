package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"meliodas/config"
	"meliodas/config/database"
	"meliodas/constant"
	"meliodas/helper"
	"meliodas/src/model"
)

type activity struct {
	cfg config.Config
	DB  database.GormDatabase
}

func (ac *activity) InsertActivity(ctx context.Context, a *model.Activity) (*model.Activity, error) {
	if err := ac.DB.Master().WithContext(ctx).Create(&a).Error; err != nil {
		return nil, helper.NewErrorMsg(constant.CodeErrQueryDB, err)
	}
	return a, nil
}

func (ac *activity) GetActivityByID(ctx context.Context, id int) (*model.Activity, error) {
	act := new(model.Activity)
	if err := ac.DB.Master().WithContext(ctx).
		Select("id", "email", "title", "created_at", "updated_at", "deleted_at").
		Where("id = ?", id).
		First(act).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewErrorRecordNotFound(constant.CodeErrDataNotFound, err)
		}
		return nil, helper.NewErrorMsg(constant.CodeErrQueryDB, err)
	}

	return act, nil
}

func (ac *activity) UpdateActivity(ctx context.Context, a *model.Activity) (*model.Activity, error) {
	res := ac.DB.Master().WithContext(ctx).
		Model(&model.Activity{ID: a.ID}).Updates(a).Find(a)

	if res.Error != nil {
		return nil, helper.NewErrorMsg(constant.CodeErrQueryDB, res.Error)
	}

	if res.RowsAffected == 0 {
		return nil, helper.NewErrorRecordNotFound(constant.CodeErrDataNotFound, fmt.Errorf("0 rows effected"))
	}

	return a, nil
}

func (ac *activity) RemoveActivity(ctx context.Context, id int) error {
	res := ac.DB.Master().WithContext(ctx).Delete(&model.Activity{ID: id})
	if res.Error != nil {
		return helper.NewErrorMsg(constant.CodeErrQueryDB, res.Error)
	}
	if res.RowsAffected == 0 {
		return helper.NewErrorMsg(constant.CodeErrDataNotFound, fmt.Errorf("0 rows effected"))
	}
	return nil
}

func (ac *activity) FetchActivity(ctx context.Context) ([]*model.Activity, error) {
	var data []*model.Activity
	if err := ac.DB.Master().WithContext(ctx).
		Select("id", "email", "title", "email", "created_at", "updated_at").
		Find(&data).Error; err != nil {
		return nil, helper.NewErrorMsg(constant.CodeErrQueryDB, err)
	}
	return data, nil
}

func NewActivityRepository(cfg config.Config) model.ActivityRepository {
	return &activity{cfg: cfg, DB: cfg.DB()}
}
