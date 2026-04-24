# Bot Commands Reference

Commands a user can send to the bot across any supported platform.

## Explicit Commands

| Command | Description |
|---------|-------------|
| `!xkcd <number>` | Link to an XKCD comic by number |
| `!explainxkcd <number>` | Link to the explainxkcd wiki entry |
| `!sd <prompt>` | Generate an AI image via Stable Diffusion |
| `!oll <question>` | Ask the Ollama LLM a question (brief, direct response) |
| `!qrcode <text>` or `/qrcode <text>` | Generate a QR code image |

## Persona Triggers (Ollama)

These phrases anywhere in a message activate a specific LLM persona:

| Trigger word | Persona |
|--------------|---------|
| `murderbot` | SecUnit — sarcastic, terse, darkly comedic |
| `demurebot` | Polite, reserved, formal |
| `angrybot` | Irritable, blunt, sarcastic |
| `depressedbot` | Melancholic, pessimistic (Marvin-inspired) |

## Phrase / Fragment Triggers (echo package)

These are pattern-matched against message content. The bot replies when the phrase appears in a message short enough relative to the phrase length.

| Example trigger | Notes |
|----------------|-------|
| `uwu` (in a word like `wieneruwurst`) | Triggers uwu transformation response |
| `hello` | Greeting response |
| `o/` | Wave greeting |
| `selamat pagi` | Malay morning greeting |
| `how do i get more aws credit` | Community in-joke response |
| `お前はもう死んでいる` | Anime meme response (`nani?!`) |
| `whymca` | YMCA ASCII art |

See [pkg/echo/data.go](../pkg/echo/data.go) for the full list.

## Image Processing (automatic)

When you send an image attachment, these handlers run automatically:

| Handler | Package | What it does |
|---------|---------|--------------|
| QR decode | `pkg/qrdecode` | Decodes QR codes in the image |
| OCR | `pkg/ocr` | Extracts text via Tesseract |
| Ollama vision | `pkg/ollama` | Describes the image if `!oll` prefixes the caption |

## SpaceTraders Commands

Available when the SpaceTraders package is loaded:

| Command | Description |
|---------|-------------|
| `!register <callsign>` | Register a new SpaceTraders agent |
| `!st status` | Show agent status |

See [SPACETRADERS.md](../SPACETRADERS.md) for the full SpaceTraders integration docs.

## APOD (Astronomy Picture of the Day)

Automatic — no command needed. The bot posts NASA's Astronomy Picture of the Day to configured channels once per day. Channel IDs are set in `~/.multibot.yaml`.
