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
