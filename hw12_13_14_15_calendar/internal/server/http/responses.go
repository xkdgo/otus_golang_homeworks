package internalhttp

import (
	"encoding/json"
	"net/http"
)

func httpJSON(w http.ResponseWriter, v interface{}, logger Logger) {
	data, err := json.Marshal(v)
	if err != nil {
		httpInternalServerError(w, "failed to encode response body", err, logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		logger.Error("failed to write response")
	}
}

func httpBadRequest(w http.ResponseWriter, msg string, err error, logger Logger) {
	httpError(w, http.StatusBadRequest, msg, err, logger)
}

func httpInternalServerError(w http.ResponseWriter, msg string, err error, logger Logger) {
	httpError(w, http.StatusInternalServerError, msg, err, logger)
}

func httpError(w http.ResponseWriter, httpStatus int, msg string, err error, logger Logger) {
	http.Error(w, msg, httpStatus)
	logger.Error(msg, err.Error())
}
