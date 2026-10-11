# Twitch Tarot Bot

A Twitch chat bot that draws tarot cards.

| Command | Reply |
| --- | --- |
| `!tarot` | One card, upright or reversed, with its meaning. |
| `!spread` | Three cards: Past, Present, Future. |

## Settings

Copy `.env.example` to `.env` and fill in the three the bot needs:

| Variable | What it is |
| --- | --- |
| `TWITCH_BOT_USERNAME` | The account the bot posts as. |
| `TWITCH_OAUTH_TOKEN` | That account's chat token, starting with `oauth:`. |
| `TWITCH_CHANNEL` | The channel to join. |

The token expires after about 60 days. When it does the bot exits with a login error; generate a new one and replace it.

## Run locally

```bash
go run ./cmd/bot
```

```bash
go test ./...
```

Or in a container:

```bash
docker build -t tarot-bot .
```

```bash
docker run --rm --env-file .env tarot-bot
```

## Deploy to Fly.io

The bot runs as its own Fly app on one small Machine. It only dials out to Twitch, so `fly.toml` has no HTTP service. The deck is compiled into the binary, so there is no volume either.

Create the app from the checked-in config. Change `app` in `fly.toml` first if the name is taken:

```bash
fly launch --no-deploy --copy-config
```

Set the three settings as secrets:

```bash
fly secrets set TWITCH_BOT_USERNAME=... TWITCH_CHANNEL=... TWITCH_OAUTH_TOKEN=oauth:...
```

Deploy one Machine. `--ha=false` matters: a second Machine would answer every command twice.

```bash
fly deploy --ha=false
```

Check on it:

```bash
fly logs
```

`Connected to Twitch chat.` means it is in. If the token is rejected the Machine stops after five restarts and `fly status` shows it stopped; set a new `TWITCH_OAUTH_TOKEN` with `fly secrets set`, which redeploys it.
