package usecase

import (
	"errors"
	"fmt"
	"task-app/internal/entity"
	"task-app/internal/usecase/repo_interface"

	"github.com/google/uuid"
)

type TaskUseCase struct {
	repo repo_interface.Repository
}

func New(repo repo_interface.Repository) TaskUseCase {
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
	switch err {
	case repo_interface.ErrRepoInternal:
		return fmt.Errorf("%w", ErrUcInternal)
	case repo_interface.ErrRepoNotFound:
		return fmt.Errorf("%w", ErrUcNotFound)
	case repo_interface.ErrRepoAlreadyExists:
		return fmt.Errorf("%w", ErrUcAlreadyExists)
	}
	return nil
}
