package handler

import (
	"log/slog"
	"net/http"
)

// HandlerFunc is the custom handler signature that works with our Context.
// Returns an error for centralized error handling.
type HandlerFunc func(*Context) error

// Wrapper wraps custom HandlerFunc to work with stdlib http.HandlerFunc.
// Provides:
//   - Context creation and injection
//   - Centralized error handling
//   - Automatic event commitment (via hxevents)
type Wrapper struct {
	logger *slog.Logger
}

// NewWrapper creates a new handler wrapper with the given logger.
func NewWrapper(logger *slog.Logger) *Wrapper {
	return &Wrapper{logger: logger}
}

// Logger returns the logger instance used by this wrapper.
// This is useful for middleware that needs access to the logger.
func (w *Wrapper) Logger() *slog.Logger {
	return w.logger
}

// Wrap converts a custom HandlerFunc to stdlib http.HandlerFunc.
//
// Flow:
//  1. Create Context with response writer, request, and logger
//  2. Call the handler
//  3. If handler returns error: log it and return 500
//  4. If handler succeeds: commit events to headers/script
//
// Example usage:
//
//	wrapper := handler.NewWrapper(slog.Default())
//	http.HandleFunc("/", wrapper.Wrap(myHandler))
func (w *Wrapper) Wrap(h HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Create context for this request
		ctx := &Context{
			Res:    res,
			Req:    req,
			logger: w.logger,
			events: make([]Event, 0),
		}

		// Call the handler
		if err := h(ctx); err != nil {
			w.logger.Error("handler error",
				"path", req.URL.Path,
				"method", req.Method,
				"error", err,
			)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Commit events to HX-Trigger headers (for HTMX requests)
		// For full-page loads, events are rendered via hxevents.RenderInitialEvents() in layout
		// Note: Event commitment is also done via ctx.NoContent() for handlers that need it
		// This is a fallback for handlers that don't explicitly commit events
		ctx.commitEvents()
	}
}
