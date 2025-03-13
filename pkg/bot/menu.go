package botlogic

import (
	"bodygraph-bot/pkg/api"
	"bodygraph-bot/pkg/common"
	"bodygraph-bot/pkg/config"
	"encoding/json"
	"fmt"
	"log"
)

func SendStartMessage(chatID int64) error {
	common.FuncLog("SendStartMessage", chatID)
	conf := config.GetConfig()
	var menuItem = common.MenuItem{URL: conf.StartUrl, Label: "main menu"}
	return processOuter(chatID, &menuItem)

}

func makeButtonsFromMenuItems(items []common.MenuItem, buttonsRows *[]int) ([][]common.MenuItem, error) {
	common.FuncLog("makeButtonsFromMenuItems")

	var buttons [][]common.MenuItem

	if buttonsRows == nil || len(*buttonsRows) == 0 {
		for _, button := range items {
			//hashData, err := button.ToJson()
			//if err != nil {
			//	return buttons, err
			//}
			buttons = append(buttons, []common.MenuItem{button})
		}
	} else {
		currentRow := 0
		buttons = make([][]common.MenuItem, 1)
		log.Println("MenuItem", items)
		log.Println("buttonsRows", *buttonsRows)
		for idx, button := range items {
			if common.Contains(*buttonsRows, idx) {
				buttons = append(buttons, make([]common.MenuItem, 0))
				currentRow++
			}
			buttons[currentRow] = append(buttons[currentRow], button)
		}
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

	_, _ = SendRawMessage(tgChatId, fmt.Sprintf("Запрос отправлен, ожидайте ответа"))
	return nil

}

func processOuter(chatId int64, menuItem *common.MenuItem) error {
	common.FuncLog("processOuter", chatId, menuItem.Label)
	log.Println(menuItem)

	userInfo := common.GetUserInfo(chatId)

	js, err := contextToJson(&api.Request{TgChatId: chatId, UserInfo: userInfo, Context: menuItem.Context})
	if err != nil {
		return err
	}

	if err = sendRequest(chatId, menuItem.URL, js); err != nil {
		return &CustomError{ChatID: chatId, URL: menuItem.URL, Params: js, Err: err}
	}
	return nil
}

func processOuterText(chatId int64, text string, callback common.TextCallback) error {
	common.FuncLog("processOuterText", chatId, text)
	js, err := contextToJson(&api.Request{TgChatId: chatId, Text: text, Context: callback.Context})
	if err != nil {
		return err
	}

	if err = sendRequest(chatId, callback.Url, js); err != nil {
		return &CustomError{ChatID: chatId, URL: callback.Url, Params: js, Err: err}
	}
	return nil
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
