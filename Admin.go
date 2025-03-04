package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/WhileCodingDoLearn/gpt_project/internal/background"
	"github.com/WhileCodingDoLearn/gpt_project/internal/database"
	"github.com/WhileCodingDoLearn/gpt_project/internal/response"
	"github.com/google/uuid"
)

type AdminConfig struct {
	dbQueries             database.IQueries
	backgroundTaskHandler *background.TaskHandler
	admin_api_key         string
	token_secret          string
	plattform             string
}

type AdminRequest struct {
	ID        uuid.UUID     `json:"id"`
	Method    string        `json:"method"`
	Interval  time.Duration `json:"interval"`
	TimeRange time.Duration `json:"timerange"`
}

func (adminCfg AdminConfig) HandlerHook(writer http.ResponseWriter, req *http.Request) {
	var adminReq AdminRequest
	errDec := json.NewDecoder(req.Body).Decode(&adminReq)
	if errDec != nil {
		response.RespondWithError(writer, http.StatusBadRequest, errDec.Error())
		return
	}

	if adminCfg.backgroundTaskHandler == nil {
		response.RespondWithError(writer, http.StatusInternalServerError, "handler not defined")
		return
	}

	if adminReq.Interval <= 0 {
		adminReq.Interval = 1
	}
	adminCfg.backgroundTaskHandler.SetInterval(adminReq.Interval * time.Second)
	if adminReq.TimeRange <= 0 {
		adminReq.TimeRange = 1
	}

	switch adminReq.Method {
	case "run_unlimited":
		{
			err := adminCfg.backgroundTaskHandler.Run()
			if err != nil {
				response.RespondWithError(writer, http.StatusNotAcceptable, errDec.Error())
				return
			}
		}
	case "run_limited":
		{
			err := adminCfg.backgroundTaskHandler.RunUntil(adminReq.TimeRange * time.Second)
			if err != nil {
				response.RespondWithError(writer, http.StatusNotAcceptable, errDec.Error())
				return
			}
		}
	case "update_interval":
		{
			adminCfg.backgroundTaskHandler.SetInterval(adminReq.Interval * time.Second)
		}
	case "reset_counter":
		{
			adminCfg.backgroundTaskHandler.ResetCounter()
		}
	case "stop":
		{
			adminCfg.backgroundTaskHandler.Stop()
		}

	default:
		{
			response.RespondWithError(writer, http.StatusBadRequest, "uknown method")
			return
		}

	}

	var handlerInfo = TaskHandler{
		Interval: adminCfg.backgroundTaskHandler.GetSetInterval() / time.Second,
		Tasks:    adminCfg.backgroundTaskHandler.GetAllTaskInfo(),
		RunUntil: adminCfg.backgroundTaskHandler.GetRunUntil(),
		IsRuning: adminCfg.backgroundTaskHandler.IsRuning(),
		Cycles:   adminCfg.backgroundTaskHandler.PastInterval(),
	}
	response.ResponseJSON(writer, http.StatusOK, handlerInfo)

}

type TaskHandler struct {
	Interval time.Duration         `json:"interval"`
	Tasks    []background.TaskInfo `json:"taks"`
	RunUntil string                `json:"run_until"`
	IsRuning bool                  `json:"is_running"`
	Cycles   int                   `json:"previous runs"`
}

func (adminCfg AdminConfig) GetHandlerInfo(writer http.ResponseWriter, req *http.Request) {
	if adminCfg.backgroundTaskHandler == nil {
		response.RespondWithError(writer, http.StatusInternalServerError, "handler not defined")
		return
	}

	var handlerInfo = TaskHandler{
		Interval: adminCfg.backgroundTaskHandler.GetSetInterval() / time.Second,
		Tasks:    adminCfg.backgroundTaskHandler.GetAllTaskInfo(),
		RunUntil: adminCfg.backgroundTaskHandler.GetRunUntil(),
		IsRuning: adminCfg.backgroundTaskHandler.IsRuning(),
		Cycles:   adminCfg.backgroundTaskHandler.PastInterval(),
	}
	response.ResponseJSON(writer, http.StatusOK, handlerInfo)
}

func (adminCfg AdminConfig) GetTaskInfo(writer http.ResponseWriter, req *http.Request) {
	idFromPath := req.PathValue("taskid")
	id, errParse := uuid.Parse(idFromPath)
	if errParse != nil {
		response.RespondWithError(writer, http.StatusBadRequest, errParse.Error())
		return
	}
	taskInfo, errTaskInfo := adminCfg.backgroundTaskHandler.GetTaskInfo(id.String())
	if errTaskInfo != nil {
		response.RespondWithError(writer, http.StatusNotFound, errTaskInfo.Error())
		return
	}
	response.ResponseJSON(writer, http.StatusOK, taskInfo)
}
