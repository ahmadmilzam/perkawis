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

type todoDelivery struct {
	todoUsecase model.TodoUsecase
}

type TodoDelivery interface {
	Mount(group *echo.Group)
}

func NewTodoHttpDelivery(act model.TodoUsecase) TodoDelivery {
	return &todoDelivery{todoUsecase: act}
}

func (d *todoDelivery) Mount(group *echo.Group) {
	group.GET("", d.FetchHandler)
	group.POST("", d.StoreTodoHandler)
	group.GET("/:id", d.DetailTodoHandler)
	group.DELETE("/:id", d.DeleteTodoHandler)
	group.PATCH("/:id", d.EdiTodoHandler)
}

func (d *todoDelivery) FetchHandler(e echo.Context) error {
	ctx := e.Request().Context()
	var id int
	var err error
	agi := e.QueryParam("activity_group_id")

	if agi != "" {
		id, err = strconv.Atoi(agi)
		if err != nil {
			id = 0
		}
	}

	key := fmt.Sprintf("todos-%d", id)
	todo, err := cache.Get(key)
	if err == notFound {
		todo, _, err := d.todoUsecase.Fetch(ctx, uint64(id))
		if err != nil {
			return helper.JsonERROR(e, err)
		}
		go cache.SetWithTTL(key, todo, 10*time.Minute)
		return helper.JsonSUCCESS(e, todo)
	}

	return helper.JsonSUCCESS(e, todo)
}

func (d *todoDelivery) StoreTodoHandler(e echo.Context) error {
	ctx := e.Request().Context()
	var req request.CreateTodo

	if err := e.Bind(&req); err != nil {
		return helper.JsonERROR(e, err)
	}

	if req.Title == "" {
		return helper.JsonValidationError(e, "title cannot be null")
	}

	if req.ActivityGroupID == 0 {
		return helper.JsonValidationError(e, "activity_group_id cannot be null")
	}

	resp, _, err := d.todoUsecase.Create(ctx, req)
	if err != nil {
		return helper.JsonERROR(e, err)
	}

	go cache.Remove("todos-0")

	return helper.JsonCreated(e, resp)
}

func (d *todoDelivery) DeleteTodoHandler(e echo.Context) error {
	ctx := e.Request().Context()

	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return helper.JsonERROR(e, err)
	}

	_, err = d.todoUsecase.Remove(ctx, uint64(id))
	if err != nil {
		return helper.JsonNotFound(e, fmt.Sprintf("Todo with ID %d Not Found", id))
	}
	key := fmt.Sprintf("todo-id-%d", id)
	go cache.Remove(key)
	go cache.Remove("todos-0")
	return helper.JsonSuccessDelete(e)
}

func (d *todoDelivery) EdiTodoHandler(e echo.Context) error {
	ctx := e.Request().Context()
	var req request.UpdateTodo

	if err := e.Bind(&req); err != nil {
		return helper.JsonERROR(e, err)
	}

	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return helper.JsonERROR(e, err)
	}

	resp, _, err := d.todoUsecase.Update(ctx, uint64(id), req)
	if err != nil {
		return helper.JsonNotFound(e, fmt.Sprintf("Todo with ID %d Not Found", id))
	}

	if resp != nil {
		key := fmt.Sprintf("todo-id-%d", resp.ID)
		go cache.SetWithTTL(key, resp, 10*time.Minute)
		go cache.Remove("todos-0")
	}

	return helper.JsonSUCCESS(e, resp)
}

func (d *todoDelivery) DetailTodoHandler(e echo.Context) error {
	ctx := e.Request().Context()

	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return helper.JsonERROR(e, err)
	}

	key := fmt.Sprintf("todo-id-%d", id)
	todo, err := cache.Get(key)

	if err == notFound {
		todo, _, err := d.todoUsecase.GetById(ctx, uint64(id))
		if err != nil {
			return helper.JsonNotFound(e, fmt.Sprintf("Todo with ID %d Not Found", id))
		}
		go cache.SetWithTTL(key, todo, 10*time.Minute)

		return helper.JsonSUCCESS(e, todo)
	}

	return helper.JsonSUCCESS(e, todo)
}
