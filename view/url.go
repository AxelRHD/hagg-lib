// Package view provides view helpers for generating basePath-aware URLs.
//
// These helpers work with middleware.BasePath to automatically prefix URLs
// with the configured base path, enabling deployment flexibility.
//
// # URL Helpers
//
//   - URLString: Get basePath-aware URL as string (for hx-*, forms, JS, redirects)
//
// # Usage Example
//
//	// In main.go
//	r.Use(middleware.BasePath("/app"))
//
//	// In views - internal links
//	A(Href(view.URLString(req, "/login")), g.Text("Login"))  // href="/app/login"
//
//	// In HTMX attributes
//	hx.Post(view.URLString(req, "/htmx/save"))  // hx-post="/app/htmx/save"
//
//	// In redirects
//	http.Redirect(w, r, view.URLString(req, "/"), http.StatusSeeOther)
//
// # Without BasePath Middleware
//
// If BasePath middleware is not used, these helpers return URLs unchanged.
//
// # Dependencies
//
// Requires: stdlib (net/http), ctxkeys package
package view

import (
	"net/http"

	"github.com/axelrhd/hagg-lib/ctxkeys"
)

// withBasePath extracts basePath from request context and prepends it to the path.
func withBasePath(req *http.Request, p string) string {
	raw := req.Context().Value(ctxkeys.BasePath)
	if raw == nil {
		return p
	}

	bp, ok := raw.(string)
	if !ok || bp == "/" {
		return p
	}

	return bp + p
}

// URLString returns a basePath-aware URL as a string.
//
// Use this for all internal URLs: links, HTMX endpoints, form actions, redirects.
//
// Example:
//
//	loginURL := view.URLString(req, "/htmx/login")
//	hx.Post(loginURL)
func URLString(req *http.Request, p string) string {
	return withBasePath(req, p)
}
