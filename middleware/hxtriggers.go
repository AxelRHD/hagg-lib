package middleware

import (
	"github.com/axelrhd/hagg-lib/hxevents"
	"github.com/gin-gonic/gin"
)

// DEPRECATED: Will be removed in Phase 4.
// This middleware is no longer needed with the new hxevents system.
func HXTriggers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		// Temporary fix: Pass empty events for Gin routes
		// Chi routes use handler.Wrapper for automatic event commitment
		_ = hxevents.Commit(ctx.Writer, ctx.Request, []hxevents.Event{})
	}
}
