package middleware

// DEPRECATED: This is the Gin version. Use basepath_chi.go for Chi.
// Will be removed in Phase 4 after main app migration is complete.

import (
	"github.com/axelrhd/hagg-lib/ctxkeys"
	"github.com/gin-gonic/gin"
)

func BasePath(base string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(ctxkeys.BasePath, base)
		ctx.Next()
	}
}
