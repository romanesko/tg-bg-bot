package botlogic

import (
	"bodygraph-bot/pkg/api"
	"bodygraph-bot/pkg/common"
	"bodygraph-bot/pkg/repo"
	database "bodygraph-bot/pkg/repo/models"
	"fmt"
	"github.com/go-telegram/bot/models"
	"log"
)

func sendWelcomeMessage(chatID int64) error {
	common.FuncLog("sendWelcomeMessage")
	return SendMessage(chatID, "Привет. Я бот-бод. Я могу рассчитать бодиграф (карту Дизайна Человека), рассказать многое о тебе, а также поведать, какое влияние прямо сейчас оказывают на тебя небесные тела.", nil)
}

func SendStartMessage(chat models.Chat) error {
	common.FuncLog("SendStartMessage", chat.ID)

	user, err := repo.GetUser(chat)
	if err != nil {
		return err
	}

	err = getMenu(user)
	if err != nil {
		return err
	}

	return nil

}

func makeButtonsFromMenuItems(items []common.MenuItem) ([][]MessageButton, error) {
	common.FuncLog("makeButtonsFromMenuItems")

	var buttons [][]MessageButton
	for _, button := range items {
		hashData, err := button.ToJson()
		if err != nil {
			return buttons, err
		}
		buttons = append(buttons, []MessageButton{NewButton(button.Label, hashData)})
	}
	return buttons, nil
}

func sendRequest(tgChatId int64, url string, params interface{}) error {
	common.FuncLog("sendRequest", tgChatId, url, params)
	resp, err := api.FetchUrl(url, params)
	if err != nil {
		println("Error fetching data", err)
		return err
	}

	common.SetUserContext(tgChatId, &resp.Response)

	if resp.Status == "complete" {
		return SendMessageData(tgChatId, resp.Response)
	}

	sentMessage, err := SendRawMessage(tgChatId, fmt.Sprintf("Запрос отправлен, ожидайте ответа"))
	if err != nil {
		return err
	}
	return repo.AddTask(int(tgChatId), url, params, sentMessage.ID)
}

func handleButtonResponse(message *models.Message, callbackData string) error {
	common.FuncLog("handleButtonResponse", message.Chat.ID, callbackData)
	log.Println(callbackData)

	menuItem := common.MenuItem{}

	err := menuItem.FromJson(callbackData)
	if err != nil {
		return err
	}

	//if menuItem.Action != "" {
	//	return processInner(message.Chat.ID, &menuItem)
	//}

	return processOuter(message.Chat.ID, &menuItem)

}

func processOuter(chatId int64, menuItem *common.MenuItem) error {
	common.FuncLog("processOuter", chatId, menuItem.Label)
	log.Println(menuItem)
	user, err := repo.GetUserByChatId(chatId)
	if err != nil {
		return err
	}

	return sendRequest(int64(user.TgChatId), menuItem.URL, database.ContextToJson(&user, interface{}(nil)))

}

func processOuterText(chatId int64, text string, callback common.TextCallback) error {
	common.FuncLog("processOuterText", chatId, text)
	log.Println(callback)
	user, err := repo.GetUserByChatId(chatId)
	if err != nil {
		return err
	}

	context := map[string]any{
		"text":    text,
		"context": callback.Context,
	}

	return sendRequest(int64(user.TgChatId), callback.Url, database.ContextToJson(&user, context))

}

//func processInner(chatId int64, menuItem *common.MenuItem) error {
//	common.FuncLog("processInner", chatId, menuItem.Label)
//	log.Println(menuItem)
//	SendMarkDownMessage(chatId, "Не реализовано", nil)
//	return nil
//}
