package registry

import (
	"meliodas/config"
	"meliodas/src/model"
	"meliodas/src/usecase"
	"os"
	"strconv"
	"sync"
	"time"
)

type usecaseRegistry struct {
	repo RepositoryRegistry
	cfg  config.Config
}

type UsecaseRegistry interface {
	Activity() model.ActivityUsecase
	Todo() model.TodoUsecase
}

func NewUsecaseRegistry(repo RepositoryRegistry, cfg config.Config) UsecaseRegistry {
	var uc UsecaseRegistry
	var loadonce sync.Once

	loadonce.Do(func() {
		uc = &usecaseRegistry{
			repo: repo,
			cfg:  cfg,
		}
	})

	return uc
}

func (u usecaseRegistry) Activity() model.ActivityUsecase {
	var uc model.ActivityUsecase
	var loadonce sync.Once

	timeout, err := strconv.Atoi(os.Getenv("CTX_TIMEOUT"))
	if err != nil {
		return nil
	}

	timeoutCtx := time.Duration(timeout) * time.Second

	loadonce.Do(func() {
		uc = usecase.NewActivityUsecase(u.repo.Activity(), timeoutCtx)
	})

	return uc

}
func (u usecaseRegistry) Todo() model.TodoUsecase {
	var uc model.TodoUsecase
	var loadonce sync.Once

	timeout, err := strconv.Atoi(os.Getenv("CTX_TIMEOUT"))
	if err != nil {
		return nil
	}

	timeoutCtx := time.Duration(timeout) * time.Second

	loadonce.Do(func() {
		uc = usecase.NewTodoUsecase(u.repo.Todo(), timeoutCtx)
	})

	return uc
}
