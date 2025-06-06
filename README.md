# multibot

A multiplatform Bot on Discord, Slack, Telegram, Mattermost and IRC) for fun and play. It was called "discordbot" before.

## Stuff to try

`wieneruwurst is a weird word, no? I can't grep uwu in /usr/share/dict/words`

`how do i get more aws credit?`

`o/`

`hello`

`お前はもう死んでいる`

`whymca?`

`!xkcd 356`

`!explainxkcd 356`

`selamat pagi!`

`!sd close up portrait of robot`

## Where to try

* <https://t.me/EngineersMY> (Bot instance is named "angchmultibot")

* <https://discord.gg/K2NVrpvBhm>

* Slack: <https://engineersmy.slack.com/> Get invite code via <https://engineers.my/> (#random)

* irc://libera.chat/engineers-my

## Why!?

<https://austinhenley.com/blog/makinguselessstuff.html>

<https://justforfunnoreally.dev/>

## Configuration

The bot supports multiple platforms. Set the appropriate environment variables for the platforms you want to use:

### Discord
- `DISCORD_BOT_TOKEN` - Your Discord bot token

### Telegram  
- `TELEGRAM_BOT_TOKEN` - Your Telegram bot token

### Slack
- `SLACK_BOT_TOKEN` - Your Slack bot token

### Mattermost
- `MATTERMOST_BOT_TOKEN` - Your Mattermost bot token
- `MATTERMOST_URL` - Your Mattermost server URL (e.g., https://your-mattermost-server.com)

### IRC
- Configure IRC settings in your environment

## How to contribute?

1. Fork

2. hack

    ```bash
    go get ./...
    CGO_ENBALED=0 go run . # Lazy to set up rust tooling
    ```

3. test

    ```bash
    go run . testbot # Test things as a CLI
    ```

4. git

    ```bash
    git add && git commit && git push
    ```

5. submit pull request
