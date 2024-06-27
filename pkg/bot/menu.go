package botlogic

import (
	"bodygraph-bot/pkg/api"
	"bodygraph-bot/pkg/common"
	"bodygraph-bot/pkg/config"
	"bodygraph-bot/pkg/repo"
	"encoding/json"
	"fmt"
	"log"
)

func SendStartMessage(chatID int64) error {
	common.FuncLog("SendStartMessage", chatID)
	api.RefreshConfig()

	cfg := config.GetConfig()

	var menuItem = common.MenuItem{URL: cfg.Api.StartUrl, Label: "main menu"}
	return processOuter(chatID, &menuItem)

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

func processOuter(chatId int64, menuItem *common.MenuItem) error {
	common.FuncLog("processOuter", chatId, menuItem.Label)
	log.Println(menuItem)

	js, err := contextToJson(&api.Request{TgChatId: chatId, Context: menuItem.Context})
	if err != nil {
		return err
	}
	return sendRequest(chatId, menuItem.URL, js)

}

func processOuterText(chatId int64, text string, callback common.TextCallback) error {
	common.FuncLog("processOuterText", chatId, text)
	js, err := contextToJson(&api.Request{TgChatId: chatId, Text: text, Context: callback.Context})
	if err != nil {
		return err
	}
	return sendRequest(chatId, callback.Url, js)
}

func contextToJson(request *api.Request) (map[string]any, error) {
	var inInterface map[string]interface{}
	inRec, _ := json.Marshal(request)
	err := json.Unmarshal(inRec, &inInterface)
	if err != nil {
		return nil, err
	}
	return inInterface, nil
}
