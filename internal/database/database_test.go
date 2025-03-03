package database_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/WhileCodingDoLearn/gpt_project/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/magiconair/properties/assert"
)

func TestUpdateUser(t *testing.T) {

	err := godotenv.Load("/home/sms_on_ubuntu/go_projects/gpt_project/.env")
	if err != nil {
		t.Fatalf("err loading: %v", err)
	}
	dbUrl := os.Getenv("DB_URL")
	db, errDB := sql.Open("postgres", dbUrl)
	if errDB != nil {
		t.Log(dbUrl)
		t.Fatalf("url %v  err: %v", dbUrl, errDB)
	}
	defer db.Close()
	dq := database.New(db)

	createdTask, errDBC := dq.CreateTask(context.Background(), database.CreateTaskParams{
		OrderID: uuid.New(),
		Name:    "Tim",
		Data:    "blub",
		Status:  "New",
	})

	if errDBC != nil {
		t.Fatalf("url %v  err: %v", dbUrl, errDBC)
	}

	errDBU := dq.UpdateTaskName(context.Background(), database.UpdateTaskNameParams{
		ID:   createdTask.ID,
		Name: "Kim",
	})
	if errDBU != nil {
		t.Fatal(errDBU)
	}

	updatedTask, errDBG := dq.GetTaskByID(context.Background(), createdTask.ID)
	if errDBG != nil {
		t.Fatal(errDBG)
	}
	assert.Equal(t, updatedTask.Name, "Kim", fmt.Sprintf("Name %v = id %v - should not be %v = id %v\n", createdTask.Name, createdTask.ID, updatedTask.Name, updatedTask.ID))
	//t.Log(updatedTask)
}
