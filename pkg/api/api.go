package api

import (
	"bodygraph-bot/pkg/common"
	"bodygraph-bot/pkg/config"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
	"strings"
)

type Request struct {
	TgChatId int64                  `json:"tg_chat_id"`
	Text     string                 `json:"text"`
	Context  map[string]interface{} `json:"context"`
}

type Response struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Response common.MessageData
}

type ConfResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Response common.ConfMessageData
}

var cfg = config.GetConfig()

func RefreshConfig() {

	conf := config.GetConfig()

	res, err := FetchConfUrl(conf.Api.ConfigUrl, nil)
	if err != nil {
		panic(err)
	}

	if res.StartUrl == nil {
		panic("Не удалось найти ссылку на главное меню")
	}

	conf.Api.StartUrl = *res.StartUrl
	conf.Settings.ButtonsHeader = res.ButtonsHeader
}

func FetchUrl(url string, params interface{}) (*Response, error) {

	bodyString, err := fetchUrl(url, params)
	if err != nil {
		return nil, err
	}

	result := Response{}
	err = json.Unmarshal(bodyString, &result)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %s", bodyString)
	}
	return &result, nil
}

func FetchConfUrl(url string, params interface{}) (*common.ConfMessageData, error) {
	bodyString, err := fetchUrl(url, params)
	if err != nil {
		return nil, err
	}

	result := ConfResponse{}
	err = json.Unmarshal(bodyString, &result)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %s", bodyString)
	}
	return &result.Response, nil
}

func fetchUrl(url string, params interface{}) ([]byte, error) {

	jsonData, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(url, "http") {
		conf := config.GetConfig()
		url = fmt.Sprintf("%s%s", conf.Api.BaseUrl, url)
	}

	if strings.Contains(url, "?") {
		url = url + "&dkey=" + cfg.Api.Token
	} else {
		url = url + "?dkey=" + cfg.Api.Token
	}

	color.Set(color.FgGreen)
	log.Println("API", url, string(jsonData))
	color.Unset()

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	//log.Println("API", url, string(jsonData))

	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil

}
