package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type AppHandler struct {
	app             Application
	logger          Logger
	calendarHandler *CalendarHandler
	// mux    *http.ServeMux
}

type CalendarHandler struct {
	app Application
}

func NewCalendarHandler(app Application) *CalendarHandler {
	return &CalendarHandler{
		app: app,
	}
}

func (h *CalendarHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	fmt.Println(head == "", r.URL.Path)
	switch {
	case head == "" && r.URL.Path == "/":
		h.HelloCalendar(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func (h *CalendarHandler) HelloCalendar(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, Calendar")
}

func NewRootHandler(app Application, logger Logger, mux *http.ServeMux) *AppHandler {
	h := &AppHandler{
		app:             app,
		logger:          logger,
		calendarHandler: NewCalendarHandler(app),
	}
	return h
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, World")
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	fmt.Println(head == "", r.URL.Path)
	switch {
	case head == "" && r.URL.Path == "/":
		HelloServer(w, r)
	case head == "api":
		h.handleAPI(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func (h *AppHandler) handleAPI(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	fmt.Println(head == "", r.URL.Path)
	if head != "v1" {
		http.NotFound(w, r)
		return
	}
	head, r.URL.Path = ShiftPath(r.URL.Path)
	switch head {
	case "calendar":
		h.calendarHandler.ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
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

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
