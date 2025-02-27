package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/WhileCodingDoLearn/gpt_project/internal/database"
	"github.com/WhileCodingDoLearn/gpt_project/internal/response"
	"github.com/WhileCodingDoLearn/gpt_project/pdf"
	"github.com/google/uuid"
)

var HTTPMethod = struct {
	POST   string
	GET    string
	DELETE string
	UPDATE string
}{POST: "POST", GET: "GET", DELETE: "DELETE", UPDATE: "PATCH"}

var PathValues = struct {
	TaksId string
}{
	TaksId: "id",
}

type ApiConfig struct {
	dbQueries     *database.Queries
	plattform     string
	user_api_key  string
	admin_api_key string
	token_secret  string
	port          string
}

type CreateTaskParams struct {
	OrderID uuid.UUID `json:"order_id"`
	Name    string    `json:"name"`
	Data    string    `json:"data"`
	Status  string    `json:"status"`
}

type Task struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	OrderID   uuid.UUID `json:"order_id"`
	Name      string    `json:"name"`
	Data      string    `json:"data"`
	Status    string    `json:"status"`
}

func (apiCfg *ApiConfig) CreateTask(w http.ResponseWriter, req *http.Request) {
	var task CreateTaskParams
	errDecode := json.NewDecoder(req.Body).Decode(&task)
	if errDecode != nil {
		response.RespondWithError(w, http.StatusBadRequest, errDecode.Error())
		return
	}

	taskFromDB, errTaskDB := apiCfg.dbQueries.CreateTask(req.Context(), database.CreateTaskParams{
		OrderID: task.OrderID,
		Name:    task.Name,
		Data:    task.Data,
		Status:  task.Status,
	})

	if errTaskDB != nil {
		response.RespondWithError(w, http.StatusNoContent, errTaskDB.Error())
		return
	}
	response.ResponseJSON(w, http.StatusCreated, Task{
		ID:        taskFromDB.ID,
		CreatedAt: taskFromDB.CreatedAt,
		UpdatedAt: taskFromDB.UpdatedAt,
		OrderID:   taskFromDB.OrderID,
		Name:      taskFromDB.Name,
		Data:      taskFromDB.Data,
		Status:    taskFromDB.Status,
	})

}

func (apiCfg *ApiConfig) GetTasks(w http.ResponseWriter, req *http.Request) {
	data, errDB := apiCfg.dbQueries.GetAllTasks(req.Context(), database.GetAllTasksParams{
		Limit:  100,
		Offset: 0,
	})
	if errDB != nil {
		response.RespondWithError(w, http.StatusInternalServerError, errDB.Error())
		return
	}
	response.ResponseJSON(w, http.StatusOK, data)
}

func (apiCfg *ApiConfig) GetTasksByID(w http.ResponseWriter, req *http.Request) {
	taskId := req.PathValue(PathValues.TaksId)
	taskUuid, errParse := uuid.Parse(taskId)
	if errParse != nil {
		response.RespondWithError(w, http.StatusBadRequest, errParse.Error())
		return
	}

	data, errDB := apiCfg.dbQueries.GetTaskByID(req.Context(), taskUuid)
	if errDB != nil {
		response.RespondWithError(w, http.StatusInternalServerError, errDB.Error())
		return
	}
	response.ResponseJSON(w, http.StatusOK, data)
}

