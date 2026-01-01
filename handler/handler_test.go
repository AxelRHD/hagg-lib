package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// TestContext_Event tests event accumulation
func TestContext_Event(t *testing.T) {
	ctx := &Context{
		events: make([]Event, 0),
	}

	// Add events
	ctx.Event("test-event", map[string]string{"key": "value"})
	ctx.Event("another-event", "simple-payload")

	// Check events were added
	events := ctx.Events()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	// Check first event
	if events[0].Name != "test-event" {
		t.Errorf("expected event name 'test-event', got '%s'", events[0].Name)
	}

	// Check second event
	if events[1].Name != "another-event" {
		t.Errorf("expected event name 'another-event', got '%s'", events[1].Name)
	}
}

// TestContext_Render tests HTML rendering
func TestContext_Render(t *testing.T) {
	// Create response recorder
	rec := httptest.NewRecorder()

	// Create context
	ctx := &Context{
		Res: rec,
	}

	// Create a simple gomponents node
	node := html.Div(nil, g.Text("Hello, World!"))

	// Render
	err := ctx.Render(node)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	// Check Content-Type header
	contentType := rec.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type 'text/html; charset=utf-8', got '%s'", contentType)
	}

	// Check response body
	body := rec.Body.String()
	if !strings.Contains(body, "Hello, World!") {
		t.Errorf("expected body to contain 'Hello, World!', got '%s'", body)
	}
}

// TestContext_Toast tests toast builder creation
func TestContext_Toast(t *testing.T) {
	ctx := &Context{
		events: make([]Event, 0),
	}

	// Create toast
	toastBuilder := ctx.Toast("Test message")

	// Toast builder should not be nil
	if toastBuilder == nil {
		t.Fatal("Toast() returned nil")
	}

	// Call Notify to emit event
	toastBuilder.Success().Notify()

	// Check event was added
	events := ctx.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event after toast, got %d", len(events))
	}

	// Check event name
	if events[0].Name != "toast" {
		t.Errorf("expected event name 'toast', got '%s'", events[0].Name)
	}
}

// TestContext_NoContent tests 204 response with event commitment
func TestContext_NoContent(t *testing.T) {
	// Create response recorder
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("HX-Request", "true") // Mark as HTMX request

	// Create context
	ctx := &Context{
		Res:    rec,
		Req:    req,
		events: make([]Event, 0),
	}

	// Add an event
	ctx.Event("test-event", "payload")

	// Call NoContent
	err := ctx.NoContent()
	if err != nil {
		t.Fatalf("NoContent() failed: %v", err)
	}

	// Check status code
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status code %d, got %d", http.StatusNoContent, rec.Code)
	}

	// Check HX-Trigger header was set (event commitment)
	hxTrigger := rec.Header().Get("HX-Trigger")
	if hxTrigger == "" {
		t.Error("expected HX-Trigger header to be set, got empty string")
	}
}

// TestWrapper_Wrap tests handler wrapping
func TestWrapper_Wrap(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	wrapper := NewWrapper(logger)

	t.Run("successful handler", func(t *testing.T) {
		// Create a handler that succeeds
		handler := func(ctx *Context) error {
			ctx.Event("success-event", "data")
			return ctx.Render(html.Div(nil, g.Text("Success")))
		}

		// Wrap and call
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)

		wrappedHandler := wrapper.Wrap(handler)
		wrappedHandler(rec, req)

		// Check response
		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		body := rec.Body.String()
		if !strings.Contains(body, "Success") {
			t.Errorf("expected body to contain 'Success', got '%s'", body)
		}
	})

	t.Run("handler with error", func(t *testing.T) {
		// Create a handler that returns an error
		handler := func(ctx *Context) error {
			return errors.New("test error")
		}

		// Wrap and call
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)

		wrappedHandler := wrapper.Wrap(handler)
		wrappedHandler(rec, req)

		// Check response
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", rec.Code)
		}

		body := rec.Body.String()
		if !strings.Contains(body, "Internal Server Error") {
			t.Errorf("expected body to contain 'Internal Server Error', got '%s'", body)
		}
	})
}

// TestWrapper_Logger tests logger accessor
func TestWrapper_Logger(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	wrapper := NewWrapper(logger)

	retrievedLogger := wrapper.Logger()
	if retrievedLogger != logger {
		t.Error("Logger() should return the same logger instance")
	}
}

// Test_hasPhasePrefix tests the phase prefix detection
func Test_hasPhasePrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"no prefix", "my-event", false},
		{"HX-Trigger prefix", "HX-Trigger:my-event", true},
		{"HX-Trigger-After-Swap prefix", "HX-Trigger-After-Swap:my-event", true},
		{"HX-Trigger-After-Settle prefix", "HX-Trigger-After-Settle:my-event", true},
		{"partial match", "HX-Trigger-my-event", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasPhasePrefix(tt.input)
			if result != tt.expected {
				t.Errorf("hasPhasePrefix(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
