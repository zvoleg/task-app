package usecase

import (
	"task-app/internal/entity"
	"task-app/internal/usecase/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestCreateTask(t *testing.T) {
	task := entity.Task{
		Id:          nil,
		Name:        "string",
		Description: "string",
		Status:      entity.Opened,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}

	r := mocks.NewRepository(t)
	r.
		On("CreateTask", mock.AnythingOfTypeArgument("entity.Task")).
		Return(func(t entity.Task) (*entity.Task, error) {
			return &t, nil
		})

	uc := New(r)

	createdTask, err := uc.CreateTask(task)
	if err != nil {
		t.Errorf("%v", err)
	}

	if createdTask.Id == task.Id {
		t.Errorf("expected that the created task id is not equalt to 'source' task.Id (%v)", task.Id)
	}
}
