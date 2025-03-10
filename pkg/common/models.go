package common

import (
	"encoding/json"
)

type Message struct {
	Text   string   `json:"text"`
	Images []string `json:"images"`
}

type TextCallback struct {
	Url     string                 `json:"url"`
	Context map[string]interface{} `json:"context"`
}

type MessageData struct {
	Messages      []Message  `json:"messages"`
	ButtonsHeader string     `json:"buttons_header"`
	Buttons       []MenuItem `json:"buttons"`
	ButtonsRows   *[]int     `json:"buttons_rows"`
	Files         *[]struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"files"`
	Callback *TextCallback `json:"callback"`
}

type ConfResponse struct {
	ID       string               `json:"id"`
	Status   string               `json:"status"`
	Response ConfResponseResponse `json:"response"`
}

type ConfResponseResponse struct {
	StartUrl        string `json:"start"`
	ButtonsHeader   string `json:"buttons_header"`
	QueueUrl        string `json:"queue_url"`
	QueueInterval   int    `json:"queue_interval"`
	ActionsUrl      string `json:"actions_url"`
	ActionsInterval int    `json:"actions_interval"`
}

type Config struct {
	BotToken        string
	ApiToken        string
	BaseUrl         string
	StartUrl        string
	ButtonsHeader   string
	QueueUrl        string
	QueueInterval   int
	ActionsUrl      string
	ActionsInterval int
}

type Menu struct {
	Title string     `json:"title"`
	Items []MenuItem `json:"items"`
}

type MenuItem struct {
	Label   string `json:"label"`
	Command string `json:"command,omitempty"`
	//Action          string `json:"action,omitempty"`
	URL             string                 `json:"url,omitempty"`
	Link            string                 `json:"link,omitempty"`
	ProfileRequired bool                   `json:"ask_profile,omitempty"`
	Context         map[string]interface{} `json:"context,omitempty"`
}

func (m *MenuItem) ToJson() (string, error) {
	res, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (m *MenuItem) FromJson(item string) error {
	err := json.Unmarshal([]byte(item), &m)
	return err
}

type QueueResponse struct {
	Status   string `json:"status"`
	Descr    string `json:"descr"`
	Response *struct {
		Items *[]QueueItem `json:"items"`
	} `json:"response"`
}

type QueueItem struct {
	MessageID int64       `json:"message_id"`
	Datetime  string      `json:"datetime"`
	TgChatID  int64       `json:"tg_chat_id"`
	Data      MessageData `json:"data"`
}

type ActionsCheckUserInChannel struct {
	TgChatID  int64  `json:"tg_chat_id"`
	TgChannel string `json:"tg_channel_name"`
	Exists    bool   `json:"exists"`
}

type ActionsDTO struct {
	CheckUserInChannel []ActionsCheckUserInChannel `json:"check_user_in_channel"`
}

type ActionsResponse struct {
	Status   string     `json:"status"`
	Response ActionsDTO `json:"response"`
}
