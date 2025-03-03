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

type StatusType struct {
	NEW       string
	PROCESS   string
	COMPLETED string
}

var status StatusType

var dbMap map[uuid.UUID]Task

func init() {
	dbMap = make(map[uuid.UUID]Task, 0)
	status = StatusType{
		NEW:       "New",
		PROCESS:   "In process",
		COMPLETED: "Completed",
	}
}

func NewMockDB() IQueries {
	return &MockQueries{
		db: dbMap,
	}
}

func (q *MockQueries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	switch arg.Status {
	case status.NEW:
	case status.PROCESS:
	case status.COMPLETED:
		q.lock.Lock()
		defer q.lock.Unlock()
	default:
		return Task{}, fmt.Errorf("status not supprted : '%v'", arg.Status)
	}

	id := uuid.New()

	newTask := Task{
		ID:        id,
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
	data := slices.Collect(maps.Values(q.db))
	result := data[arg.Offset:arg.Limit]
	return result, nil
}

func (q *MockQueries) GetTaskByID(ctx context.Context, id uuid.UUID) (Task, error) {

	t, found := q.db[id]
	if !found {
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
	switch arg.Status {
	case status.NEW:
	case status.PROCESS:
	case status.COMPLETED:
		q.lock.Lock()
		defer q.lock.Unlock()
	default:
		return fmt.Errorf("status not supprted : '%v'", arg.Status)
	}

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
