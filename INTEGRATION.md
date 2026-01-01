# Integration Guide: Phase 1 Complete

## Summary

Phase 1 implementation is complete:
- ✅ `handler/` - Context wrapper
- ✅ `hxevents/` - Event system (fresh rewrite)
- ✅ `toast/` - Toast notification system
- ✅ Frontend JS (events.js, toast.js)

## Integration with Existing App (Phase 2)

The packages are ready to use, but the main `hagg` app still uses `gin.Context`.
Here's how to integrate when migrating to Chi:

### 1. Update skeleton.go (or equivalent layout)

**Add scripts after surreal.js:**

```go
// In Head():
view.Script(ctx, "/static/js/surreal_v1.3.4.js"),

// NEW: Event system
view.Script(ctx, "/static/js/events.js"),
view.Script(ctx, "/static/js/toast.js"),
```

**Add initial-events rendering in Body:**

```go
// In Body():
Body(
    RenderEvents(ctx),  // NEW: Initial events script
    // ... existing content
),
```

### 2. Update Layout Type

Change from `gin.Context` to `handler.Context`:

```go
// OLD:
func Skeleton(ctx *gin.Context, content ...g.Node) g.Node

// NEW:
func Skeleton(ctx *handler.Context, content ...g.Node) g.Node
```

### 3. CSS Updates

Add toast styles to `static/css/base.css` (or create if using new Tailwind setup):

```css
@layer components {
  .toast {
    @apply bg-white rounded-pico shadow-pico-lg;
    @apply border-l-4;
    @apply p-4 mb-2;
    @apply min-w-[300px] max-w-md;
    @apply transition-opacity duration-300;
  }

  .toast-success {
    @apply border-success;
  }

  .toast-error {
    @apply border-error;
  }

  .toast-warning {
    @apply border-warning;
  }

  .toast-info {
    @apply border-info;
  }
}
```

### 4. Handler Migration Pattern

**OLD (Gin):**
```go
func LoginHandler(ctx *gin.Context) {
    // ...
    notie.NewAlert("Login successful").Success().Notify(ctx)
}
```

**NEW (handler.Context):**
```go
func LoginHandler(ctx *handler.Context) error {
    // ...
    ctx.Toast("Login successful").Success().Notify()
    return nil
}
```

### 5. Router Integration (Chi)

**Setup:**
```go
// Create wrapper
wrapper := handler.NewWrapper(slog.Default())

// Register routes
r := chi.NewRouter()
r.Get("/", wrapper.Wrap(pages.Home))
```

## Testing Phase 1

To test without migrating the entire app:

1. Create a test handler with `handler.Context`
2. Register it on a test route (e.g., `/test-phase1`)
3. Test both full-page loads and HTMX requests
4. Verify toasts appear correctly

## Next Steps (Phase 2)

After validating Phase 1:
1. Migrate router (Gin → Chi)
2. Migrate all handlers to `handler.Context`
3. Migrate session management (gin-sessions → scs)
4. Remove old hxevents package
5. Remove notie package

See `REFACTORING_PLAN.md` for full migration plan.
