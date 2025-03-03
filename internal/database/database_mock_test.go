package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/magiconair/properties/assert"
)

func initUsers() IQueries {
	mc := MockQueries{}
	mc.db = make(map[uuid.UUID]Task, 0)
	users := []CreateTaskParams{
		CreateTaskParams{
			OrderID: uuid.New(),
			Name:    "aaaa",
			Data:    "987654",
			Status:  "New",
		},
		CreateTaskParams{
			OrderID: uuid.New(),
			Name:    "abbb",
			Data:    "987654",
			Status:  "New",
		},
		CreateTaskParams{
			OrderID: uuid.New(),
			Name:    "bbbbb",
			Data:    "987654",
			Status:  "New",
		},
		CreateTaskParams{
			OrderID: uuid.New(),
			Name:    "cccc",
			Data:    "987654",
			Status:  "New",
		},
	}

	for _, v := range users {
		id := uuid.New()
		mc.db[id] = Task{ID: id, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), OrderID: v.OrderID, Name: v.Name, Data: v.Data, Status: v.Status}
	}
	return &mc
}

func TestCreateUser(t *testing.T) {
	mc := initUsers()
	id := uuid.New()
	newTask, errCT := mc.CreateTask(context.Background(), CreateTaskParams{
		Name:    "awdawd",
		OrderID: id,
		Data:    "awa",
		Status:  "New",
	})
	if errCT != nil {
		t.Fatal(errCT)
	}
	taskInDB, errGT := mc.GetTaskByID(context.Background(), newTask.ID)
	if errGT != nil {
		t.Fatal(errGT)
	}

	assert.Equal(t, taskInDB, newTask, "Task is not equal")
}

func TestGetTaskByMonth(t *testing.T) {
	mc := initUsers()
	data, errData := mc.GetTaskByMonth(context.Background(), time.Now().AddDate(0, 0, -2).UTC())
	if errData != nil {
		t.Fatal(errData)
	}
	assert.Equal(t, len(data), 4, fmt.Sprint(data))
}

func TestDeleteUser(t *testing.T) {
	mc := initUsers()
	id := uuid.New()
	newTask, errCT := mc.CreateTask(context.Background(), CreateTaskParams{
		Name:    "awdawd",
		OrderID: id,
		Data:    "awa",
		Status:  "New",
	})
	if errCT != nil {
		t.Fatal(errCT)
	}
	errDel := mc.DeleteTaskById(context.Background(), newTask.ID)
	if errDel != nil {
		t.Fatal(errDel)
	}
	taskInDB, _ := mc.GetTaskByID(context.Background(), newTask.ID)
	assert.Equal(t, taskInDB, Task{}, "Task is not empty")
}
