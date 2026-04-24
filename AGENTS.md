# AGENTS.md

This file provides guidance to AI coding agents when working in this repository.

## Commands

```bash
# Build (includes Rust uwuify library)
make build

# Run tests + govulncheck vulnerability scan
make test

# Run a single package's tests
go test ./pkg/echo/...

# Run bot in dev loop (restarts every 10 min)
make run

# Interactive CLI testing — no platform tokens needed
go run . testbot

# Build without Rust CGO requirement
CGO_ENABLED=0 go build .

# Nil pointer analysis
make nilaway

# Update Go modules + Rust deps
make updatemod
```

`dev.sh` contains environment variables (tokens) sourced by the Makefile via `include dev.sh`. Copy `dev.template.sh` to `dev.sh` to get started.

## Architecture Summary

Multibot is a multi-platform chatbot (Discord, Slack, Telegram, Mattermost, IRC).

**Plugin system:** feature packages are blank-imported in `main.go`; each `init()` registers handlers with the global registry in `pkg/bothandler/struct.go`. To add a feature, create a package with an `init()` that calls the appropriate `Register*` function, then add the blank import.

**Handler dispatch:** every inbound message is tried against all registered handlers in sequence (not first-match-wins). A handler returns `""` or `nil` to pass through.

**Platform adapters:** `cmd/run.go` reads env vars, instantiates each platform, and runs `go platform.ProcessMessages()` in its own goroutine. Graceful shutdown on SIGINT/SIGTERM via `bothandler.Shutdown()`.

**Rust FFI:** `lib/uwu/` is a Rust cdylib (`libuwu.so`) called by `pkg/echo/` via CGo. Set `CGO_ENABLED=0` during development to skip it; a no-op stub is provided.

## Documentation

- [Architecture](doc/architecture.md) — handler types, request flow, platform adapter interface, background services
- [Adding a Handler](doc/adding-a-handler.md) — step-by-step guide to creating a new plugin
- [Platform Configuration](doc/platforms.md) — environment variables for each platform and AI service
- [Commands Reference](doc/commands.md) — user-facing bot commands and phrase triggers
- [Security Findings](TODO.md) — audit log of known security issues (do not introduce new ones in the same categories)
- [SpaceTraders Integration](SPACETRADERS.md) — SpaceTraders sub-binary and debugging guide
