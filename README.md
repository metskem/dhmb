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
``` 

## TODO

* - provide  response ti,me graphs using https://github.com/go-echarts/go-echarts and
       photoConfig := tgbotapi.NewDocumentUpload(chatter.ChatId, f.Name())
     	_, err := Bot.Send(photoConfig)
* v - monitor definition (data model)
* v - put chat in a (bolt?) db
* v - if update comes in with [new_chat_members](https://stackoverflow.com/questions/52271498/can-i-detect-my-bots-groups-with-telegram-bot-api): put it on a persistent list
* x - can we have the tg debug to a separate file?
* v - dynamic turn debug on/off with a bot cmd
* v - dynamic add/delete usernames with a bot cmd
* v - show current usernames