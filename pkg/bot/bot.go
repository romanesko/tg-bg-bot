package bot

import (
	"bodygraph-bot/pkg/config"
	"context"
	tglib "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"log"
	"os"
	"os/signal"
)

var BotInstance *tglib.Bot

type MessageButton struct {
	Label string
	Data  string
}

func NewButton(label string, data string) MessageButton {
	return MessageButton{
		Label: label,
		Data:  data,
	}
}

func ButtonsSet(buttons ...MessageButton) [][]MessageButton {
	var out [][]MessageButton
	for _, button := range buttons {
		out = append(out, []MessageButton{button})
	}
	return out
}

func Init() {
	cfg := config.GetConfig()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	telegramBotToken := cfg.Bot.Token

	opts := []tglib.Option{
		tglib.WithDefaultHandler(defaultHandler),
	}

	var err error
	BotInstance, err = tglib.New(telegramBotToken, opts...)

	if err != nil {
		panic(err)
	}
	println("Bot instance created", BotInstance)

	BotInstance.RegisterHandler(tglib.HandlerTypeMessageText, "/start", tglib.MatchTypeExact, startHandler)

	BotInstance.Start(ctx)
}

func sendMessage(params tglib.SendMessageParams) int {
	log.Println("→ sendMessage")
	ctx := context.Background()
	_, err := BotInstance.SendMessage(ctx, &params)
	if err != nil {
		log.Println(err)
		return 1
	}
	return 0
}
func SendErrorMessage(chatID int64, err error) int {
	log.Println("got err: %v", err)
	return sendMessage(tglib.SendMessageParams{ChatID: chatID, Text: "Что-то пошло не так, мы уже разибраемся с проблемой"})

}

func SendMessage(chatID int64, msg string, callback func(message *models.Message, callbackData string) int) int {
	sendMessage(tglib.SendMessageParams{ChatID: chatID, Text: msg})
	SetUserContext(chatID, callback)
	return 0
}
func SendMarkDownMessage(chatID int64, msg string, callback func(message *models.Message, callbackData string) int) int {
	SetUserContext(chatID, callback)
	return sendMessage(tglib.SendMessageParams{ChatID: chatID, Text: msg, ParseMode: models.ParseModeMarkdown})
}

func SendMessageWithReplyMarkup(chatID int64, msg string, buttons [][]MessageButton, callback func(message *models.Message, callbackData string) int) int {
	SetUserContext(chatID, callback)

	rm := inline.New(BotInstance)
	for _, row := range buttons {
		rm = rm.Row()
		for _, button := range row {
			rm.Button(button.Label, []byte(button.Data), OnInlineKeyboardSelect)
		}
	}
	return sendMessage(tglib.SendMessageParams{ChatID: chatID, Text: msg, ReplyMarkup: rm})
}

func startHandler(ctx context.Context, b *tglib.Bot, update *models.Update) {
	log.Println("→ startHandler")
	SendStartMessage(update.Message.Chat)
}

func defaultHandler(ctx context.Context, b *tglib.Bot, update *models.Update) {
	log.Println("→ defaultHandler")

	if update.Message == nil {
		log.Println("update.Message is nil")
		return
	}

	log.Println("update.Message.Chat.ID", update.Message.Chat.ID)

	userContext := GetUserContext(update.Message.Chat.ID)

	if userContext == nil {
		SendMessage(update.Message.Chat.ID, "Не понимаю о чём это вы (", nil)
		SendStartMessage(update.Message.Chat)
		return
	}

	userContext(update.Message, "")

}

func OnInlineKeyboardSelect(ctx context.Context, b *tglib.Bot, mes models.MaybeInaccessibleMessage, data []byte) {

	log.Println("→ OnInlineKeyboardSelect")

	userContext := GetUserContext(mes.Message.Chat.ID)

	if userContext == nil {
		SendMessage(mes.Message.Chat.ID, "Не понимаю о чём это вы (", nil)
		SendStartMessage(mes.Message.Chat)
		return
	}

	userContext(mes.Message, string(data))
}
