package botlogic

import (
	//botlogic "bodygraph-bot/pkg/bot"
	"bodygraph-bot/pkg/common"
	"bodygraph-bot/pkg/config"
	"bodygraph-bot/pkg/repo"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	tglib "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

var red = color.New(color.FgRed).SprintFunc()

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
	BotInstance.RegisterHandler(tglib.HandlerTypeMessageText, "/refresh_cities", tglib.MatchTypeExact, citiesHandler)
	//BotInstance.RegisterHandler(tglib.HandlerTypeMessageText, "/refresh_menu", tglib.MatchTypeExact, refreshMenuHandler)

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
	//return sendMessage(tglib.SendMessageParams{ChatID: chatID, Text: "Что-то пошло не так, мы уже разибраемся с проблемой"})

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
	sendMessage(tglib.SendMessageParams{ChatID: chatID, Text: msg})
	common.SetUserContext(chatID, callback)
	return nil
}

//func SendMarkDownMessage(chatID int64, msg string, callback *common.MessageData) error {
//	common.FuncLog("SendMarkDownMessage", chatID, msg)
//	common.SetUserContext(chatID, callback)
//	return sendMessage(tglib.SendMessageParams{ChatID: chatID, Text: msg, ParseMode: models.ParseModeMarkdown})
//}

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

//func SendMessageWithReplyMarkup(chatID int64, msg string, buttons [][]MessageButton, callback *common.TextCallback) error {
//	common.SetUserContext(chatID, callback)
//
//	rm := inline.New(BotInstance)
//	for _, row := range buttons {
//		rm = rm.Row()
//		for _, button := range row {
//			rm.Button(button.Label, []byte(button.Data), OnInlineKeyboardSelect)
//		}
//	}
//	return sendMessage(tglib.SendMessageParams{ChatID: chatID, Text: msg, ReplyMarkup: rm})
//}

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

func startHandler(ctx context.Context, b *tglib.Bot, update *models.Update) {
	common.FuncLog("startHandler")
	err := SendStartMessage(update.Message.Chat)
	SendErrorMessage(update.Message.Chat.ID, err)
}

func citiesHandler(ctx context.Context, b *tglib.Bot, update *models.Update) {
	common.FuncLog("citiesHandler")
	SendMessage(update.Message.Chat.ID, "Обновляю справочник городов", nil)
	err := repo.InsertCitiesFromApi()
	if err != nil {
		common.ErrorLog(err)
		SendMessage(update.Message.Chat.ID, err.Error(), nil)
		return
	}
	SendMessage(update.Message.Chat.ID, "Города загружены успешно", nil)
}

//func refreshMenuHandler(ctx context.Context, b *tglib.Bot, update *models.Update) {
//	common.FuncLog("refreshMenuHandler")
//	menu, err := getMenu()
//	if err != nil {
//		SendErrorMessage(update.Message.Chat.ID, err)
//		return
//	}
//
//	commands := make([]models.BotCommand, 0)
//
//	for _, menuItem := range menu.Items {
//		println("menuItem.Command", menuItem.Command)
//		commands = append(commands, models.BotCommand{Command: menuItem.Command, Description: menuItem.Label})
//
//	}
//	params := tglib.SetMyCommandsParams{
//		Commands: commands,
//		Scope:    &models.BotCommandScopeDefault{},
//	}
//	_, err = BotInstance.SetMyCommands(ctx, &params)
//
//	if err != nil {
//		SendErrorMessage(update.Message.Chat.ID, err)
//		return
//	}
//
//	SendMessage(update.Message.Chat.ID, "Меню обновлено успешно.\n\nМожет потребоваться некоторое время для того, чтобы обновлённый список появился на клиенте\n\nВозможно потребуется выйти и зайти в чат, чтобы увидеть изменения", nil)
//
//}

//func processCommand(chatID int64, cmd string) error {
//	common.FuncLog("processCommand", chatID, cmd)
//
//	log.Printf("Got command %s\n", cmd)
//
//	user, err := repo.GetUserByChatId(chatID)
//	if err != nil {
//		return SendErrorMessage(chatID, err)
//	}
//
//	menuItem, err := getMenuItemByCommand(cmd, user)
//
//	if err != nil {
//		return SendErrorMessage(chatID, err)
//
//	}
//
//	if menuItem.Action != "" {
//		err = processInner(chatID, &menuItem)
//	} else {
//		err = processOuter(chatID, &menuItem)
//	}
//
//	if err != nil {
//		return SendErrorMessage(chatID, err)
//
//	}
//	return nil
//}

func defaultHandler(ctx context.Context, b *tglib.Bot, update *models.Update) {

	if update.Message == nil {
		log.Println("update.Message is nil")
		return
	}

	chatID := update.Message.Chat.ID

	common.FuncLog("defaultHandler", update.Message.Text)

	if strings.HasPrefix(update.Message.Text, "/start") {
		user, err := repo.GetUserByChatId(chatID)
		if err != nil {
			SendErrorMessage(update.Message.Chat.ID, err)
			return
		}
		err = getMenu(user)
		if err != nil {
			SendErrorMessage(chatID, err)
			return
		}
		return

	}
	userContext := common.GetUserContext(chatID)

	if userContext == nil {
		SendMessage(chatID, "Не понимаю о чём это вы (", nil)
		return
	}

	if userContext.Callback == nil {
		SendMessage(chatID, "Выберите пункт из меню, чтобы продолжить, или /start для выхода в главное меню", nil)
		return
	}

	err := processOuterText(chatID, update.Message.Text, *userContext.Callback)
	if err != nil {
		SendErrorMessage(chatID, err)
	}

}

func OnInlineKeyboardSelect(ctx context.Context, b *tglib.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	common.FuncLog("OnInlineKeyboardSelect")

	userContext := common.GetUserContext(mes.Message.Chat.ID)

	//if strings.HasPrefix(string(data), "/") {
	//	processCommand(mes.Message.Chat.ID, string(data))
	//}

	if userContext == nil {
		SendMessage(mes.Message.Chat.ID, "Не понимаю о чём это вы (", nil)
		SendStartMessage(mes.Message.Chat)
		return
	}

	var selectedButton = common.MenuItem{}

	err := json.Unmarshal(data, &selectedButton)
	if err != nil {
		SendErrorMessage(mes.Message.Chat.ID, err)
		return
	}

	for i := range userContext.Buttons {
		menuItem := userContext.Buttons[i]
		if menuItem.Label == selectedButton.Label {
			log.Println("Pressed button", menuItem)
			err = processOuter(mes.Message.Chat.ID, &menuItem)
			if err != nil {
				SendErrorMessage(mes.Message.Chat.ID, err)
				return
			}
			return
		}
	}

	SendErrorMessage(mes.Message.Chat.ID, fmt.Errorf("Не найдена кнопка из контекста"))

	//err := userContext(mes.Message, string(data))
	//if err != nil {
	//	SendErrorMessage(mes.Message.Chat.ID, err)
	//}
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
