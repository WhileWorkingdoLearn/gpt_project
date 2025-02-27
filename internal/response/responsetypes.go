package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResponseWithoutBody(writer http.ResponseWriter, statuscode int) {
	writer.WriteHeader(statuscode)
}

func ResponseJSON(writer http.ResponseWriter, statuscode int, payload interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	msg, errJSON := json.Marshal(payload)
	if errJSON != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error encoding response: %s", errJSON)
		return
	}
	writer.WriteHeader(statuscode)
	writer.Write(msg)
}

type errorResp struct {
	Err string `json:"err"`
}

func RespondWithError(writer http.ResponseWriter, statusCode int, msg string) {
	errorStruct := errorResp{
		Err: msg,
	}
	writer.Header().Set("Content-Type", "application/json")

	data, errJSON := json.Marshal(errorStruct)

	if errJSON != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error encoding response: %s", errJSON)
		return
	}
	writer.WriteHeader(statusCode)
	writer.Write(data)
}
