package toast

// EventEmitter is the interface for emitting events (avoids import cycle with handler package).
type EventEmitter interface {
	Event(name string, payload any)
}

// Toast represents a toast notification with fluent builder API.
// Toast notifications are sent to the frontend via the event system.
//
// The toast struct is NOT serialized directly - only its fields (excluding ctx).
// The ctx reference is used to emit the event when .Notify() is called.
type Toast struct {
	Message  string `json:"message"`  // Message text
	Level    string `json:"level"`    // success, error, warning, info
	Timeout  int    `json:"timeout"`  // Milliseconds, 0 = stay forever
	Position string `json:"position"` // bottom-right, top-right, bottom-left, top-left
	ctx      EventEmitter `json:"-"` // Context reference (not serialized)
}

// New creates a new toast builder with default values.
// Default: info level, 3 second timeout, bottom-right position.
//
// Example:
//
//	toast.New("User created", ctx).Success().Notify()
func New(msg string, ctx EventEmitter) *Toast {
	return &Toast{
		Message:  msg,
		Level:    "info",
		Timeout:  3000,
		Position: "bottom-right",
		ctx:      ctx,
	}
}

// Success sets the toast level to success (green).
// Returns self for method chaining.
func (t *Toast) Success() *Toast {
	t.Level = "success"
	return t
}

// Error sets the toast level to error (red).
// Returns self for method chaining.
func (t *Toast) Error() *Toast {
	t.Level = "error"
	return t
}

// Warning sets the toast level to warning (orange).
// Returns self for method chaining.
func (t *Toast) Warning() *Toast {
	t.Level = "warning"
	return t
}

// Info sets the toast level to info (blue).
// Returns self for method chaining.
func (t *Toast) Info() *Toast {
	t.Level = "info"
	return t
}

// SetLevel sets the toast level from a string.
// Useful for converting flash messages to toasts.
// Valid values: success, error, warning, info
// Returns self for method chaining.
func (t *Toast) SetLevel(level string) *Toast {
	t.Level = level
	return t
}

// Stay makes the toast persistent (no auto-dismiss).
// Sets timeout to 0, which tells the frontend not to auto-remove.
// Returns self for method chaining.
func (t *Toast) Stay() *Toast {
	t.Timeout = 0
	return t
}

// SetTimeout sets a custom timeout in milliseconds.
// Use 0 for no auto-dismiss (same as Stay()).
// Returns self for method chaining.
func (t *Toast) SetTimeout(ms int) *Toast {
	t.Timeout = ms
	return t
}

// SetPosition sets the toast position on screen.
// Valid values: bottom-right, top-right, bottom-left, top-left
// Returns self for method chaining.
func (t *Toast) SetPosition(pos string) *Toast {
	t.Position = pos
	return t
}

// Notify emits the toast as an event.
// The toast is sent to the frontend via:
//   - HX-Trigger header (for HTMX requests)
//   - Initial-events script (for full-page loads)
//
// Example:
//
//	ctx.Toast("Operation successful").Success().Notify()
func (t *Toast) Notify() {
	// Create a copy without the context reference for JSON serialization
	toastData := struct {
		Message  string `json:"message"`
		Level    string `json:"level"`
		Timeout  int    `json:"timeout"`
		Position string `json:"position"`
	}{
		Message:  t.Message,
		Level:    t.Level,
		Timeout:  t.Timeout,
		Position: t.Position,
	}

	// Emit as regular event (works for both HTMX and initial-events)
	t.ctx.Event("toast", toastData)
}
