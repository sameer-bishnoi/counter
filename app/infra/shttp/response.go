package shttp

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code"`
}

func ResponseJSON(w http.ResponseWriter, httpStatus int, response interface{}) {
	w.WriteHeader(httpStatus)
	if response == nil {
		return
	}

	res := Response{
		Data: response,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("unable to send the response: %v", err)
	}
}

func FailedResponseJSON(w http.ResponseWriter, httpStatus int, errorCode, errorMessage string) {
	res := ErrorResponse{
		ErrorMessage: errorMessage,
		ErrorCode:    errorCode,
	}
	ResponseJSON(w, httpStatus, res)
}
