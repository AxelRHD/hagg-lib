package hxevents

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// RenderInitialEvents creates a <script> tag with initial events for full-page loads.
// This is the counterpart to Commit() for non-HTMX requests.
//
// Returns an empty text node if:
//   - This is an HTMX request (use Commit() instead)
//   - There are no non-phased events
//
// Only renders events added with ctx.Event() (no phase prefix).
// Events added with hxevents.Add() are filtered out (they're HTMX-only).
//
// The frontend processes this script on DOMContentLoaded and triggers the same
// event handlers as HX-Trigger events, creating a unified event system.
//
// Example output:
//
//	<script type="application/json" id="initial-events">
//	[
//	  {"name":"toast","payload":{"message":"Welcome!","level":"info"}},
//	  {"name":"auth-changed","payload":null}
//	]
//	</script>
func RenderInitialEvents(req *http.Request, events []Event) g.Node {
	// HTMX requests use headers, not initial-events script
	if IsHtmxRequest(req.Header) {
		return g.Text("")
	}

	// Filter out phase-prefixed events (those are for HTMX only)
	var initialEvents []Event
	for _, evt := range events {
		// Skip events that start with phase prefixes
		isPhaseEvent := strings.HasPrefix(evt.Name, string(Immediate)+":") ||
			strings.HasPrefix(evt.Name, string(AfterSwap)+":") ||
			strings.HasPrefix(evt.Name, string(AfterSettle)+":")

		if !isPhaseEvent {
			initialEvents = append(initialEvents, evt)
		}
	}

	// No events to render
	if len(initialEvents) == 0 {
		return g.Text("")
	}

	// Marshal events to JSON
	jsonData, err := json.Marshal(initialEvents)
	if err != nil {
		// Log error but don't fail rendering
		return g.Text("")
	}

	return h.Script(
		h.Type("application/json"),
		h.ID("initial-events"),
		g.Raw(string(jsonData)),
	)
}

// RenderToasts creates self-destructing toast elements for full-page loads.
// Uses Surreal.js pattern: Script tag with showToast() call + me().remove().
//
// Returns nil if:
//   - This is an HTMX request (toasts are sent via HX-Trigger header)
//   - There are no toast events
//
// Example output:
//
//	<div>
//	  <script>showToast({"message":"Welcome!","level":"success",...})</script>
//	  <script>me().remove()</script>
//	</div>
func RenderToasts(req *http.Request, events []Event) g.Node {
	// HTMX requests use headers, not DOM elements
	if IsHtmxRequest(req.Header) {
		return nil
	}

	// Filter out phase-prefixed events
	var toasts []g.Node
	for _, evt := range events {
		// Skip phase-prefixed events
		isPhaseEvent := strings.HasPrefix(evt.Name, string(Immediate)+":") ||
			strings.HasPrefix(evt.Name, string(AfterSwap)+":") ||
			strings.HasPrefix(evt.Name, string(AfterSettle)+":")

		if isPhaseEvent {
			continue
		}

		// Only render toast events
		if evt.Name != "toast" {
			continue
		}

		// Marshal payload to JSON
		toastJSON, err := json.Marshal(evt.Payload)
		if err != nil {
			continue
		}

		// Create self-destructing toast element
		toasts = append(toasts, h.Div(
			h.Script(g.Raw(fmt.Sprintf(`showToast(%s)`, toastJSON))),
			h.Script(g.Raw(`me().remove()`)),
		))
	}

	// Return nil if no toasts (g.Group with nil elements causes issues)
	if len(toasts) == 0 {
		return nil
	}

	return g.Group(toasts)
}
