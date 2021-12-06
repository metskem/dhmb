### dhmb - Dirty Harry Monitor Bot

Can monitor http resources at specified intervals, and optionally send Telegram updates if a threshold is exceeded.  
Can partly be operated using Telegram.  
Provides a prometheus exporter for response time metrics.

## Resources

* [awesome Go](https://github.com/avelino/awesome-go)
* [bots: an introduction for developers](https://core.telegram.org/bots)
* [go telegram bot api](https://github.com/go-telegram-bot-api/telegram-bot-api)
* [go telegram bot api docs](https://godoc.org/github.com/go-telegram-bot-api/telegram-bot-api)

## Supported Bot Commands

The following commands are supported (and can/should be configured with BotFather (Edit Commands) for convenience):
```
status - show the current status of all monitors
chart - <monname> - create a response time chart for monname
members - show active chats (users and groups)
usernames - show usernames (names and their role)
debug - [on|off] - dynamically turn Telegram Bot debugging on/off
restart - restart the bot to pick up DB updates
usernameadd - <username> [reader|admin] - add a username
usernamedelete - <username> - delete a username
silence - <monname> - silence a monitor (keep monitoring, but no alerts)
unsilence - <monname> - unsilence a monitor (revert the silence)
``` 

## Configuration

All configuration is done with environment variables. The following envvars are available:
* **BOT_TOKEN** - the Telegram bot token, should be in the format `<number>:<token>`
* **MAX_ROWS_RESPTIME** - The oldest rows in the resptime table are deleted each 10 minutes. This envvar defines how many rows remain per monitor. Default = 1000
* **DEBUG** - true/false, whether the bot debug should be on or off. Default is false
* **PROMETHEUS_EXPORTER_PORT** - The port to use for exposing the prometheus exporter (default is 9094). The metric `dhmb_resptime`, with labels `name` and `status` is exposed.
