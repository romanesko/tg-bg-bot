package botlogic

import (
	"bodygraph-bot/pkg/common"
	"bodygraph-bot/pkg/config"
	"context"
	"encoding/json"
	"fmt"
	tglib "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
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

func BotIsRunning() bool {
	return BotInstance != nil
}

func sendImage(params tglib.SendPhotoParams) error {
	common.FuncLog("sentImage", params.ChatID, params.Photo)
	ctx := context.Background()
	_, err := BotInstance.SendPhoto(ctx, &params)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func sendMessage(params tglib.SendMessageParams) error {
	common.FuncLog("sendMessage", params.ChatID, params.Text)
	ctx := context.Background()
	_, err := BotInstance.SendMessage(ctx, &params)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func SendErrorMessage(chatID int64, err error) error {
	if err == nil {
		return nil
	}
	common.ErrorLog(err)
	return sendMessage(tglib.SendMessageParams{ChatID: chatID, Text: fmt.Sprintf("Error: %s", err.Error())})

}

func SendRawMessage(chatID int64, message string) (*models.Message, error) {
	common.FuncLog("SendRawMessage", chatID, message)
	ctx := context.Background()
	msgInfo, err := BotInstance.SendMessage(ctx, &tglib.SendMessageParams{ChatID: chatID, Text: message})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return msgInfo, nil
}

func SendMessage(chatID int64, msg string, callback *common.MessageData) error {
	common.FuncLog("SendMessage", chatID, msg)
	common.SetUserContext(chatID, callback)
	return sendMessage(tglib.SendMessageParams{ChatID: chatID, Text: msg})
}

func SendHtmlMessageWithPictures(chatID int64, msg string, images []string) error {
	common.FuncLog("SendHtmlMessageWithPictures", chatID, msg)
	//common.SetUserContext(chatID, nil)

	if len(images) == 0 {
		return sendMessage(tglib.SendMessageParams{ChatID: chatID, Text: msg, ParseMode: models.ParseModeHTML})
	}

	for i := range images {
		photo := &models.InputFileString{Data: images[i]}
		if i == len(images)-1 {

			return sendImage(tglib.SendPhotoParams{ChatID: chatID, Caption: msg, ParseMode: models.ParseModeHTML, Photo: photo})
		} else {
			err := sendImage(tglib.SendPhotoParams{ChatID: chatID, Photo: photo})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DeleteMessage(chatID int64, messageID int) error {
	common.FuncLog("DeleteMessage", chatID, messageID)
	ctx := context.Background()
	_, err := BotInstance.DeleteMessage(ctx, &tglib.DeleteMessageParams{ChatID: chatID, MessageID: messageID})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func SendHtmlMessageMessageWithReplyMarkup(chatID int64, msg string, buttons [][]MessageButton, callback *common.MessageData) error {
	common.SetUserContext(chatID, callback)

	rm := inline.New(BotInstance)
	for _, row := range buttons {
		rm = rm.Row()
		for _, button := range row {
			rm.Button(button.Label, []byte(button.Data), OnInlineKeyboardSelect)
		}
	}
	return sendMessage(tglib.SendMessageParams{ChatID: chatID, Text: msg, ReplyMarkup: rm, ParseMode: models.ParseModeHTML})
}

func startHandler(_ context.Context, _ *tglib.Bot, update *models.Update) {
	common.FuncLog("startHandler")
	err := SendStartMessage(update.Message.Chat.ID)
	if err != nil {
		_ = SendErrorMessage(update.Message.Chat.ID, err)
	}
}

func defaultHandler(_ context.Context, _ *tglib.Bot, update *models.Update) {

	if update.Message == nil {
		log.Println("update.Message is nil")
		return
	}

	chatID := update.Message.Chat.ID

	common.FuncLog("defaultHandler", update.Message.Text)

	if strings.HasPrefix(update.Message.Text, "/start") {
		err := SendStartMessage(chatID)
		if err != nil {
			_ = SendErrorMessage(chatID, err)
			return
		}
		return

	}
	userContext := common.GetUserContext(chatID)

	if userContext == nil {
		_ = SendMessage(chatID, "Не понимаю о чём это вы (", nil)
		_ = SendStartMessage(chatID)
		return
	}

	if userContext.Callback == nil {
		_ = SendMessage(chatID, "Выберите пункт из меню, чтобы продолжить, или /start для выхода в главное меню", nil)
		return
	}

	err := processOuterText(chatID, update.Message.Text, *userContext.Callback)
	if err != nil {
		_ = SendErrorMessage(chatID, err)
	}

}

func OnInlineKeyboardSelect(_ context.Context, _ *tglib.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	common.FuncLog("OnInlineKeyboardSelect")

	chatId := mes.Message.Chat.ID
	userContext := common.GetUserContext(chatId)

	if userContext == nil {
		_ = SendMessage(chatId, "Не понимаю о чём это вы (", nil)
		_ = SendStartMessage(chatId)
		return
	}

	var selectedButton = common.MenuItem{}

	err := json.Unmarshal(data, &selectedButton)
	if err != nil {
		_ = SendErrorMessage(mes.Message.Chat.ID, err)
		return
	}

	for i := range userContext.Buttons {
		menuItem := userContext.Buttons[i]
		if menuItem.Label == selectedButton.Label {
			log.Println("Pressed button", menuItem)
			err = processOuter(mes.Message.Chat.ID, &menuItem)
			if err != nil {
				_ = SendErrorMessage(mes.Message.Chat.ID, err)
				return
			}
			return
		}
	}

	_ = SendErrorMessage(mes.Message.Chat.ID, fmt.Errorf("не найдена кнопка из контекста"))

}

func SendMessageData(chatID int64, data common.MessageData) error {
	common.FuncLog("SendMessageData", chatID, data)
	var err error

	for idx, msg := range data.Messages {
		fmt.Printf("idx: %d, len: %d\n", idx, len(data.Messages))

		err = SendHtmlMessageWithPictures(chatID, msg.Text, msg.Images)
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}

	buttons, err := makeButtonsFromMenuItems(data.Buttons)
	if err != nil {
		return err
	}

	cfg := config.GetConfig()

	if len(data.Buttons) > 0 {
		err = SendHtmlMessageMessageWithReplyMarkup(chatID, cfg.Settings.ButtonsHeader, buttons, &data)
		if err != nil {
			return err
		}
	}

	return nil
}
