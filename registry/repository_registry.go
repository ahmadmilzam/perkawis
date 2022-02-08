package registry

import (
	"meliodas/config"
	"meliodas/src/model"
	"meliodas/src/repository"
	"sync"
)

type repositoryRegistry struct {
	cfg config.Config
}

type RepositoryRegistry interface {
	Activity() model.ActivityRepository
	Todo() model.TodoRepository
}

func NewRepositoryRegistry(cfg config.Config) RepositoryRegistry {
	var repoRegistry RepositoryRegistry
	var loadonce sync.Once

	loadonce.Do(func() {
		repoRegistry = &repositoryRegistry{
			cfg: cfg,
		}
	})

	return repoRegistry
}

func (r repositoryRegistry) Activity() model.ActivityRepository {
	var activityRepository model.ActivityRepository
	var loadonce sync.Once

	loadonce.Do(func() {
		activityRepository = repository.NewActivityRepository(r.cfg)
	})

	return activityRepository
}

func (r repositoryRegistry) Todo() model.TodoRepository {
	var todoRepository model.TodoRepository
	var loadonce sync.Once

	loadonce.Do(func() {
		todoRepository = repository.NewTodoRepository(r.cfg)
	})

	return todoRepository
}
