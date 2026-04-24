# Adding a New Handler

This guide shows how to add a new feature to multibot. All features are self-contained packages that register via `init()`.

## 1. Create the package

```bash
mkdir pkg/myfeature
```

Create `pkg/myfeature/myfeature.go`:

```go
package myfeature

import "github.com/angch/multibot/pkg/bothandler"

func init() {
    // Choose one or more registration functions — see below
    bothandler.RegisterCatchallHandler(myHandler)
}

func myHandler(r bothandler.Request) string {
    // r.Content  — the full message text
    // r.Platform — "discord", "slack", "telegram", "mattermost", "irc", "readline"
    // r.Channel  — channel/room identifier
    // r.From     — sender identifier
    if r.Content == "ping" {
        return "pong"
    }
    return "" // return "" to pass through to the next handler
}
```

## 2. Register the import in main.go

Add a blank import to [main.go](../main.go):

```go
_ "github.com/angch/multibot/pkg/myfeature"
```

## 3. Choose the right handler type

| Goal | Use |
|------|-----|
| React to any message (pattern matching) | `RegisterCatchallHandler` |
| React to `!mycmd arg` | `RegisterMessageWithInputHandler("!mycmd", fn)` |
| React to an exact keyword | `RegisterMessageHandler("keyword", fn)` |
| Return an image alongside text | `RegisterCatchallExtendeHandler` |
| Process uploaded image files | `RegisterImageHandler` |

### MessageWithInputHandler

`Request.Content` holds everything **after** the command keyword, trimmed.

```go
bothandler.RegisterMessageWithInputHandler("!hello", func(r bothandler.Request) string {
    if r.Content == "" {
        return "Hello, world!"
    }
    return "Hello, " + r.Content + "!"
})
```

### CatchallExtendedHandler

Return `nil` to pass through. Return an `*ExtendedMessage` to respond with text and/or image bytes.

```go
bothandler.RegisterCatchallExtendeHandler(func(msg bothandler.ExtendedMessage) *bothandler.ExtendedMessage {
    if !strings.HasPrefix(msg.Text, "!img ") {
        return nil
    }
    imgBytes := generateImage(msg.Text[5:])
    return &bothandler.ExtendedMessage{Text: "Here you go", Image: imgBytes}
})
```

### ImageHandler

Called when a message has an attachment. `filename` is the path to a temp file containing the downloaded image.

```go
bothandler.RegisterImageHandler(func(filename string, r bothandler.Request) string {
    // process filename, return response string or ""
    return ""
})
```

## 4. Test it

```bash
go run . testbot
```

The `testbot` command starts an interactive readline session on the `readline` platform — no tokens needed. Type messages and see responses.

## 5. Tips

- Return `""` / `nil` when your handler doesn't apply. All other handlers still run.
- Use `r.Platform` to restrict behaviour to specific platforms.
- Catchall handlers run on **every** message; keep them fast. Defer slow work (HTTP calls, LLM inference) only when the message clearly matches your feature.
- If you need a background goroutine (scheduled tasks, polling), start it in `init()` and use `bothandler.ActiveMessagePlatforms` to broadcast.
