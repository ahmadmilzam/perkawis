package delivery

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"meliodas/helper"
	"meliodas/src/model"
	"meliodas/src/request"
	"strconv"
	"time"
)

type activityDelivery struct {
	actUsecase model.ActivityUsecase
}

type ActivityDelivery interface {
	Mount(group *echo.Group)
}

func NewActivityHttpDelivery(act model.ActivityUsecase) ActivityDelivery {
	return &activityDelivery{actUsecase: act}
}

func (d *activityDelivery) Mount(group *echo.Group) {
	group.GET("", d.FetchHandler)
	group.POST("", d.StoreActivityHandler)
	group.GET("/:id", d.DetailActivityHandler)
	group.DELETE("/:id", d.DeleteActivityHandler)
	group.PATCH("/:id", d.EdiActivityHandler)
}

func (d *activityDelivery) FetchHandler(e echo.Context) error {
	ctx := e.Request().Context()

	key := "activities"
	activities, err := cache.Get(key)
	if err == notFound {
		activities, _, err := d.actUsecase.Fetch(ctx)
		if err != nil {
			return helper.JsonERROR(e, err)
		}
		go cache.SetWithTTL(key, activities, time.Hour)
		return helper.JsonSUCCESS(e, activities)
	}

	return helper.JsonSUCCESS(e, activities)
}

func (d *activityDelivery) StoreActivityHandler(e echo.Context) error {
	ctx := e.Request().Context()
	var req request.CreateActivity

	if err := e.Bind(&req); err != nil {
		return helper.JsonERROR(e, err)
	}

	if req.Title == "" {
		return helper.JsonValidationError(e, "title cannot be null")
	}

	resp, _, err := d.actUsecase.Create(ctx, req)
	if err != nil {
		return helper.JsonERROR(e, err)
	}

	if resp != nil {
		key := fmt.Sprintf("activity-id-%d", resp.ID)
		go cache.SetWithTTL(key, resp, time.Hour)
		go cache.Remove("activities")
	}

	return helper.JsonCreated(e, resp)
}

func (d *activityDelivery) DeleteActivityHandler(e echo.Context) error {
	ctx := e.Request().Context()

	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return helper.JsonERROR(e, err)
	}

	_, err = d.actUsecase.Remove(ctx, id)
	if err != nil {
		return helper.JsonNotFound(e, fmt.Sprintf("Activity with ID %d Not Found", id))
	}

	key := fmt.Sprintf("activity-id-%d", id)
	go cache.Remove(key)
	go cache.Remove("activities")

	return helper.JsonSuccessDelete(e)
}

func (d *activityDelivery) EdiActivityHandler(e echo.Context) error {
	ctx := e.Request().Context()
	var req request.UpdateActivity

	if err := e.Bind(&req); err != nil {
		return helper.JsonERROR(e, err)
	}

	if req.Title == "" {
		return helper.JsonValidationError(e, "title cannot be null")
	}

	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return helper.JsonERROR(e, err)
	}

	resp, _, err := d.actUsecase.Update(ctx, id, req)
	if err != nil {
		return helper.JsonNotFound(e, fmt.Sprintf("Activity with ID %d Not Found", id))
	}

	if resp != nil {
		key := fmt.Sprintf("activity-id-%d", id)
		go cache.SetWithTTL(key, resp, time.Hour)
		go cache.Remove("activities")
	}

	return helper.JsonSUCCESS(e, resp)
}

func (d *activityDelivery) DetailActivityHandler(e echo.Context) error {
	ctx := e.Request().Context()

	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return helper.JsonERROR(e, err)
	}

	key := fmt.Sprintf("activity-id-%d", id)
	activity, err := cache.Get(key)

	if err == notFound {
		activity, _, err := d.actUsecase.GetById(ctx, id)
		if err != nil {
			return helper.JsonNotFound(e, fmt.Sprintf("Activity with ID %d Not Found", id))
		}
		go cache.SetWithTTL(key, activity, time.Hour)
		return helper.JsonSUCCESS(e, activity)
	}

	return helper.JsonSUCCESS(e, activity)
}
