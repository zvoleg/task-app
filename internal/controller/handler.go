package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"task-app/internal/controller/codegen/httpapi"
	"task-app/internal/entity"
	"task-app/internal/usecase"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func NewHandler(si httpapi.ServerInterface, r chi.Router) http.Handler {
	options := httpapi.ChiServerOptions{
		BaseRouter:       r,
		ErrorHandlerFunc: ErrorHandler,
	}
	return httpapi.HandlerWithOptions(si, options)
}

type UseCase interface {
	CreateTask(entity.Task) (*entity.Task, error)
	DeleteTask(uuid.UUID) error
	GetTask(uuid.UUID) (*entity.Task, error)
	GetTaskList(entity.FilterParameters) ([]entity.Task, error)
	UpdateTask(uuid.UUID, entity.Task) (*entity.Task, error)
}

type TaskHandler struct {
	uc UseCase
}

func NewServer(uc UseCase) httpapi.ServerInterface {
	return TaskHandler{
		uc: uc,
	}
}

// CreateTask implements codegen.ServerInterface
func (th TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var body httpapi.CreateTaskBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	task := toEntityTask(body)
	createdTask, err := th.uc.CreateTask(task)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	httpapiTask := fromEntityTask(*createdTask)
	responseBody, err := json.Marshal(httpapiTask)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(responseBody)
}

// DeleteTask implements codegen.ServerInterface
func (th TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request, taskId string) {
	id, err := uuid.Parse(taskId)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	err = th.uc.DeleteTask(id)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetTask implements codegen.ServerInterface
func (th TaskHandler) GetTask(w http.ResponseWriter, r *http.Request, taskId string) {
	id, err := uuid.Parse(taskId)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	task, err := th.uc.GetTask(id)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	httpapiTask := fromEntityTask(*task)
	responseBody, err := json.Marshal(httpapiTask)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

// GetTaskList implements codegen.ServerInterface
func (th TaskHandler) GetTaskList(w http.ResponseWriter, r *http.Request, requestParameters httpapi.GetTaskListParams) {
	params := toEntityFilterParameters(*requestParameters.FilterParameters)
	taskList, err := th.uc.GetTaskList(params)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	httpapiTaskList := fromEntityTaskList(taskList)
	responseBody, err := json.Marshal(httpapiTaskList)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

// UpdateTask implements codegen.ServerInterface
func (th TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request, taskId string) {
	id, err := uuid.Parse(taskId)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	var body httpapi.CreateTaskBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	task := toEntityTask(body)

	updatedTask, err := th.uc.UpdateTask(id, task)

	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	httpapiTask := fromEntityTask(*updatedTask)
	responseBody, err := json.Marshal(httpapiTask)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(responseBody)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	var invalidParam *httpapi.InvalidParamFormatError

	var statusCode int
	if errors.Is(err, usecase.ErrUcInternal) {
		statusCode = http.StatusInternalServerError
	} else if errors.Is(err, usecase.ErrUcNotFound) {
		statusCode = http.StatusNotFound
	} else if errors.Is(err, usecase.ErrUcAlreadyExists) {
		statusCode = http.StatusConflict
	} else if errors.As(err, &invalidParam) {
		statusCode = http.StatusBadRequest
	} else if uuid.IsInvalidLengthError(err) {
		statusCode = http.StatusBadRequest
	} else {
		statusCode = http.StatusInternalServerError
	}

	message := err.Error()
	errorMessage := httpapi.Error{
		Message: &message,
	}

	body, err := json.Marshal(errorMessage)
	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte("Internal error with error message generation"))
	}

	w.WriteHeader(statusCode)
	w.Write(body)
}

func toEntityStatus(status httpapi.CreateTaskBodyStatus) entity.TaskStatus {
	switch status {
	case httpapi.CreateTaskBodyStatusOpened:
		return entity.Opened
	case httpapi.CreateTaskBodyStatusClosed:
		return entity.Closed
	default:
		return entity.Undefined
	}
}

func fromEntityStatus(status entity.TaskStatus) httpapi.TaskStatus {
	switch status {
	case entity.Opened:
		return httpapi.TaskStatusOpened
	case entity.Closed:
		return httpapi.TaskStatusClosed
	default:
		return httpapi.TaskStatusClosed
	}
}

func toEntityTask(t httpapi.CreateTaskBody) entity.Task {
	return entity.Task{
		Id:          nil,
		Name:        t.Name,
		Description: *t.Description,
		Status:      toEntityStatus(t.Status),
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}
}

func fromEntityTask(t entity.Task) httpapi.Task {
	id := t.Id.String()
	return httpapi.Task{
		Id:          id,
		Name:        t.Name,
		Description: t.Description,
		Status:      fromEntityStatus(t.Status),
		CreatedAt:   t.CreatedAt,
		UpdateAt:    t.UpdatedAt,
	}
}

func toEntityFilterParameters(params httpapi.FilterParameters) entity.FilterParameters {
	return entity.FilterParameters{
		Name:          params.Name,
		CreatedAtFrom: params.CreatedAtFrom,
		CreatedAtTo:   params.CreatedAtTo,
		UpdatedAtFrom: params.UpdatedAtFrom,
		UpdatedAtTo:   params.UpdatedAtTo,
	}
}

func fromEntityTaskList(taskList []entity.Task) httpapi.TaskList {
	amount := len(taskList)
	httpapiTaskList := make([]httpapi.Task, amount)

	for _, task := range taskList {
		httpapiTaskList = append(httpapiTaskList, fromEntityTask(task))
	}

	return httpapi.TaskList{
		Amount:   amount,
		Entities: httpapiTaskList,
	}
}
