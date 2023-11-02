package restapihandlers

import (
	"encoding/json"
	"net/http"

	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/pkg/logger"
)

type errResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data any, logger logger.AppLogger, headers ...http.Header) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func errorJSON(w http.ResponseWriter, err error, logger logger.AppLogger) {
	var payload errResponse
	payload.Error = true
	payload.Message = err.Error()

	logger.Error(err.Error())
	status := models.ErrCodeToHTTPStatus(err)

	writeJSON(w, status, payload, logger)
}
