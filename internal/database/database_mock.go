package database

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MockQueries struct {
	db   map[uuid.UUID]Task
	lock sync.Mutex
}

func NewMockDB() IQueries {
	return &MockQueries{
		db: make(map[uuid.UUID]Task),
	}
}

func (q *MockQueries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {

	if arg.Status != "New" || arg.Status != "In process" || arg.Status != "Completed" {
		return Task{}, nil
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	newTask := Task{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		OrderID:   arg.OrderID,
		Name:      arg.Name,
		Data:      arg.Data,
		Status:    arg.Status,
	}

	q.db[newTask.ID] = newTask

	return newTask, nil
}

func (q *MockQueries) DeleteTaskById(ctx context.Context, id uuid.UUID) error {
	q.lock.Lock()
	defer q.lock.Unlock()
	delete(q.db, id)
	return nil
}

func (q *MockQueries) DeleteAllTasks(ctx context.Context) error {
	q.lock.Lock()
	defer q.lock.Unlock()
	clear(q.db)
	return nil
}

func (q *MockQueries) GetAllTasks(ctx context.Context, arg GetAllTasksParams) ([]Task, error) {
	return slices.Collect(maps.Values(q.db)), nil
}

func (q *MockQueries) GetTaskByID(ctx context.Context, id uuid.UUID) (Task, error) {
	t, ok := q.db[id]
	if !ok {
		return Task{}, fmt.Errorf("id not found in table")
	}
	return t, nil
}

func (q *MockQueries) GetTaskByMonth(ctx context.Context, dollar_1 time.Time) ([]Task, error) {
	result := make([]Task, 0)

	rangeEnd := dollar_1.AddDate(0, 1, 0)
	for _, value := range q.db {
		if value.CreatedAt.After(dollar_1) && value.CreatedAt.Before(rangeEnd) {
			result = append(result, value)
		}
	}

	return result, nil
}

func (q *MockQueries) UpdateTaskName(ctx context.Context, arg UpdateTaskNameParams) error {
	t, ok := q.db[arg.ID]
	if !ok {
		return fmt.Errorf("id not found in table")
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	t.UpdatedAt = time.Now().UTC()
	t.Name = arg.Name
	q.db[t.ID] = t

	return nil
}

func (q *MockQueries) UpdateTaskStatus(ctx context.Context, arg UpdateTaskStatusParams) error {
	t, ok := q.db[arg.ID]
	if !ok {
		return fmt.Errorf("id not found in table")
	}
	q.lock.Lock()
	defer q.lock.Unlock()

	t.UpdatedAt = time.Now().UTC()
	t.Status = arg.Status
	q.db[t.ID] = t
	return nil
}
