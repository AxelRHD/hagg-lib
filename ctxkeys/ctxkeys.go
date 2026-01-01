// Package ctxkeys provides shared context key constants used across hagg-lib packages.
//
// # BasePath
//
// The BasePath constant is used by middleware.BasePathChi to inject a base path
// into the request context, and by view helpers to generate basePath-aware URLs.
//
// Example:
//
//	// In main.go
//	r.Use(middleware.BasePathChi("/app"))
//
//	// In handler
//	url := view.URLStringChi(ctx.Req, "/login")  // Returns "/app/login"
//
// # Why a Separate Package?
//
// Context keys are defined in a separate package to avoid import cycles between
// middleware and view packages.
//
// # Dependencies
//
// None - stdlib only.
package ctxkeys

const BasePath = "basePath"
