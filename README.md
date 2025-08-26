# DeltaChat webhook bot

Create webhooks to send messages on Delta Chat.

## Getting started

```sh
go install github.com/jgimenez/deltachat-webhook-bot@latest
```

This program depends on a standalone Delta Chat RPC server `deltachat-rpc-server` program that must be available in your `PATH`. For installation instructions check: https://github.com/deltachat/deltachat-core-rust/tree/master/deltachat-rpc-server

### Setting up an account and sending a message
 * Use your favorite Delta Chat client to create an account for your bot
 * Export a backup of your profile: **Settings** > **Chats and media** > **Export Backup**. This will result in a tar file.
 * Set `DELTA_CHAT_BOT_IMPORT_ACCOUNT` variable to point to the exported profile (tar file)
 * Run `deltachat-webhook-bot`. It will listen on the address `:8080` unless you set `DELTA_CHAT_BOT_LISTEN_ADDR` to a different one.
 * Use your Delta Chat client to connect to any users you want to send messages with (for example, `youracct@nine.testrun.org`)
 * Send them a message like this: `curl -d '{"text":"hello world!"}' localhost:8080/youracct@nine.testrun.org`

### Available configuration variables

Configuration variables may be provided in the environment or in an `.env` file.