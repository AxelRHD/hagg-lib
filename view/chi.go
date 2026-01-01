// Package view provides Chi-compatible view helpers for generating basePath-aware URLs.
//
// These helpers work with middleware.BasePathChi to automatically prefix URLs
// with the configured base path, enabling deployment flexibility.
//
// # URL Helpers
//
//   - AChi: Create <a> tags with basePath-aware hrefs
//   - ScriptChi: Create <script> tags with basePath-aware src
//   - StylesheetChi: Create <link rel="stylesheet"> with basePath-aware href
//   - URLStringChi: Get basePath-aware URL as string (for hx-*, forms, JS)
//
// # Usage Example
//
//	// In main.go
//	r.Use(middleware.BasePathChi("/app"))
//
//	// In views
//	view.AChi(req, "/login", g.Text("Login"))  // <a href="/app/login">Login</a>
//	view.ScriptChi(req, "/static/app.js")      // <script src="/app/static/app.js"></script>
//
//	// In HTMX attributes
//	html.Div(
//	    g.Attr("hx-get", view.URLStringChi(req, "/api/data")),  // hx-get="/app/api/data"
//	    g.Text("Click me"),
//	)
//
// # Without BasePathChi Middleware
//
// If BasePathChi middleware is not used, these helpers return URLs unchanged.
//
// # Dependencies
//
// Requires: stdlib (net/http), gomponents, ctxkeys package
package view

import (
	"net/http"

	"github.com/axelrhd/hagg-lib/ctxkeys"
)

// withBasePathChi extracts basePath from request context (Chi version).
func withBasePathChi(req *http.Request, p string) string {
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

// URLStringChi returns a basePath-aware URL as a string for Chi routes.
// This is the Chi-compatible version of URLString.
//
// Intended for hx-* attributes, form actions, redirects, JS, etc.
//
// Example:
//
//	loginURL := view.URLStringChi(ctx.Req, "/htmx/login")
//	Form(hx.Post(loginURL), ...)
func URLStringChi(req *http.Request, p string) string {
	return withBasePathChi(req, p)
}
