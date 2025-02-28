package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	OrderID   uuid.UUID
	Name      string
	Data      string
	Status    string
}

type Queries struct {
	db map[uuid.UUID]Task
}

func NewMockDB() *Queries {
	return &Queries{
		db: make(map[uuid.UUID]Task),
	}
}

type CreateTaskParams struct {
	OrderID uuid.UUID
	Name    string
	Data    string
	Status  string
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	var newTask Task

	return newTask, nil
}

func (q *Queries) DeleteTaskById(ctx context.Context, id uuid.UUID) error {
	return nil
}

type GetAllTasksParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetAllTasks(ctx context.Context, arg GetAllTasksParams) ([]Task, error) {
	return nil, nil
}

func (q *Queries) GetTaskByID(ctx context.Context, id uuid.UUID) (Task, error) {
	return Task{}, nil
}

func (q *Queries) GetTaskByMonth(ctx context.Context, dollar_1 time.Time) ([]Task, error) {
	return nil, nil
}

type UpdateTaskNameParams struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) UpdateTaskName(ctx context.Context, arg UpdateTaskNameParams) error {
	return nil
}

type UpdateTaskStatusParams struct {
	ID     uuid.UUID
	Status string
}

func (q *Queries) UpdateTaskStatus(ctx context.Context, arg UpdateTaskStatusParams) error {
	return nil
}
