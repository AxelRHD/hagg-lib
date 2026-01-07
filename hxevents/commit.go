package hxevents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Event represents a single event (matches handler.Event to avoid import cycle).
type Event struct {
	Name    string `json:"name"`
	Payload any    `json:"payload"`
}

// Commit writes accumulated events to HX-Trigger response headers.
// This function should be called after the handler completes, before the response is sent.
//
// Only works for HTMX requests (checks HX-Request header).
// For full-page loads, use hxevents.RenderInitialEvents() instead.
//
// Events are grouped by phase:
//   - Events added with hxevents.Add() are sent via the appropriate HX-Trigger header
//   - Events added with ctx.Event() (no phase) are ignored by this function
//
// Example output headers:
//
//	HX-Trigger: {"toast":{"message":"Success!","level":"success"},"auth-changed":true}
//	HX-Trigger-After-Swap: {"refresh-stats":{"count":42}}
func Commit(res http.ResponseWriter, req *http.Request, events []Event) error {
	// Only commit for HTMX requests
	if !IsHtmxRequest(req.Header) {
		return nil
	}

	// Group events by phase
	phases := map[Phase]map[string]any{
		Immediate:   make(map[string]any),
		AfterSwap:   make(map[string]any),
		AfterSettle: make(map[string]any),
	}

	// Parse events and group by phase
	for _, evt := range events {
		// Check each phase to see if event name has that phase prefix
		for _, phase := range []Phase{Immediate, AfterSwap, AfterSettle} {
			prefix := string(phase) + ":"
			if strings.HasPrefix(evt.Name, prefix) {
				// Remove phase prefix from event name
				name := strings.TrimPrefix(evt.Name, prefix)
				phases[phase][name] = evt.Payload
				break // Event matched a phase, don't check others
			}
		}
	}

	// Write headers for each phase that has events
	for phase, events := range phases {
		if len(events) == 0 {
			continue // Skip phases with no events
		}

		jsonData, err := marshalASCIISafe(events)
		if err != nil {
			return fmt.Errorf("marshal events for %s: %w", phase, err)
		}

		res.Header().Set(string(phase), string(jsonData))
	}

	return nil
}

// marshalASCIISafe marshals data to JSON with non-ASCII characters escaped as \uXXXX.
// This is required for HTTP headers which should only contain ASCII characters.
func marshalASCIISafe(v any) ([]byte, error) {
	// First, marshal normally
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// Escape non-ASCII characters by iterating over runes
	var buf bytes.Buffer
	str := string(data)
	for _, r := range str {
		if r > 127 {
			// Non-ASCII rune - escape as \uXXXX
			if r <= 0xFFFF {
				fmt.Fprintf(&buf, "\\u%04x", r)
			} else {
				// Supplementary character - use surrogate pair
				r -= 0x10000
				high := 0xD800 + (r >> 10)
				low := 0xDC00 + (r & 0x3FF)
				fmt.Fprintf(&buf, "\\u%04x\\u%04x", high, low)
			}
		} else {
			buf.WriteRune(r)
		}
	}

	return buf.Bytes(), nil
}
