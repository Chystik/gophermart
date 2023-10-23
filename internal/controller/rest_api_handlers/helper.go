package restapihandlers

import (
	"encoding/json"
	"net/http"

	"github.com/Chystik/gophermart/pkg/logger"
)

type errResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, contentType string, data any, logger logger.AppLogger, headers ...http.Header) {
	out, err := json.Marshal(data)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func errorJSON(w http.ResponseWriter, err error, status int, logger logger.AppLogger) {
	var payload errResponse
	payload.Error = true
	payload.Message = err.Error()

	logger.Error(err.Error())
	writeJSON(w, status, "application/json", payload, logger)
}
