package repository

import (
	"database/sql"
	"errors"
	"task-app/internal/entity"
	"task-app/internal/usecase/repo_interface"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

const TASK_TABLE_NAME = "task"

type TaskRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) repo_interface.Repository {
	return TaskRepository{
		db: db,
	}
}

func (r TaskRepository) CreateTask(task entity.Task) (*entity.Task, error) {
	dbTask := fromEntityTask(task)

	sql, args, err := squirrel.Insert(TASK_TABLE_NAME).
		Columns("taskId", "name", "description", "status", "createdAt", "updatedAt", "isDeleted").
		Values(dbTask.Id, dbTask.Name, dbTask.Description, dbTask.Status, dbTask.CreatedAt, dbTask.UpdatedAt, false).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, repo_interface.InternalRepoErr{Err: err}
	}

	row := r.db.QueryRowx(sql, args...)

	var createdDbTask repoTask
	err = row.StructScan(&createdDbTask)
	if err != nil {
		return nil, wrapErr(err)
	}

	createdTask := toEntityTask(createdDbTask)

	return createdTask, nil
}

func (r TaskRepository) DeleteTask(taskId uuid.UUID) error {
	sql, args, err := squirrel.Update(TASK_TABLE_NAME).
		Set("isDelted", 1).
		Where(squirrel.Eq{"taskId": taskId.String()}).
		ToSql()
	if err != nil {
		return repo_interface.InternalRepoErr{Err: err}
	}

	_, err = r.db.Exec(sql, args...)
	if err != nil {
		return wrapErr(err)
	}

	return nil
}

func (r TaskRepository) GetTask(taskId uuid.UUID) (*entity.Task, error) {
	sql, args, err := squirrel.Select("*").
		From(TASK_TABLE_NAME).
		Where(squirrel.Eq{"taskId": taskId.String()}).
		ToSql()
	if err != nil {
		return nil, repo_interface.InternalRepoErr{Err: err}
	}

	row := r.db.QueryRowx(sql, args...)
	var task repoTask
	err = row.StructScan(&task)
	if err != nil {
		return nil, wrapErr(err)
	}

	if task.IsDeleted == 1 {
		return nil, repo_interface.ErrRepoNotFound
	}

	entityTask := toEntityTask(task)

	return entityTask, nil
}

func (r TaskRepository) GetTaskList(params entity.FilterParameters) ([]entity.Task, error) {
	sqlBuilder := squirrel.Select("*").
		From(TASK_TABLE_NAME).
		Where(squirrel.Eq{"isDeleted": 0})

	if params.Name != nil {
		sqlBuilder = sqlBuilder.Where(squirrel.Like{"name": params.Name})
	}

	if params.CreatedAtFrom != nil {
		createdAtFrom := params.CreatedAtFrom.Unix()
		sqlBuilder = sqlBuilder.Where(squirrel.Gt{"createdAt": createdAtFrom})
	}
	if params.CreatedAtTo != nil {
		createdAtTo := params.CreatedAtTo.Unix()
		sqlBuilder = sqlBuilder.Where(squirrel.Gt{"createdAt": createdAtTo})
	}

	if params.UpdatedAtFrom != nil {
		updatedAtFrom := params.UpdatedAtFrom.Unix()
		sqlBuilder = sqlBuilder.Where(squirrel.Gt{"UpdatedAt": updatedAtFrom})
	}
	if params.UpdatedAtTo != nil {
		UpdatedAtTo := params.UpdatedAtTo.Unix()
		sqlBuilder = sqlBuilder.Where(squirrel.Gt{"UpdatedAt": UpdatedAtTo})
	}

	sql, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, repo_interface.InternalRepoErr{Err: err}
	}

	rows, err := r.db.Queryx(sql, args...)
	if err != nil {
		return nil, wrapErr(err)
	}

	taskList := make([]entity.Task, 0)
	var rTask repoTask
	for rows.Next() {
		rows.StructScan(&rTask)
		task := toEntityTask(rTask)
		taskList = append(taskList, *task)
	}

	return taskList, nil
}

func (r TaskRepository) UpdateTask(taskId uuid.UUID, task entity.Task) (*entity.Task, error) {
	rTask := fromEntityTask(task)

	sql, args, err := squirrel.Update(TASK_TABLE_NAME).
		Set("name", rTask.Name).
		Set("description", rTask.Description).
		Set("status", rTask.Status).
		Set("createdAt", rTask.CreatedAt).
		Set("updatedAt", time.Now().Unix()).
		Where(squirrel.Eq{"taskId": taskId.String()}).
		ToSql()
	if err != nil {
		return nil, repo_interface.InternalRepoErr{Err: err}
	}

	row := r.db.QueryRowx(sql, args...)
	var updatedRepoTask repoTask
	err = row.StructScan(&updatedRepoTask)
	if err != nil {
		return nil, wrapErr(err)
	}

	updatedTask := toEntityTask(updatedRepoTask)

	return updatedTask, nil
}

func fromEntityTask(task entity.Task) repoTask {
	taskId := task.Id.String()
	createdAt := task.CreatedAt.Unix()
	var updatedAt *int64
	if task.UpdatedAt != nil {
		unixTime := task.UpdatedAt.Unix()
		updatedAt = &unixTime
	}
	return repoTask{
		Id:          taskId,
		Name:        task.Name,
		Description: task.Description,
		Status:      int16(task.Status),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		IsDeleted:   0,
	}
}

func toEntityTask(task repoTask) *entity.Task {
	taskId, err := uuid.Parse(task.Id)
	if err != nil {
		return nil
	}
	createdAt := time.Unix(task.CreatedAt, 0)
	var updatedAt *time.Time
	if task.UpdatedAt != nil {
		unixTime := time.Unix(*task.UpdatedAt, 0)
		updatedAt = &unixTime
	}

	return &entity.Task{
		Id:          &taskId,
		Name:        task.Name,
		Description: task.Description,
		Status:      entity.TaskStatus(task.Status),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func wrapErr(err error) error {
	var liteErr sqlite3.Error
	if errors.As(err, &liteErr) {
		switch liteErr.Code {
		case sqlite3.ErrConstraint:
			return repo_interface.ErrRepoAlreadyExists
		default:
			return repo_interface.InternalRepoErr{Err: err}
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		return repo_interface.ErrRepoNotFound
	} else {
		return repo_interface.InternalRepoErr{Err: err}
	}
}
