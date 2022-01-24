package internalhttp

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

type RootHandler struct {
	app             Application
	logger          Logger
	calendarHandler *CalendarHandler
}

func NewRootHandler(app Application, logger Logger) *RootHandler {
	h := &RootHandler{
		app:             app,
		logger:          logger,
		calendarHandler: NewCalendarHandler(app, logger),
	}
	return h
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, World")
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h *RootHandler) handleAPI(w http.ResponseWriter, r *http.Request) {
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
