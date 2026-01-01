package toast

import (
	"testing"
)

// mockEventEmitter implements EventEmitter for testing
type mockEventEmitter struct {
	events []struct {
		name    string
		payload any
	}
}

func (m *mockEventEmitter) Event(name string, payload any) {
	m.events = append(m.events, struct {
		name    string
		payload any
	}{name, payload})
}

// TestNew tests toast creation with default values
func TestNew(t *testing.T) {
	ctx := &mockEventEmitter{}
	toast := New("Test message", ctx)

	if toast.Message != "Test message" {
		t.Errorf("expected message 'Test message', got '%s'", toast.Message)
	}
	if toast.Level != "info" {
		t.Errorf("expected default level 'info', got '%s'", toast.Level)
	}
	if toast.Timeout != 3000 {
		t.Errorf("expected default timeout 3000, got %d", toast.Timeout)
	}
	if toast.Position != "bottom-right" {
		t.Errorf("expected default position 'bottom-right', got '%s'", toast.Position)
	}
	if toast.ctx != ctx {
		t.Error("context reference should be set")
	}
}

// TestLevelSetters tests all level setter methods
func TestLevelSetters(t *testing.T) {
	tests := []struct {
		name     string
		setter   func(*Toast) *Toast
		expected string
	}{
		{"Success", (*Toast).Success, "success"},
		{"Error", (*Toast).Error, "error"},
		{"Warning", (*Toast).Warning, "warning"},
		{"Info", (*Toast).Info, "info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &mockEventEmitter{}
			toast := New("Test", ctx)
			result := tt.setter(toast)

			if toast.Level != tt.expected {
				t.Errorf("expected level '%s', got '%s'", tt.expected, toast.Level)
			}

			// Check method chaining
			if result != toast {
				t.Error("setter should return self for chaining")
			}
		})
	}
}

// TestSetLevel tests custom level setting
func TestSetLevel(t *testing.T) {
	ctx := &mockEventEmitter{}
	toast := New("Test", ctx)

	result := toast.SetLevel("custom-level")

	if toast.Level != "custom-level" {
		t.Errorf("expected level 'custom-level', got '%s'", toast.Level)
	}

	// Check method chaining
	if result != toast {
		t.Error("SetLevel should return self for chaining")
	}
}

// TestStay tests persistent toast setting
func TestStay(t *testing.T) {
	ctx := &mockEventEmitter{}
	toast := New("Test", ctx)

	// Initially has timeout
	if toast.Timeout != 3000 {
		t.Errorf("expected initial timeout 3000, got %d", toast.Timeout)
	}

	result := toast.Stay()

	if toast.Timeout != 0 {
		t.Errorf("expected timeout 0 after Stay(), got %d", toast.Timeout)
	}

	// Check method chaining
	if result != toast {
		t.Error("Stay should return self for chaining")
	}
}

// TestSetTimeout tests custom timeout setting
func TestSetTimeout(t *testing.T) {
	ctx := &mockEventEmitter{}
	toast := New("Test", ctx)

	result := toast.SetTimeout(5000)

	if toast.Timeout != 5000 {
		t.Errorf("expected timeout 5000, got %d", toast.Timeout)
	}

	// Check method chaining
	if result != toast {
		t.Error("SetTimeout should return self for chaining")
	}

	// Test setting to 0 (same as Stay())
	toast.SetTimeout(0)
	if toast.Timeout != 0 {
		t.Errorf("expected timeout 0, got %d", toast.Timeout)
	}
}

// TestSetPosition tests position setting
func TestSetPosition(t *testing.T) {
	positions := []string{"bottom-right", "top-right", "bottom-left", "top-left"}

	for _, pos := range positions {
		t.Run(pos, func(t *testing.T) {
			ctx := &mockEventEmitter{}
			toast := New("Test", ctx)

			result := toast.SetPosition(pos)

			if toast.Position != pos {
				t.Errorf("expected position '%s', got '%s'", pos, toast.Position)
			}

			// Check method chaining
			if result != toast {
				t.Error("SetPosition should return self for chaining")
			}
		})
	}
}

// TestNotify tests event emission
func TestNotify(t *testing.T) {
	ctx := &mockEventEmitter{}
	toast := New("Test notification", ctx)
	toast.Success().SetTimeout(5000).SetPosition("top-right")

	toast.Notify()

	// Check event was emitted
	if len(ctx.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(ctx.events))
	}

	event := ctx.events[0]
	if event.name != "toast" {
		t.Errorf("expected event name 'toast', got '%s'", event.name)
	}

	// Check payload structure
	payload, ok := event.payload.(struct {
		Message  string `json:"message"`
		Level    string `json:"level"`
		Timeout  int    `json:"timeout"`
		Position string `json:"position"`
	})
	if !ok {
		t.Fatal("payload should be toast data struct")
	}

	if payload.Message != "Test notification" {
		t.Errorf("expected message 'Test notification', got '%s'", payload.Message)
	}
	if payload.Level != "success" {
		t.Errorf("expected level 'success', got '%s'", payload.Level)
	}
	if payload.Timeout != 5000 {
		t.Errorf("expected timeout 5000, got %d", payload.Timeout)
	}
	if payload.Position != "top-right" {
		t.Errorf("expected position 'top-right', got '%s'", payload.Position)
	}
}

