# hagg-lib

Reusable building blocks for the **HAGG stack** (HTMX + Alpine.js + Gomponents + Go).

This library provides Chi-compatible packages for server-side rendered Go web applications.

---

## Status: Chi-Based (v1.0)

**Migration Complete:** This library has been migrated from Gin to Chi router.
All packages are Chi-compatible or framework-independent.

---

## Packages

### Core (Chi-Compatible)

#### **handler/** - Context Wrapper
Custom context wrapper for Chi handlers.

**Purpose:**
- Provides `handler.Context` with explicit Res/Req fields
- Fluent handler pattern: `func(*Context) error`
- Automatic error handling and event commitment

**Dependencies:** stdlib (net/http), gomponents

#### **hxevents/** - Event System
Server-driven event bus for HTMX applications.

**Purpose:**
- Emit events via HX-Trigger headers (HTMX requests)
- Emit events via initial-events script (full page loads)
- Phase support (Immediate, AfterSwap, AfterSettle)

**Dependencies:** stdlib (net/http, encoding/json), gomponents

#### **toast/** - Toast Notifications
Fluent toast notification builder.

**Purpose:**
- Create toast notifications with builder API
- Level support (success, error, warning, info)
- Timeout and position configuration
- Integrates with event system

**Dependencies:** None (uses EventEmitter interface)

### Utilities (Chi-Compatible)

#### **middleware/** - Chi Middleware
- `basepath_chi.go` - Base path injection for Chi

#### **view/** - View Helpers
- `chi.go` - Chi-compatible URL helpers (basePath-aware)

### Framework-Independent

#### **ctxkeys/** - Context Keys
Shared context key constants (e.g., BasePath).

#### **casbinx/** - Casbin Helpers
Thin helpers around Casbin integration (enforcer setup).

---

## Deprecated Packages (Phase 4 Removal)

The following packages are deprecated and will be removed after the main application completes migration:

- **flash/** - Gin-sessions based flash messages (use SCS directly)
- **middleware/basepath.go** - Gin version (use basepath_chi.go)
- **middleware/hxtriggers.go** - Old middleware (no longer needed)
- **view/render.go** - Gin version (use handler.Context.Render())
- **view/links.go** - Gin helpers (use chi.go)

---

## Installation

```bash
go get github.com/axelrhd/hagg-lib@latest
```

During development (with local hagg-lib):
```go
// go.mod
replace github.com/axelrhd/hagg-lib => ../hagg-lib
```

---

## Usage Example

```go
import (
    "github.com/axelrhd/hagg-lib/handler"
    "github.com/axelrhd/hagg-lib/toast"
    "github.com/go-chi/chi/v5"
)

func main() {
    r := chi.NewRouter()
    wrapper := handler.NewWrapper(slog.Default())

    r.Get("/", wrapper.Wrap(func(ctx *handler.Context) error {
        ctx.Toast("Welcome!").Success().Notify()
        return ctx.Render(myPage())
    }))
}
```

See [INTEGRATION.md](INTEGRATION.md) for detailed integration guide.

---

## Design Principles

- No framework magic
- Explicit imports
- Predictable behavior
- Preference for clarity over abstraction
- Chi-compatible (stdlib http.Handler)

---

## License

MIT
