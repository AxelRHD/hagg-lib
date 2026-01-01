package hxevents

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestIsHtmxRequest tests HTMX request detection
func TestIsHtmxRequest(t *testing.T) {
	tests := []struct {
		name     string
		setupHeader func() http.Header
		expected bool
	}{
		{
			name: "HTMX request",
			setupHeader: func() http.Header {
				h := make(http.Header)
				h.Set("HX-Request", "true")
				return h
			},
			expected: true,
		},
		{
			name: "non-HTMX request",
			setupHeader: func() http.Header {
				return make(http.Header)
			},
			expected: false,
		},
		{
			name: "HX-Request with wrong value",
			setupHeader: func() http.Header {
				h := make(http.Header)
				h.Set("HX-Request", "false")
				return h
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := tt.setupHeader()
			result := IsHtmxRequest(headers)
			if result != tt.expected {
				t.Errorf("IsHtmxRequest() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestCommit tests event commitment to HX-Trigger headers
func TestCommit(t *testing.T) {
	t.Run("non-HTMX request - no headers set", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		// No HX-Request header

		events := []Event{
			{Name: "HX-Trigger:test-event", Payload: "data"},
		}

		err := Commit(rec, req, events)
		if err != nil {
			t.Fatalf("Commit() failed: %v", err)
		}

		// No headers should be set for non-HTMX requests
		if rec.Header().Get("HX-Trigger") != "" {
			t.Error("HX-Trigger header should not be set for non-HTMX request")
		}
	})

	t.Run("HTMX request - immediate event", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("HX-Request", "true")

		events := []Event{
			{Name: "HX-Trigger:test-event", Payload: "test-data"},
		}

		err := Commit(rec, req, events)
		if err != nil {
			t.Fatalf("Commit() failed: %v", err)
		}

		// Check HX-Trigger header
		hxTrigger := rec.Header().Get("HX-Trigger")
		if hxTrigger == "" {
			t.Fatal("HX-Trigger header should be set")
		}

		// Parse JSON
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(hxTrigger), &data); err != nil {
			t.Fatalf("failed to parse HX-Trigger JSON: %v", err)
		}

		// Check event data
		if data["test-event"] != "test-data" {
			t.Errorf("expected event payload 'test-data', got %v", data["test-event"])
		}
	})

	t.Run("HTMX request - multiple phases", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("HX-Request", "true")

		events := []Event{
			{Name: "HX-Trigger:immediate", Payload: "immediate-data"},
			{Name: "HX-Trigger-After-Swap:after-swap", Payload: "swap-data"},
			{Name: "HX-Trigger-After-Settle:after-settle", Payload: "settle-data"},
		}

		err := Commit(rec, req, events)
		if err != nil {
			t.Fatalf("Commit() failed: %v", err)
		}

		// Check all three headers
		headers := []string{"HX-Trigger", "HX-Trigger-After-Swap", "HX-Trigger-After-Settle"}
		expectedEvents := []string{"immediate", "after-swap", "after-settle"}

		for i, header := range headers {
			value := rec.Header().Get(header)
			if value == "" {
				t.Errorf("%s header should be set", header)
				continue
			}

			var data map[string]interface{}
			if err := json.Unmarshal([]byte(value), &data); err != nil {
				t.Errorf("failed to parse %s JSON: %v", header, err)
				continue
			}

			if _, ok := data[expectedEvents[i]]; !ok {
				t.Errorf("%s should contain event '%s'", header, expectedEvents[i])
			}
		}
	})

	t.Run("HTMX request - multiple events same phase", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("HX-Request", "true")

		events := []Event{
			{Name: "HX-Trigger:event1", Payload: "data1"},
			{Name: "HX-Trigger:event2", Payload: "data2"},
		}

		err := Commit(rec, req, events)
		if err != nil {
			t.Fatalf("Commit() failed: %v", err)
		}

		hxTrigger := rec.Header().Get("HX-Trigger")
		if hxTrigger == "" {
			t.Fatal("HX-Trigger header should be set")
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(hxTrigger), &data); err != nil {
			t.Fatalf("failed to parse HX-Trigger JSON: %v", err)
		}

		// Both events should be present
		if len(data) != 2 {
			t.Errorf("expected 2 events, got %d", len(data))
		}
		if data["event1"] != "data1" {
			t.Error("event1 should have payload 'data1'")
		}
		if data["event2"] != "data2" {
			t.Error("event2 should have payload 'data2'")
		}
	})

	t.Run("event without phase prefix - ignored", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("HX-Request", "true")

		events := []Event{
			{Name: "no-prefix-event", Payload: "data"},
		}

		err := Commit(rec, req, events)
		if err != nil {
			t.Fatalf("Commit() failed: %v", err)
		}

		// No headers should be set for events without phase prefix
		if rec.Header().Get("HX-Trigger") != "" {
			t.Error("HX-Trigger header should not be set for events without phase prefix")
		}
	})

	t.Run("complex payload", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("HX-Request", "true")

		complexPayload := map[string]interface{}{
			"message": "Hello",
			"level":   "success",
			"timeout": 3000,
		}

		events := []Event{
			{Name: "HX-Trigger:toast", Payload: complexPayload},
		}

		err := Commit(rec, req, events)
		if err != nil {
			t.Fatalf("Commit() failed: %v", err)
		}

		hxTrigger := rec.Header().Get("HX-Trigger")
		if hxTrigger == "" {
			t.Fatal("HX-Trigger header should be set")
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(hxTrigger), &data); err != nil {
			t.Fatalf("failed to parse HX-Trigger JSON: %v", err)
		}

		// Check complex payload structure
		toast, ok := data["toast"].(map[string]interface{})
		if !ok {
			t.Fatal("toast event should be a map")
		}

		if toast["message"] != "Hello" {
			t.Errorf("expected message 'Hello', got %v", toast["message"])
		}
	})
}

// TestPhaseConstants tests phase constant values
func TestPhaseConstants(t *testing.T) {
	if Immediate != "HX-Trigger" {
		t.Errorf("Immediate phase should be 'HX-Trigger', got '%s'", Immediate)
	}
	if AfterSwap != "HX-Trigger-After-Swap" {
		t.Errorf("AfterSwap phase should be 'HX-Trigger-After-Swap', got '%s'", AfterSwap)
	}
	if AfterSettle != "HX-Trigger-After-Settle" {
		t.Errorf("AfterSettle phase should be 'HX-Trigger-After-Settle', got '%s'", AfterSettle)
	}
}