func (apiCfg *ApiConfig) GetTaskByMonth(w http.ResponseWriter, req *http.Request) {
	yearParam := req.PathValue("year")
	year, errYear := strconv.Atoi(yearParam)
	if errYear != nil || year < 2000 || year > 2200 {
		response.RespondWithError(w, http.StatusBadRequest, errYear.Error())
		return
	}

	monthParam := req.PathValue("month")
	month, errMonth := strconv.Atoi(monthParam)
	if errMonth != nil || month < 1 || month > 12 {
		response.RespondWithError(w, http.StatusBadRequest, errMonth.Error())
		return
	}

	timeRange := time.Date(year, time.Month(month+1), 0,
		0, 0, 0, 0, time.UTC)

	orders, errDB := apiCfg.dbQueries.GetTaskByMonth(req.Context(), timeRange)
	if errDB != nil {
		response.RespondWithError(w, http.StatusInternalServerError, errDB.Error())
		return
	}

	data := make([]pdf.Task, 0)
	for _, o := range orders {
		data = append(data, pdf.Task{
			ID:        o.ID,
			CreatedAt: o.CreatedAt,
			UpdatedAt: o.UpdatedAt,
			OrderID:   o.OrderID,
			Name:      o.Name,
			Data:      o.Data,
			Status:    o.Status,
		})
	}

	datapdf, errPDF := pdf.GenerateDocument(data)
	if errPDF != nil {
		response.RespondWithError(w, http.StatusInternalServerError, errPDF.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(datapdf)

}

type UpdateTaskNameRequest struct {
	Id      uuid.UUID `json:"id"`
	NewName string    `json:"new_name"`
}

func (apiCfg *ApiConfig) UpdateTasksNameByID(w http.ResponseWriter, req *http.Request) {
	var newName UpdateTaskNameRequest
	errDecode := json.NewDecoder(req.Body).Decode(&newName)
	if errDecode != nil {
		response.RespondWithError(w, http.StatusBadRequest, errDecode.Error())
		return
	}
	if len(newName.NewName) < 3 {
		response.ResponseWithoutBody(w, http.StatusBadRequest)
		return
	}

	errDB := apiCfg.dbQueries.UpdateTaskName(req.Context(), database.UpdateTaskNameParams{
		ID:   newName.Id,
		Name: newName.NewName,
	})
	if errDB != nil {
		response.RespondWithError(w, http.StatusNotFound, errDB.Error())
		return
	}
	response.ResponseWithoutBody(w, http.StatusAccepted)
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status"`
}

func (apiCfg *ApiConfig) UpdateTaskStatus(w http.ResponseWriter, req *http.Request) {
	orderId := req.PathValue("id")
	idFromOrder, errUuid := uuid.Parse(orderId)
	if errUuid != nil {
		response.RespondWithError(w, http.StatusBadRequest, "order could not be parsed to uuid: "+orderId)
		return
	}

	var reqBody UpdateTaskStatusRequest

	errDecode := json.NewDecoder(req.Body).Decode(&reqBody)
	if errDecode != nil {
		response.RespondWithError(w, http.StatusBadRequest, errDecode.Error())
		return
	}

	newStatus := strings.Trim(strings.ToLower(reqBody.Status), " ")

	errUpdate := apiCfg.dbQueries.UpdateTaskStatus(req.Context(), database.UpdateTaskStatusParams{
		ID:     idFromOrder,
		Status: newStatus,
	})
	if errUpdate != nil {
		response.RespondWithError(w, http.StatusBadRequest, errUpdate.Error())
		return
	}

	response.ResponseWithoutBody(w, http.StatusOK)

}

type CVUploadResponse struct {
	created []Task
	failed  []Task
}

func (apiCfg *ApiConfig) UpdloadTasksCSV(w http.ResponseWriter, req *http.Request) {
	var resp CVUploadResponse
	req.ParseMultipartForm(10 << 20)

	// Datei aus Request abrufen
	file, _, err := req.FormFile("file")
	if err != nil {
		response.ResponseWithoutBody(w, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// CSV-Datei parsen
	reader := csv.NewReader(bufio.NewReader(file))
	records, err := reader.ReadAll()
	if err != nil {
		response.ResponseWithoutBody(w, http.StatusInternalServerError)
		return
	}

	// Header überprüfen
	if len(records) < 1 || records[0][0] != "order_id" {
		response.ResponseWithoutBody(w, http.StatusBadRequest)
		return
	}

	// Daten verarbeiten und einfügen
	failed := make([]Task, 0)
	created := make([]Task, 0)
	for _, record := range records[1:] {
		orderID, err := uuid.Parse(record[0])
		if err != nil {
			failed = append(failed, Task{
				ID: orderID,
			})
			continue
		}
		createdTask, errDB := apiCfg.dbQueries.CreateTask(req.Context(), database.CreateTaskParams{
			OrderID: orderID,
			Name:    record[1],
			Data:    record[2],
			Status:  record[3],
		})
		if errDB != nil {
			failed = append(failed, Task{
				OrderID: orderID,
				Name:    record[1],
				Data:    record[2],
				Status:  record[3],
			})
			continue
		}
		created = append(created, Task{
			ID:        createdTask.ID,
			CreatedAt: createdTask.CreatedAt,
			UpdatedAt: createdTask.UpdatedAt,
			OrderID:   createdTask.OrderID,
			Name:      createdTask.Name,
			Data:      createdTask.Data,
			Status:    createdTask.Status,
		})

	}
	if len(created) == 0 {
		response.ResponseWithoutBody(w, http.StatusNoContent)
		return
	}
	resp.created = created
	resp.failed = failed
	response.ResponseJSON(w, http.StatusCreated, resp)

}

func (apiCfg *ApiConfig) Reset(w http.ResponseWriter, req *http.Request) {
	err := apiCfg.dbQueries.DeleteAllTasks(req.Context())
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	response.ResponseWithoutBody(w, http.StatusOK)
}
