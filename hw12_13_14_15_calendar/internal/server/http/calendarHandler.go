package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/server/http/internal/models"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
)

type ContextKey string

const ContextUserKey ContextKey = "user"

// "api/v1/calendar/"  handler.
type CalendarHandler struct {
	app     Application
	storage storage.Storage
	logger  Logger
}

func NewCalendarHandler(app Application, logger Logger) *CalendarHandler {
	return &CalendarHandler{
		app:     app,
		storage: app.GetStorage(),
		logger:  logger,
	}
}

func (h *CalendarHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	fmt.Println(head == "", r.URL.Path)
	switch {
	case head == "" && r.URL.Path == "/":
		h.HelloCalendar(w, r)
	case head == "event":
		h.handleEvent(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func (h *CalendarHandler) HelloCalendar(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, Calendar")
}

func (h *CalendarHandler) handleEvent(w http.ResponseWriter, r *http.Request) {
	// "event/..." handler.
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	fmt.Println(head == "", r.URL.Path)
	switch {
	case head == "" && r.URL.Path == "/":
		HelloEventHandler(w, r)
	case head == "create" && r.URL.Path == "/":
		if r.Method == "POST" {
			h.handleCreateEvent(w, r)
		} else {
			UnsupportedMethod(w, r)
		}
	}
}

func (h *CalendarHandler) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	userIDValue := ctx.Value(ContextUserKey)
	userID, ok := userIDValue.(string)
	if !ok {
		InvalidUser(w, r)
		return
	}
	jsDecoder := json.NewDecoder(r.Body)
	var event models.Event
	err := jsDecoder.Decode(&event)
	if err != nil {
		httpBadRequest(w, "failed to parse create event", err, h.logger)
		return
	}
	_ = userID
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, your event is %#v", event)
}

func HelloEventHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, This is Event Handler")
}

func UnsupportedMethod(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "This method not allowed")
}

func InvalidUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, "Bad User")
}
