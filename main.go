package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/WhileCodingDoLearn/gpt_project/internal/background"
	"github.com/WhileCodingDoLearn/gpt_project/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

/*
goose postgres postgres://postgres:postgres@localhost:5432/task_db up/down

*/

var VERSION string

func init() {
	VERSION = "/V1"
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	portFromEnv := os.Getenv("PORT")
	if len(portFromEnv) == 0 {
		portFromEnv = "8080"
	}

	userApiKey := os.Getenv("API_KEY_USER")
	if len(userApiKey) == 0 {
		fmt.Println("Warning!: No Password Provided. Used SECRET_PASSWORD instead")
		userApiKey = "SECRET_PASSWORD"
	}

	adminApiKey := os.Getenv("API_KEY_ADMIN")
	if len(adminApiKey) == 0 {
		fmt.Println("Warning!: No Password Provided. Used SECRET_ADMIN_PASSWORD instead")
		adminApiKey = "SECRET_ADMIN_PASSWORD"
	}

	dbUrl := os.Getenv("DB_URL")
	var queries database.IQueries
	if len(dbUrl) != 0 {
		db, errDB := sql.Open("postgres", dbUrl)
		if errDB != nil {
			log.Println(errDB)
		}
		queries = database.New(db)
	} else {
		queries = database.NewMockDB()
	}

	cfg := ApiConfig{
		dbQueries:    queries,
		port:         portFromEnv,
		user_api_key: userApiKey,
	}

	smux := CustomSmux("Api.log")

	smux.HandleFunc(HTTPMethod.GET+" /v1/orders", cfg.GetTasks)

	smux.HandleFunc(HTTPMethod.GET+" /v1/orders/{id}", cfg.GetTasksByID)

	smux.HandleFunc(HTTPMethod.GET+" /v1/orders/pdf/{year}/{month}", cfg.GetTaskByMonth)

	smux.HandleFunc(HTTPMethod.POST+" /v1/orders", Authorized("ApiKey", cfg.user_api_key, cfg.CreateTask))

	smux.HandleFunc(HTTPMethod.POST+" /v1/orders/csv", Authorized("ApiKey", cfg.user_api_key, cfg.UpdloadTasksCSV))

	smux.HandleFunc(HTTPMethod.UPDATE+" /v1/orders/{id}", Authorized("ApiKey", cfg.user_api_key, cfg.UpdateTaskStatus))

	smux.HandleFunc(HTTPMethod.POST+" /v1/reset", Authorized("ApiKey", cfg.user_api_key, cfg.Reset))

	adminCfg := AdminConfig{
		dbQueries:             queries,
		backgroundTaskHandler: background.NewBackgroundWorker(),
		admin_api_key:         adminApiKey,
	}

	smux.HandleFunc(HTTPMethod.POST+" /v1/amdin/sheduler", Authorized("ApiKey", adminCfg.admin_api_key, adminCfg.HandlerHook))

	smux.HandleFunc(HTTPMethod.GET+" /v1/amdin/sheduler", Authorized("ApiKey", adminCfg.admin_api_key, adminCfg.GetHandlerInfo))

	smux.HandleFunc(HTTPMethod.GET+" /v1/amdin/sheduler/{taskid}", Authorized("ApiKey", adminCfg.admin_api_key, adminCfg.GetTaskInfo))

	server := http.Server{Handler: smux, Addr: ":" + portFromEnv}
	fmt.Println("listening on Port: ", server.Addr)
	errServer := server.ListenAndServe()
	if errServer != nil {
		log.Fatal(errServer)
	}

}
