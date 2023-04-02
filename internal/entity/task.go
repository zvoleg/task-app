package entity

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id          *uuid.UUID
	Name        string
	Description string
	Status      TaskStatus
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type TaskStatus int16

const (
	Undefined TaskStatus = iota
	Opened
	Closed
)
