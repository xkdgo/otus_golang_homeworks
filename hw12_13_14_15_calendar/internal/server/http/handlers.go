package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type AppHandler struct {
	app    Application
	logger Logger
	mux    *http.ServeMux
}

func NewHandler(app Application, logger Logger, mux *http.ServeMux) *AppHandler {
	h := &AppHandler{app: app, logger: logger, mux: mux}
	mux.Handle("/", http.HandlerFunc(HelloServer))
	mux.Handle("/calendar/event", http.HandlerFunc(h.CRUDEvent))
	return h
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, World")
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *AppHandler) CRUDEvent(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "This is CRUD API")
	case r.Method == "POST":
		// TODO CreateEvent via AppHandler app.
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		status := http.StatusNotFound
		w.WriteHeader(status)
		ans := TestResponse{
			Status:  status,
			Message: "method not implemented yet",
		}
		response, err := json.Marshal(ans)
		if err != nil {
			h.logger.Error(time.Now().Format("02/Jan/2006:15:04:05 -0700"),
				errors.Wrapf(err, ":something wrong when encode json answer"))
			fmt.Fprintf(w, "something wrong when encode json answer")
		}
		_, err = w.Write(response)
		if err != nil {
			h.logger.Error(time.Now().Format("02/Jan/2006:15:04:05 -0700"),
				errors.Wrapf(err, ":something wrong when when write response"))
		}
	case r.Method == "PATCH":
		// TODO UpdateEvent via AppHandler app.
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(time.Now().Format("02/Jan/2006:15:04:05 -0700"), "unimplemented method", r.Method)
		fmt.Fprintf(w, "method %s not implemented", r.Method)
	case r.Method == "DELETE":
		// TODO DeleteEvent via AppHandler app.
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(time.Now().Format("02/Jan/2006:15:04:05 -0700"), "unimplemented method", r.Method)
		fmt.Fprintf(w, "method %s not implemented", r.Method)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(time.Now().Format("02/Jan/2006:15:04:05 -0700"), "unimplemented method", r.Method)
		fmt.Fprintf(w, "method %s not implemented", r.Method)
	}
}

type TestResponse struct {
	Status  int    `json:"status"`
	Message string `json:"messsage"`
}
