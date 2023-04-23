package usecase

import (
	"errors"
	"fmt"
	"task-app/internal/entity"

	"github.com/google/uuid"
)

type TaskUseCase struct {
	repo Repository
}

func New(repo Repository) TaskUseCase {
	return TaskUseCase{
		repo: repo,
	}
}

func (uc TaskUseCase) CreateTask(task entity.Task) (*entity.Task, error) {
	taskId := uuid.New()
	task.Id = &taskId

	createdTask, err := uc.repo.CreateTask(task)
	if err != nil {
		return nil, wrapErr(err)
	}

	return createdTask, nil
}

func (uc TaskUseCase) DeleteTask(taskId uuid.UUID) error {
	err := uc.repo.DeleteTask(taskId)
	if err != nil {
		return wrapErr(err)
	}
	return nil
}

func (uc TaskUseCase) GetTask(taskId uuid.UUID) (*entity.Task, error) {
	task, err := uc.repo.GetTask(taskId)
	if err != nil {
		return nil, wrapErr(err)
	}
	return task, nil
}

func (uc TaskUseCase) GetTaskList(params entity.FilterParameters) ([]entity.Task, error) {
	taskList, err := uc.repo.GetTaskList(params)
	if err != nil {
		return nil, wrapErr(err)
	}
	return taskList, nil
}

func (uc TaskUseCase) UpdateTask(taskId uuid.UUID, task entity.Task) (*entity.Task, error) {
	updatedTask, err := uc.repo.UpdateTask(taskId, task)
	if err != nil {
		return nil, wrapErr(err)
	}
	return updatedTask, nil
}

var (
	ErrUcInternal      = errors.New("err UseCase: internal error")
	ErrUcNotFound      = errors.New("err UseCase: entity not found")
	ErrUcAlreadyExists = errors.New("err UseCase: entity already exists")
)

func wrapErr(err error) error {
	var errRepoInternal ErrRepoInternal
	if errors.As(err, &errRepoInternal) {
		return ErrUcInternal
	} else if errors.Is(err, ErrRepoNotFound) {
		return ErrUcNotFound
	} else if errors.Is(err, ErrRepoAlreadyExists) {
		return ErrUcAlreadyExists
	} else {
		return nil
	}
}

//go:generate mockery --name Repository --with-expecter=false
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
