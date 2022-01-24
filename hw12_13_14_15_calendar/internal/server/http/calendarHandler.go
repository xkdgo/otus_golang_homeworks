package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/server/http/internal/models"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
)

const (
	timelayoutWithMin = "02 Jan 06 15:04 -0700"
	timelayout        = "2006-01-02"
)

type ContextKey string

const ContextUserKey ContextKey = "user"

// "api/v1/calendar/"  handler.
type CalendarHandler struct {
	app    Application
	logger Logger
}

func NewCalendarHandler(app Application, logger Logger) *CalendarHandler {
	return &CalendarHandler{
		app:    app,
		logger: logger,
	}
}

func (h *CalendarHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)

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

	switch {
	case head == "" && r.URL.Path == "/":
		HelloEventHandler(w, r)
	case head == "create" && r.URL.Path == "/":
		if r.Method == http.MethodPost {
			h.handleCreateEvent(w, r)
		} else {
			UnsupportedMethod(w, r)
		}
	case head == "delete":
		if r.Method == http.MethodPost {
			h.handleDeleteEvent(w, r)
		} else {
			UnsupportedMethod(w, r)
		}
	case head == "update" && r.URL.Path == "/":
		if r.Method == http.MethodPost {
			h.handleUpdateEvent(w, r)
		} else {
			UnsupportedMethod(w, r)
		}
	case head == "day" || head == "week" || head == "month":
		if r.Method == http.MethodGet {
			h.handleGetEventsDayWeekMonth(w, r, head)
		} else {
			UnsupportedMethod(w, r)
		}
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
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

	dateTimeStart, err := time.Parse(timelayoutWithMin, event.DateTimeStart)
	if err != nil {
		httpBadRequest(w, "failed to parse datetimestart", err, h.logger)
		return
	}
	duration := event.Duration.Duration
	alarmTime, err := time.Parse(timelayoutWithMin, event.AlarmTime)
	if err != nil {
		httpBadRequest(w, "failed to parse alarmtime", err, h.logger)
		return
	}
	id, err := h.app.CreateEvent(ctx, event.ID, event.Title, userID, dateTimeStart, duration, alarmTime)
	if err != nil {
		httpBadRequest(w, "failed to create event", err, h.logger)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, your event is created with id %v", id)
}

func (h *CalendarHandler) handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
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
		httpBadRequest(w, "failed to parse update event", err, h.logger)
		return
	}

	dateTimeStart, err := time.Parse(timelayoutWithMin, event.DateTimeStart)
	if err != nil {
		httpBadRequest(w, "failed to parse datetimestart", err, h.logger)
		return
	}
	duration := event.Duration.Duration
	alarmTime, err := time.Parse(timelayoutWithMin, event.AlarmTime)
	if err != nil {
		httpBadRequest(w, "failed to parse alarmTime", err, h.logger)
		return
	}
	err = h.app.UpdateEvent(ctx, event.ID, event.Title, userID, dateTimeStart, duration, alarmTime)
	if err != nil {
		httpBadRequest(w, "failed to update event", err, h.logger)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, your event  with id %v is updated", event.ID)
}

func (h *CalendarHandler) handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	var eventID string
	eventID, _ = ShiftPath(r.URL.Path)
	if !IsValidUUID(eventID) {
		httpBadRequest(w, "failed to delete event not valid uuid", fmt.Errorf("%s is not uuid", eventID), h.logger)
		return
	}
	ctx := r.Context()
	err := h.app.DeleteEvent(ctx, eventID)
	if err != nil {
		httpInternalServerError(w, "failed to create event", err, h.logger)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, your event is deleted with id %v", eventID)
}

func (h *CalendarHandler) handleGetEventsDayWeekMonth(w http.ResponseWriter, r *http.Request, period string) {
	day, _ := ShiftPath(r.URL.Path)
	dayTime, err := time.Parse(timelayout, day)
	if err != nil {
		httpBadRequest(w, "failed to parse dayTime", err, h.logger)
		return
	}
	ctx := r.Context()
	userIDValue := ctx.Value(ContextUserKey)
	userID, ok := userIDValue.(string)
	if !ok {
		InvalidUser(w, r)
		return
	}
	var events []storage.Event
	switch period {
	case "day":
		events, err = h.app.ListEventsDay(ctx, userID, dayTime)
		if err != nil {
			httpInternalServerError(w, "failed to list events day", err, h.logger)
			return
		}
	case "week":
		events, err = h.app.ListEventsWeek(ctx, userID, dayTime)
		if err != nil {
			httpInternalServerError(w, "failed to list events week", err, h.logger)
			return
		}
	case "month":
		events, err = h.app.ListEventsMonth(ctx, userID, dayTime)
		if err != nil {
			httpInternalServerError(w, "failed to list events month", err, h.logger)
			return
		}
	}
	modelEvents := convertToModelsEvents(events)
	w.WriteHeader(http.StatusOK)
	httpJSON(w, modelEvents, h.logger)
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

func convertToModelsEvents(events []storage.Event) (modelsEvents []models.Event) {
	modelsEvents = make([]models.Event, 0, len(events))
	for _, event := range events {
		modelEvent := models.Event{}
		modelEvent.ID = event.ID
		modelEvent.Title = event.Title
		modelEvent.DateTimeStart = event.DateTimeStart.Format(timelayoutWithMin)
		modelEvent.Duration.Duration = event.Duration
		modelEvent.AlarmTime = event.AlarmTime.Format(timelayoutWithMin)
		modelsEvents = append(modelsEvents, modelEvent)
	}
	return modelsEvents
}
