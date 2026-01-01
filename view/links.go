package view

// DEPRECATED: These are Gin helpers. Use chi.go for Chi-compatible helpers.
// Will be removed in Phase 4.

import (
	"github.com/axelrhd/hagg-lib/ctxkeys"
	"github.com/gin-gonic/gin"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func withBasePath(ctx *gin.Context, p string) string {
	raw, ok := ctx.Get(ctxkeys.BasePath)
	if !ok {
		return p
	}

	bp, ok := raw.(string)
	if !ok || bp == "/" {
		return p
	}

	return bp + p
}

func A(ctx *gin.Context, href string, nodes ...g.Node) g.Node {
	return html.A(
		html.Href(withBasePath(ctx, href)),
		g.Group(nodes),
	)
}

func Script(ctx *gin.Context, src string, nodes ...g.Node) g.Node {
	return html.Script(
		html.Src(withBasePath(ctx, src)),
		g.Group(nodes),
	)
}

func Stylesheet(ctx *gin.Context, href string, nodes ...g.Node) g.Node {
	return html.Link(
		html.Rel("stylesheet"),
		html.Href(withBasePath(ctx, href)),
		g.Group(nodes),
	)
}

// URLString returns a basePath-aware URL as a string.
// Intended for hx-* attributes, form actions, redirects, JS, etc.
func URLString(ctx *gin.Context, p string) string {
	return withBasePath(ctx, p)
}
