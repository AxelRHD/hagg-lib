// Package middleware provides HTTP middleware for hagg-lib.
//
// # BasePath
//
// BasePath injects a base path into the request context, enabling
// basePath-aware URL generation in view helpers.
//
// Use case: Deploy the same app at different URL prefixes without code changes.
//
// Example:
//
//	r := chi.NewRouter()
//	r.Use(middleware.BasePath("/app"))
//
//	// In handlers
//	url := view.URLString(req, "/login")  // Returns "/app/login"
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

// BasePath is a middleware that injects the basePath into the request context.
//
// The basePath is used by view helpers to generate basePath-aware URLs.
//
// Example:
//
//	r.Use(middleware.BasePath("/app"))
//	// Now view.URLString(req, "/login") returns "/app/login"
func BasePath(base string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ctxkeys.BasePath, base)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
