package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/WhileCodingDoLearn/gpt_project/docs"
	"github.com/WhileCodingDoLearn/gpt_project/internal/background"
	"github.com/WhileCodingDoLearn/gpt_project/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

var VERSION string

func init() {
	VERSION = "/V1"
}

// @title           Meine API
// @version         1.0
// @description     Dies ist eine Beispiel-API mit Swagger-Dokumentation
// @host           localhost:8080
// @BasePath       /v1

// @Summary        Holt einen Benutzer
// @Description    Holt einen Benutzer anhand der ID
// @Tags           Benutzer
// @Accept         json
// @Produce        json
// @Param          id  path  int  true  "Benutzer ID"
// @Success        200  {object}  map[string]string
// @Failure        400  {object}  map[string]string
// @Router         /users/{id} [get]

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
	/**/
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

	smux := CustomSmux("Api.log")

	cfg := ApiConfig{
		dbQueries:    queries,
		port:         portFromEnv,
		user_api_key: userApiKey,
	}

	smux.HandleFunc(HTTPMethod.GET+" /swagger/", httpSwagger.WrapHandler)

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

	server := &http.Server{Handler: smux, Addr: ":" + portFromEnv}
	fmt.Println("listening on Port: ", server.Addr)
	errServer := server.ListenAndServe()

	if errServer != nil {
		log.Fatal(errServer)
	}

}
