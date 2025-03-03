package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/WhileCodingDoLearn/gpt_project/internal/background"
	"github.com/WhileCodingDoLearn/gpt_project/internal/database"
	"github.com/WhileCodingDoLearn/gpt_project/internal/response"
)

type AdminConfig struct {
	dbQueries             database.IQueries
	backgroundTaskHandler *background.TaskHandler
	admin_api_key         string
	token_secret          string
	plattform             string
}

type AdminRequest struct {
	Method   string `json:"method`
	Interval int    `json:timerange`
}

func (adminCfg AdminConfig) StartHandler(writer http.ResponseWriter, req *http.Request) {
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

	switch adminReq.Method {
	case "run_unlimited":
		{
			err := adminCfg.backgroundTaskHandler.Run()
			if err != nil {
				response.RespondWithError(writer, http.StatusNotAcceptable, errDec.Error())
				return
			}
			response.ResponseWithoutBody(writer, http.StatusOK)
		}
	case "run_limited":
		{
			err := adminCfg.backgroundTaskHandler.RunUntil(time.Duration(adminReq.Interval) * time.Second)
			if err != nil {
				response.RespondWithError(writer, http.StatusNotAcceptable, errDec.Error())
				return
			}
			response.ResponseWithoutBody(writer, http.StatusOK)
		}
	case "stop":
		{
			adminCfg.backgroundTaskHandler.Stop()
			response.ResponseWithoutBody(writer, http.StatusOK)
			return
		}

	default:
		{
			response.RespondWithError(writer, http.StatusBadRequest, "uknown method")
			return
		}

	}

}
