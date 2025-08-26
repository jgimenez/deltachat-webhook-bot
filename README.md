# DeltaChat webhook bot

Create webhooks to send messages on Delta Chat.

## Getting started

```sh
go install github.com/jgimenez/deltachat-webhook-bot@latest
```

This program depends on a standalone Delta Chat RPC server `deltachat-rpc-server` program that must be available in your `PATH`. For installation instructions check: https://github.com/deltachat/deltachat-core-rust/tree/master/deltachat-rpc-server

### Setting up a Delta Chat account
 * Use your favorite Delta Chat client to create an account for your bot
 * Export a backup of your profile: **Settings** > **Chats and media** > **Export Backup**. This will result in a tar file.
 * Set `DELTA_CHAT_BOT_IMPORT_ACCOUNT` variable to point to the exported profile (tar file)
 * Run `deltachat-webhook-bot`. It will listen on port 8080 by default.

### Sending a message
 * Use your Delta Chat client to connect to any users you want to send messages with (for example, `youracct@nine.testrun.org`)
 * Send them a message like this: `curl -d '{"text":"hello world!"}' localhost:8080/youracct@nine.testrun.org`

## Available configuration variables

 * `DELTA_CHAT_BOT_IMPORT_ACCOUNT`: account backup file to import. Only one account can be imported. Once an account is imported, this variable is ignored.
 * `DELTA_CHAT_BOT_LISTEN_ADDR`: listen address (default: `:8080`)

Configuration variables may be provided in the environment or in an `.env` file.

## Available endpoints

 * `POST /{destination}`: send a message to `destination` email address
   * The `destination` must be on your bot's contact list
   * The request should contain a `Content-Type: application/json; charset=utf-8` header
   * The request body should contain a JSON document with:
     * `text`: the text of the message to send
   * Expect a `204 No content` response
 
Notes:
 * The `Content-Type` header is currently not required, but it is recommended for forward compatibility.
 * The body of the JSON message is similar to Slack's, on purpose. You may use this endpoint where a Slack webhook endpoint is accepted.
 * At the moment it's only possible to send text. PRs are welcome :)