package main

import (
	tbgot "bodygraph-bot/pkg/bot"
	"bodygraph-bot/pkg/config"
)

func main() {

	config.Init()

	tbgot.Init()

}

//func onInlineKeyboardSelect(ctx context.Context, b *tgbot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
//	bot.SendMessage(ctx, &tgbot.SendMessageParams{
//		ChatID: mes.Message.Chat.ID,
//		Text:   "You selected: " + string(data),
//	})
//}
