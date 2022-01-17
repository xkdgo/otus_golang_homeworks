package internalhttp

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				s.logger.Info(
					time.Now().Format("2006-01-02T15:04:05.999999999Z07:00"),
					"err", err,
					"trace", string(debug.Stack()),
				)
			}
		}()
		reqpath := r.URL.EscapedPath()
		start := time.Now()
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
		s.logger.Info(
			r.RemoteAddr,
			start.Format("[02/Jan/2006:15:04:05 -0700]"),
			r.Method,
			reqpath,
			r.Proto,
			wrapped.Status(),
			time.Since(start),
			r.UserAgent(),
		)
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDKey := r.Header.Get("x-user-id")
		switch {
		case IsValidUUID(userIDKey):
			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextUserKey, userIDKey)
			r = r.WithContext(ctx)
		default:
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
