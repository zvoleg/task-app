package repo_interface

import (
	"errors"
	"fmt"
	"task-app/internal/entity"

	"github.com/google/uuid"
)

type Repository interface {
	CreateTask(entity.Task) (*entity.Task, error)
	DeleteTask(uuid.UUID) error
	GetTask(uuid.UUID) (*entity.Task, error)
	GetTaskList(entity.FilterParameters) ([]entity.Task, error)
	UpdateTask(uuid.UUID, entity.Task) (*entity.Task, error)
}

var (
	ErrRepoAlreadyExists = errors.New("entity alredy excists")
	ErrRepoNotFound      = errors.New("entity not found")
)

type ErrRepoInternal struct {
	Err error
}

func (e ErrRepoInternal) Error() string {
	return fmt.Sprintf("repo: internal repo error: %s", e.Err.Error())
}

func (e ErrRepoInternal) Unwrap() error {
	return e.Err
}
