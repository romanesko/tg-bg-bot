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
	"net/url"
	"strings"
)

type Request struct {
	TgChatId int64                  `json:"tg_chat_id"`
	UserInfo common.UserInfo        `json:"user_info"`
	Text     string                 `json:"text"`
	Context  map[string]interface{} `json:"context"`
}

type Response struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Response common.MessageData
}

func FetchUrl(url string, params interface{}) (*Response, error) {

	bodyString, err := fetchUrl(url, params)
	if err != nil {
		return nil, err
	}

	result := Response{}
	err = json.Unmarshal(bodyString, &result)
	if err != nil {
		var generic interface{}
		err = json.Unmarshal(bodyString, &generic)
		if err != nil {
			return nil, fmt.Errorf("NOT A JSON: %s", bodyString)
		}
		return nil, fmt.Errorf("response doesn't match required model: %s", generic)
	}
	return &result, nil
}

func paramsToString(params map[string]interface{}) string {
	urlValues := url.Values{}
	for key, value := range params {
		valueStr := fmt.Sprintf("%v", value)
		urlValues.Add(key, valueStr)
	}
	return urlValues.Encode()
}

func FetchUrlAbstract(url string, params interface{}) ([]byte, error) {
	return fetchUrl(url, params)
}

func GetUrl(url string, params map[string]interface{}) ([]byte, error) {
	var cfg = config.GetConfig()
	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("%s%s", cfg.BaseUrl, url)
	}

	params["dkey"] = cfg.ApiToken

	url = url + "?" + paramsToString(params)

	color.Set(color.FgGreen)
	log.Println("API", url)
	color.Unset()

	resp, err := http.Get(url)

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

func fetchUrl(url string, params interface{}) ([]byte, error) {
	var cfg = config.GetConfig()
	jsonData, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("%s%s", cfg.BaseUrl, url)
	}

	if strings.Contains(url, "?") {
		url = url + "&dkey=" + cfg.ApiToken
	} else {
		url = url + "?dkey=" + cfg.ApiToken
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

	jsonData, textPart, err := extractJSON(bodyBytes)

	if err != nil {

		return nil, err
	}

	if textPart != nil {
		color.Set(color.FgYellow)
		log.Println(string(textPart))
		color.Unset()
	}

	return jsonData, nil

}

func extractJSON(data []byte) (jsonData []byte, textPart []byte, err error) {

	strData := string(data)

	if strings.HasPrefix(strings.TrimSpace(strData), "{") {
		jsonData = data
		return
	}

	idx := strings.Index(strData, "{")

	if idx == -1 {
		err = fmt.Errorf("invalid json: %s", strData)
		return
	}

	textPart = []byte(strData[:idx]) // everything before needle
	jsonData = []byte(strData[idx:])

	return

}
