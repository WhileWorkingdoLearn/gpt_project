package main

import (
	"fmt"
	"log"
	"os"

	"github.com/WhileCodingDoLearn/gpt_project/repository"
	endpoint "github.com/WhileCodingDoLearn/gpt_project/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func goDotEnvVar(key string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}

	return os.Getenv(key), nil
}

func main() {
	fmt.Println("Start Server")

	port, errPort := goDotEnvVar("SERVER_PORT")
	if errPort != nil {
		port = "8080"
	}

	middlewares := []gin.HandlerFunc{
		endpoint.LoggerMiddleware(endpoint.SetupLogger()),
	}

	dbconn, errEnv := goDotEnvVar("DB_CONNECTUON")
	if errEnv != nil {
		log.Fatalf("Environment variable not found %v \n", errEnv)
	}

	database, errDB := repository.ConnectToDB(dbconn)
	if errDB != nil {
		//log.Fatalf("Connection was not able %v \n", errDB)
		fmt.Printf("Connection was not able %v \n", errDB)
	}

	err := endpoint.StartSever(port, middlewares, database)
	if err != nil {
		log.Fatalf("server could not run")
	}
	/*

	 */

}