// TestFluentAPI tests method chaining workflow
func TestFluentAPI(t *testing.T) {
	tests := []struct {
		name     string
		builder  func(*Toast) *Toast
		expected struct {
			level    string
			timeout  int
			position string
		}
	}{
		{
			name: "success toast with custom timeout",
			builder: func(toast *Toast) *Toast {
				return toast.Success().SetTimeout(2000)
			},
			expected: struct {
				level    string
				timeout  int
				position string
			}{"success", 2000, "bottom-right"},
		},
		{
			name: "error toast that stays",
			builder: func(toast *Toast) *Toast {
				return toast.Error().Stay()
			},
			expected: struct {
				level    string
				timeout  int
				position string
			}{"error", 0, "bottom-right"},
		},
		{
			name: "warning toast in top-left",
			builder: func(toast *Toast) *Toast {
				return toast.Warning().SetPosition("top-left")
			},
			expected: struct {
				level    string
				timeout  int
				position string
			}{"warning", 3000, "top-left"},
		},
		{
			name: "info toast with all custom settings",
			builder: func(toast *Toast) *Toast {
				return toast.Info().SetTimeout(10000).SetPosition("bottom-left")
			},
			expected: struct {
				level    string
				timeout  int
				position string
			}{"info", 10000, "bottom-left"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &mockEventEmitter{}
			toast := New("Fluent API test", ctx)

			result := tt.builder(toast)

			// Check method chaining returns same instance
			if result != toast {
				t.Error("fluent methods should return self")
			}

			// Check final values
			if toast.Level != tt.expected.level {
				t.Errorf("expected level '%s', got '%s'", tt.expected.level, toast.Level)
			}
			if toast.Timeout != tt.expected.timeout {
				t.Errorf("expected timeout %d, got %d", tt.expected.timeout, toast.Timeout)
			}
			if toast.Position != tt.expected.position {
				t.Errorf("expected position '%s', got '%s'", tt.expected.position, toast.Position)
			}
		})
	}
}

// TestGetIcon tests icon retrieval
func TestGetIcon(t *testing.T) {
	tests := []struct {
		level    string
		expected string
	}{
		{"success", IconSuccess},
		{"error", IconError},
		{"warning", IconWarning},
		{"info", IconInfo},
		{"unknown", IconInfo}, // Default fallback
		{"", IconInfo},        // Empty string fallback
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			icon := GetIcon(tt.level)
			if icon != tt.expected {
				t.Errorf("GetIcon(%q) returned wrong icon", tt.level)
			}
		})
	}
}

// TestIconConstants tests that icon constants are not empty
func TestIconConstants(t *testing.T) {
	icons := map[string]string{
		"IconSuccess": IconSuccess,
		"IconError":   IconError,
		"IconWarning": IconWarning,
		"IconInfo":    IconInfo,
	}

	for name, icon := range icons {
		if icon == "" {
			t.Errorf("%s should not be empty", name)
		}
		// Basic SVG validation
		if len(icon) < 10 {
			t.Errorf("%s should be a valid SVG string, got length %d", name, len(icon))
		}
	}
}

// TestMultipleNotify tests calling Notify multiple times
func TestMultipleNotify(t *testing.T) {
	ctx := &mockEventEmitter{}
	toast := New("Test", ctx).Success()

	// Call Notify multiple times
	toast.Notify()
	toast.Notify()
	toast.Notify()

	if len(ctx.events) != 3 {
		t.Errorf("expected 3 events after 3 Notify() calls, got %d", len(ctx.events))
	}

	// All should be identical
	for i, event := range ctx.events {
		if event.name != "toast" {
			t.Errorf("event %d: expected name 'toast', got '%s'", i, event.name)
		}
	}
}

// TestContextReference tests that context is not included in payload
func TestContextReference(t *testing.T) {
	ctx := &mockEventEmitter{}
	toast := New("Test", ctx)

	toast.Notify()

	event := ctx.events[0]
	payload := event.payload

	// Use type assertion to check the payload structure
	// It should NOT have a ctx field
	switch v := payload.(type) {
	case struct {
		Message  string `json:"message"`
		Level    string `json:"level"`
		Timeout  int    `json:"timeout"`
		Position string `json:"position"`
	}:
		// This is expected - anonymous struct without ctx
		if v.Message != "Test" {
			t.Error("payload should have correct message")
		}
	default:
		t.Errorf("payload has unexpected type: %T", payload)
	}
}
