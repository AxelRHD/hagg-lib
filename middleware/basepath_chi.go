// Package middleware provides Chi-compatible middleware for hagg-lib.
//
// # BasePathChi
//
// BasePathChi injects a base path into the request context, enabling
// basePath-aware URL generation in view helpers.
//
// Use case: Deploy the same app at different URL prefixes without code changes.
//
// Example:
//
//	r := chi.NewRouter()
//	r.Use(middleware.BasePathChi("/app"))
//
//	// In handlers
//	url := view.URLStringChi(req, "/login")  // Returns "/app/login"
//	html := view.AChi(req, "/about", g.Text("About"))  // href="/app/about"
//
// # Dependencies
//
// Requires: stdlib (net/http, context), ctxkeys package
package middleware

import (
	"context"
	"net/http"

	"github.com/axelrhd/hagg-lib/ctxkeys"
)

// BasePathChi is a Chi middleware that injects the basePath into the request context.
// This is the Chi-compatible version of BasePath.
//
// The basePath is used by view helpers to generate basePath-aware URLs.
//
// Example:
//
//	r.Use(middleware.BasePathChi("/app"))
//	// Now view.URLStringChi(req, "/login") returns "/app/login"
func BasePathChi(base string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ctxkeys.BasePath, base)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
