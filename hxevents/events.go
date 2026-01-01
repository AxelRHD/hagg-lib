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
