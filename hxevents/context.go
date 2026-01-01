package hxevents

// EventAdder is the interface for adding events (avoids import cycle with handler package).
type EventAdder interface {
	Event(name string, payload any)
}

// Add adds an event for a specific HTMX phase.
// Events added with a phase are only sent via HX-Trigger headers (not initial-events).
//
// The event is stored with a phase prefix (e.g., "HX-Trigger:eventname") to distinguish
// it from non-phased events that should appear in initial-events scripts.
//
// Example:
//
//	hxevents.Add(ctx, hxevents.Immediate, "auth-changed", map[string]any{"user": "alice"})
func Add(ctx EventAdder, phase Phase, name string, payload any) {
	// Convention: Store phase-events with "phase:name" format
	// This allows us to filter them out from initial-events rendering
	eventName := string(phase) + ":" + name
	ctx.Event(eventName, payload)
}
