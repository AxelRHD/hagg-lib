package view

import (
	"log"
	"net/http"

	"github.com/axelrhd/hagg-lib/hxevents"
	"github.com/gin-gonic/gin"
	g "maragu.dev/gomponents"
)

// DEPRECATED: This is the Gin version. Use handler.Context.Render() for Chi.
// Will be removed in Phase 4.
func Render(ctx *gin.Context, node g.Node) {
	// Temporary fix: Pass empty events for Gin routes
	// Chi routes use handler.Context.Render() with automatic event commitment
	_ = hxevents.Commit(ctx.Writer, ctx.Request, []hxevents.Event{})

	if err := node.Render(ctx.Writer); err != nil {
		log.Println(err)
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}
