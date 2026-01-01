# CLAUDE.md

This file provides guidance to Claude Code when working with hagg-lib.

## Project Overview

**hagg-lib** is a reusable library providing Chi-compatible building blocks for server-side rendered Go web applications using the HAGG stack (HTMX + Alpine.js + Gomponents + Go).

**Status:** Phase 2 Complete - Gin → Chi migration finished.

## Package Structure

### Core Packages (Chi-Compatible)

**handler/** - Context wrapper for Chi handlers
- Files: `context.go`, `wrapper.go`
- Provides `handler.Context` with explicit Res/Req fields
- Handler pattern: `func(*Context) error`
- Auto error handling and event commitment

**hxevents/** - Event system for HTMX
- Files: `events.go`, `commit.go`, `context.go`, `initial.go`
- Supports HX-Trigger headers and initial-events scripts
- Phase support: Immediate, AfterSwap, AfterSettle
- Framework-independent (uses http.ResponseWriter)

**toast/** - Toast notification builder
- Files: `toast.go`, `icons.go`
- Fluent API: `.Success()`, `.Error()`, `.Warning()`, `.Info()`
- Timeout and position control
- Integrates with event system via EventEmitter interface

**middleware/** - Chi middleware
- `basepath_chi.go` - Base path injection for Chi

**view/** - View helpers
- `chi.go` - Chi-compatible URL helpers

### Framework-Independent

**ctxkeys/** - Context key constants
**casbinx/** - Casbin authorization helpers

### Deprecated (Will Remove in Phase 4)

- **flash/** - Gin-sessions based (use SCS directly instead)
- **middleware/basepath.go**, **hxtriggers.go** - Gin versions
- **view/render.go**, **links.go** - Gin helpers

## Development Guidelines

### Adding New Packages

1. Packages should be Chi-compatible or framework-independent
2. Use interfaces to avoid import cycles (see toast.EventEmitter)
3. Keep packages focused and minimal
4. Document dependencies in package comments

### Handler Pattern

```go
// Use handler.Context, not gin.Context or raw http.ResponseWriter
func MyHandler(ctx *handler.Context) error {
    // Access request
    id := chi.URLParam(ctx.Req, "id")

    // Emit events
    ctx.Event("custom-event", data)
    ctx.Toast("Success!").Success().Notify()

    // Render
    return ctx.Render(myPage())
}
```

### Event System Usage

```go
// In handler - events auto-commit via wrapper
ctx.Event("auth-changed", nil)

// Toast is just an event
ctx.Toast("Message").Success().Notify()  // emits "toast" event
```

### Avoiding Import Cycles

Use interfaces when packages need to reference each other:

```go
// In toast/toast.go - don't import handler
type EventEmitter interface {
    Event(name string, payload any)
}

// handler.Context implements EventEmitter
```

## Dependencies

**Core:**
- stdlib (net/http, encoding/json, context)
- maragu.dev/gomponents (HTML rendering)
- github.com/casbin/casbin/v2 (authorization)

**Deprecated (will remove):**
- github.com/gin-gonic/gin
- github.com/gin-contrib/sessions

## Testing

(No tests currently - add tests before v1.0.0 release)

## Versioning

Currently using replace directive for local development.
Target: v1.0.0 release after Phase 4 cleanup (remove all Gin code).

## Critical Files

- `handler/context.go` - Core context wrapper
- `handler/wrapper.go` - Wrapper that creates context and handles errors
- `hxevents/events.go` - Event definitions and storage
- `hxevents/commit.go` - Event commitment (HX-Trigger headers)
- `toast/toast.go` - Toast builder API
- `README.md` - Package documentation
