### dhmb - Dirty Harry Monitor Bot

Can monitor resources at specified intervals, and optionally send updates if a threshold is exceeded.

## Resources

* [awesome Go](https://github.com/avelino/awesome-go)
* [bots: an introduction for developers](https://core.telegram.org/bots)
* [go telegram bot api](https://github.com/go-telegram-bot-api/telegram-bot-api)
* [go telegram bot api docs](https://godoc.org/github.com/go-telegram-bot-api/telegram-bot-api)

## Supported Bot Commands

The following commands are supported (and can/should be configured with BotFather (Edit Commands) for convenience):

```
status - show the current status of all monitors
members - show active chats (users and groups)
usernames - show usernames (names and their role)
debug - [on|off] - dynamically turn Telegram Bot debugging on/off
restart - restart the bot to pick up DB updates
usernameadd - <username> [reader|admin] - add a username
usernamedelete - <username> - delete a username
silence - <monname> - silence a monitor (keep monitoring, but no alerts)
unsilence - <monname> - unsilence a monitor (revert the silence)
chart - <monname> - create a response time chart for monname
``` 

## Configuration

All configuration is done with environment variables. The following envvars are available:
* **BOT_TOKEN** - the Telegram bot token, should be in the format <number>:<token>
* **MAX_ROWS_RESPTIME** - The oldest rows in the resptime table are deleted each 10 minutes. This envvar defines how many rows remain per monitor. Default = 1000
* **DEBUG** - true/false, whether the bot debug should be on or off. Default is false 

## TODO

* - make the row cleanup configurable with an envvar
* v - provide  response time graphs using https://github.com/go-echarts/go-echarts and
       photoConfig := tgbotapi.NewDocumentUpload(chatter.ChatId, f.Name())
     	_, err := Bot.Send(photoConfig)
* v - refactor the reading of envvars to a separate conf package
* v - provide background thread that cleans up the resptime table, use maxnumresptimes as an envvar
* v - option to (un)silence a monitor
* v - monitor definition (data model)
* v - put chat in a (bolt?) db
* v - if update comes in with [new_chat_members](https://stackoverflow.com/questions/52271498/can-i-detect-my-bots-groups-with-telegram-bot-api): put it on a persistent list
* x - can we have the tg debug to a separate file?
* v - dynamic turn debug on/off with a bot cmd
* v - dynamic add/delete usernames with a bot cmd
* v - show current usernames