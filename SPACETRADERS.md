# Dev

SpaceTraders is a discordbot/multibot package that *WILL* talk to
`https://spacetraders.io/`

There is an alternative main package to discordbot that links and
activates only selected pkg for spacetrader, to make debugging much
easier.

The following will run the full bot:

```bash
go run .
```

But this runs only the relevant spacetraders modules:

```bash
go run ./spacetraders
```

An interactive dev looks like this, running on the readline dev platform:

```bash
user@server:~/discordbot$ go run  ./spacetraders testbot
2023/07/10 13:43:12 pkg/spacetraders/init
Test Bot is now running.  Press CTRL-C to exit.
» hello
2023/07/10 13:43:15 pkg/spacetraders/SpaceTradersHandler {Content:hello Platform:readline Channel: From:}
```

You may need to set the following env vars:

```bash
export LD_LIBRARY_PATH=`pwd`/lib/uwu/target/release
```

to silence the annoying Rust linker.

Or just turn oFF `CGO_ENABLED=0`:

```bash
export CGO_ENABLED=0
```
