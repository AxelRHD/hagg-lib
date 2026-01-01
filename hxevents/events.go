// Package hxevents provides a server-driven event system for HTMX applications.
//
// This package enables servers to emit events to the frontend via:
//   - HX-Trigger headers (for HTMX requests)
//   - Initial-events scripts (for full-page loads)
//
// # Event Phases
//
// HTMX supports three event timing phases:
//
//   - Immediate: Fires as soon as response is received (HX-Trigger)
//   - AfterSwap: Fires after content is swapped into DOM (HX-Trigger-After-Swap)
//   - AfterSettle: Fires after settling completes (HX-Trigger-After-Settle)
//
// # Usage with handler.Context
//
// The most common way to use events is through handler.Context:
//
//	func MyHandler(ctx *handler.Context) error {
//	    ctx.Event("HX-Trigger:user-updated", userData)
//	    ctx.Toast("User updated").Success().Notify()
//	    return ctx.Render(page())
//	}
//
// Events are automatically committed via wrapper.Wrap().
//
// # Direct Usage
//
// For framework-independent usage, use hxevents.Commit directly:
//
//	events := []hxevents.Event{
//	    {Name: "HX-Trigger:refresh", Payload: nil},
//	}
//	hxevents.Commit(w, r, events)
//
// # Event Format
//
// Events use the format "Phase:EventName" where Phase is one of:
//   - HX-Trigger
//   - HX-Trigger-After-Swap
//   - HX-Trigger-After-Settle
//
// Events without a phase prefix are ignored.
//
// # Dependencies
//
// Requires: stdlib (net/http, encoding/json), gomponents
package hxevents

import "net/http"

// Phase represents when an HTMX event should be triggered.
// HTMX supports three trigger timings, each with a corresponding response header.
type Phase string

const (
	// Immediate events fire as soon as the response is received.
	// Header: HX-Trigger
	Immediate Phase = "HX-Trigger"

	// AfterSwap events fire after HTMX swaps the content into the DOM.
	// Header: HX-Trigger-After-Swap
	AfterSwap Phase = "HX-Trigger-After-Swap"

	// AfterSettle events fire after the settling process completes.
	// Header: HX-Trigger-After-Settle
	AfterSettle Phase = "HX-Trigger-After-Settle"
)

// IsHtmxRequest checks if the request is an HTMX request.
// HTMX adds the "HX-Request: true" header to all requests it makes.
func IsHtmxRequest(headers http.Header) bool {
	return headers.Get("HX-Request") == "true"
}
