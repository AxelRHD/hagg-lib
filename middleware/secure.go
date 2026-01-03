package middleware

import "net/http"

// Secure is a middleware that sets common security headers.
// It is recommended to use this middleware in all production deployments.
//
// Headers set:
//   - X-Frame-Options: DENY (prevents clickjacking)
//   - X-Content-Type-Options: nosniff (prevents MIME sniffing)
//   - X-XSS-Protection: 1; mode=block (legacy XSS filter)
//   - Referrer-Policy: strict-origin-when-cross-origin (limits referrer info)
//
// Usage:
//
// Add early in the middleware chain, after session middleware:
//
//	import libmw "github.com/axelrhd/hagg-lib/middleware"
//
//	r := chi.NewRouter()
//	r.Use(chimw.RealIP)
//	r.Use(chimw.Compress(5))
//	r.Use(session.Manager.LoadAndSave)
//	r.Use(libmw.BasePath("/app"))  // if needed
//	r.Use(libmw.Secure)            // recommended
//
// Note: These headers are safe defaults. X-Frame-Options: DENY prevents
// embedding in iframes - if you need iframe support, consider a custom
// implementation with SAMEORIGIN or CSP frame-ancestors.
func Secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}
