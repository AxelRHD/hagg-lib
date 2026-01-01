package toast

// SVG icons for toast notifications.
// These match the color scheme defined in static/css/base.css

const (
	// IconSuccess - Green checkmark icon for success toasts
	IconSuccess = `<svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
  <circle cx="10" cy="10" r="9" stroke="#43a047" stroke-width="2"/>
  <path d="M6 10l2.5 2.5L14 7" stroke="#43a047" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
</svg>`

	// IconError - Red X icon for error toasts
	IconError = `<svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
  <circle cx="10" cy="10" r="9" stroke="#e53935" stroke-width="2"/>
  <path d="M7 7l6 6M13 7l-6 6" stroke="#e53935" stroke-width="2" stroke-linecap="round"/>
</svg>`

	// IconWarning - Orange warning triangle for warning toasts
	IconWarning = `<svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
  <path d="M10 2L2 17h16L10 2z" stroke="#fb8c00" stroke-width="2" stroke-linejoin="round"/>
  <path d="M10 8v3M10 14h.01" stroke="#fb8c00" stroke-width="2" stroke-linecap="round"/>
</svg>`

	// IconInfo - Blue info icon for info toasts
	IconInfo = `<svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
  <circle cx="10" cy="10" r="9" stroke="#1095c1" stroke-width="2"/>
  <path d="M10 9v5M10 6h.01" stroke="#1095c1" stroke-width="2" stroke-linecap="round"/>
</svg>`
)

// GetIcon returns the SVG icon for the given toast level.
// Returns IconInfo if the level is not recognized.
func GetIcon(level string) string {
	switch level {
	case "success":
		return IconSuccess
	case "error":
		return IconError
	case "warning":
		return IconWarning
	case "info":
		return IconInfo
	default:
		return IconInfo
	}
}
