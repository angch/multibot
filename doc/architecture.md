# Architecture

## Overview

Multibot is a multi-platform chatbot (Discord, Slack, Telegram, Mattermost, IRC) built with a plugin-based handler system. Each feature registers itself via `init()` and the core framework dispatches messages across all registered handlers.

## Plugin Registration

All feature packages are imported as blank imports in [main.go](../main.go). Each package's `init()` registers its handlers with the global registry in [pkg/bothandler/struct.go](../pkg/bothandler/struct.go). To add a feature: create a package, register in `init()`, and add a blank import in `main.go`.

## Handler Types

Defined in [pkg/bothandler/struct.go](../pkg/bothandler/struct.go):

| Type | Registration function | Signature | When invoked |
|------|-----------------------|-----------|--------------|
| `MessageHandler` | `RegisterMessageHandler(keyword, fn)` | `func() string` | Exact keyword match |
| `MessageWithInputHandler` | `RegisterMessageWithInputHandler("!cmd", fn)` | `func(Request) string` | `!cmd args` — `Request.Content` holds args |
| `CatchallHandler` | `RegisterCatchallHandler(fn)` | `func(Request) string` | Every message, in registration order |
| `CatchallExtendedHandler` | `RegisterCatchallExtendeHandler(fn)` | `func(ExtendedMessage) *ExtendedMessage` | Every message; can return text + image bytes |
| `ImageHandler` | `RegisterImageHandler(fn)` | `func(string, Request) string` | Messages with file attachments; first arg is temp file path |

Handlers are **not** first-match-wins. All applicable handlers in each category are tried in order. A handler signals "not applicable" by returning `""` or `nil`.

## Request Flow

```
Platform event received
        │
        ▼
Exact keyword lookup (Handlers map)
        │ if miss
        ▼
CatchallHandlers (in order)
        │
        ▼
CatchallExtendedHandlers (in order, can return image)
        │
        ▼
MsgInputHandlers (if message matches "!cmd" prefix)
        │
        ▼
ImageHandlers (if message has attachments)
        │
        ▼
Response sent back to originating channel
```

Long responses are split at platform-specific limits (Discord: 2000 chars) with a 200 ms delay between chunks.

## Platform Adapters

Each platform lives in [pkg/bothandler/](../pkg/bothandler/) and implements `MessagePlatform`:

```go
type MessagePlatform interface {
    Send(string)
    SendWithOptions(string, SendOptions)
    ProcessMessages()
    Close()
    ChannelMessageSend(channel string, message string) error
}
```

`cmd/run.go` reads env vars, instantiates each platform, and calls `go platform.ProcessMessages()` in its own goroutine. A `SIGINT`/`SIGTERM` triggers `bothandler.Shutdown()` which calls `Close()` on every registered platform.

## Background Services

Some packages start goroutines in `init()` independent of message flow:

- **`pkg/apod/`** — polls NASA APOD on a daily schedule and broadcasts to all active platforms via `bothandler.ActiveMessagePlatforms`.

## Rust FFI

[lib/uwu/](../lib/uwu/) is a Rust `cdylib` exporting `uwuify()` over C FFI, compiled to `lib/uwu/target/release/libuwu.so`. Used by [pkg/echo/](../pkg/echo/) via CGo. Set `CGO_ENABLED=0` to disable; `pkg/echo/echo-nocgo.go` provides a no-op stub so the package compiles without the library.

## Data Flow for AI Features

- **Ollama** (`pkg/ollama/`): registered as both `CatchallHandler` and `ImageHandler`. Triggers on `!oll` prefix or persona keywords (`murderbot`, `demurebot`, `angrybot`, `depressedbot`). Connects to Ollama API at `http://localhost:11434` by default.
- **Stable Diffusion** (`pkg/stablediffusion/`): registered as `CatchallExtendedHandler`. Triggers on `!sd ` prefix. Returns raw image bytes. Supports two API modes: `SDAPI_URL` (structured API) or `SD_URL` (simple GET endpoint).

## SpaceTraders Sub-Binary

`spacetraders/` is an alternative `main` package that includes only SpaceTraders-related plugins — useful for isolated debugging. Run with `go run ./spacetraders testbot`.
