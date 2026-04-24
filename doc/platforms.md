# Platform Configuration

Each platform is optional. Set its environment variables to enable it; omit them to skip it. Multiple platforms can run simultaneously — messages are handled independently per platform.

## Discord

| Variable | Description |
|----------|-------------|
| `DISCORDTOKEN` | Bot token from Discord Developer Portal |

The bot auto-joins all channels it has access to. Message splitting occurs at 2000 characters with 200 ms between chunks.

## Slack

| Variable | Description |
|----------|-------------|
| `SLACK_BOT_TOKEN` | Bot token (`xoxb-…`) from Slack app settings |
| `SLACK_APP_TOKEN` | App-level token (`xapp-…`) for Socket Mode |

Both tokens are required. Socket Mode must be enabled in the Slack app configuration.

## Telegram

| Variable | Description |
|----------|-------------|
| `TELEGRAM_BOT_TOKEN` | Token from `@BotFather` |

Default channel label is `offtopic` (used internally for routing; not an actual Telegram channel name).

## Mattermost

| Variable | Description |
|----------|-------------|
| `MATTERMOST_BOT_TOKEN` | Bot token from Mattermost admin panel |
| `MATTERMOST_URL` | Server URL, e.g. `https://your-server.com` |
| `MATTERMOST_CHANNEL` | (Optional) Default channel name |

Uses WebSocket for real-time message delivery. The URL must use `https://` — `http://` sends tokens over plaintext.

## IRC

| Variable | Format |
|----------|--------|
| `IRC_CONN` | `irc://user:password@host:port/channel` |

Example: `irc://mybot:s3cr3t@irc.libera.chat:6667/engineers-my`

The path component (after `/`) becomes the default channel.

## AI Services

These are not platforms but are required by specific feature packages:

| Variable | Package | Description |
|----------|---------|-------------|
| `SDAPI_URL` | `pkg/stablediffusion` | Structured Stable Diffusion API endpoint |
| `SD_URL` | `pkg/stablediffusion` | Simple GET-based SD endpoint (alternative to SDAPI_URL) |
| `OLLAMA_HOST` | `pkg/ollama` | Ollama server URL (defaults to `http://localhost:11434`) |

`pkg/stablediffusion` calls `log.Fatal` at startup if neither `SD_URL` nor `SDAPI_URL` is set. Remove the blank import from `main.go` if you don't want Stable Diffusion.

## Development (readline)

The `testbot` command creates a `readline` platform — it reads from stdin and writes to stdout. No tokens needed.

```bash
go run . testbot
```

Handlers can check `r.Platform == "readline"` to gate admin-only behaviours during development.

## dev.sh

`dev.sh` is sourced by the Makefile (via `include dev.sh`) and sets all env vars for a development session. Copy `dev.template.sh` to `dev.sh` and fill in your tokens. Do not commit `dev.sh`.
