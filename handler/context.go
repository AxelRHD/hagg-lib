package handler

import (
	"log/slog"
	"net/http"
	"strings"

	g "maragu.dev/gomponents"

	"github.com/axelrhd/hagg-lib/hxevents"
	"github.com/axelrhd/hagg-lib/toast"
)

// Context is the custom request context that replaces gin.Context.
// It provides a minimal, framework-agnostic abstraction over HTTP request/response.
type Context struct {
	Res http.ResponseWriter // Response writer (explicit field, no embedding)
	Req *http.Request       // Request (explicit field, no embedding)

	logger *slog.Logger // Structured logger
	events []Event      // Event accumulator for frontend communication
}

// Event represents a single event to be sent to the frontend.
// Events can be delivered via HX-Trigger headers (HTMX) or initial-events script (full-page).
type Event struct {
	Name    string `json:"name"`    // Event name (e.g., "toast", "auth-changed")
	Payload any    `json:"payload"` // Event payload (must be JSON-serializable)
}

// Render renders a gomponents node to the HTTP response.
// Sets Content-Type to text/html and writes the node's HTML output.
func (c *Context) Render(node g.Node) error {
	c.Res.Header().Set("Content-Type", "text/html; charset=utf-8")
	return node.Render(c.Res)
}

// Event adds an event to the context's event queue.
// Events are committed to HX-Trigger headers or initial-events script at the end of the request.
func (c *Context) Event(name string, payload any) {
	c.events = append(c.events, Event{
		Name:    name,
		Payload: payload,
	})
}

// Events returns all accumulated events.
// Used internally by hxevents package for committing events.
func (c *Context) Events() []Event {
	return c.events
}

// Logger returns the structured logger for this request.
func (c *Context) Logger() *slog.Logger {
	return c.logger
}

// Toast creates a new toast builder for this request.
// Returns a fluent builder for configuring and emitting toast notifications.
//
// Example:
//
//	ctx.Toast("User created").Success().Notify()
//	ctx.Toast("Error occurred").Error().Stay().Notify()
func (c *Context) Toast(msg string) *toast.Toast {
	return toast.New(msg, c)
}

// NoContent writes a 204 No Content response with HX-Trigger headers.
// This is a helper that commits events before writing the status code.
// Use this instead of manually calling WriteHeader(http.StatusNoContent).
//
// Example:
//
//	ctx.Toast("Operation successful").Success().Notify()
//	return ctx.NoContent()
func (c *Context) NoContent() error {
	c.commitEvents()
	c.Res.WriteHeader(http.StatusNoContent)
	return nil
}

// commitEvents commits accumulated events to HX-Trigger headers.
// This is called internally before writing response headers.
func (c *Context) commitEvents() {
	// Convert handler.Event to hxevents.Event and add default phase prefix
	hxEvents := make([]hxevents.Event, len(c.events))
	for i, e := range c.events {
		// Add "HX-Trigger:" prefix if event doesn't have a phase prefix
		name := e.Name
		if !hasPhasePrefix(name) {
			name = "HX-Trigger:" + name
		}
		hxEvents[i] = hxevents.Event{Name: name, Payload: e.Payload}
	}

	// Commit events (errors are logged but don't fail the request)
	_ = hxevents.Commit(c.Res, c.Req, hxEvents)
}

// hasPhasePrefix checks if an event name has a phase prefix.
func hasPhasePrefix(name string) bool {
	return strings.HasPrefix(name, "HX-Trigger:") ||
		strings.HasPrefix(name, "HX-Trigger-After-Swap:") ||
		strings.HasPrefix(name, "HX-Trigger-After-Settle:")
}
