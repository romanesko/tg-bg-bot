package common

import (
	"encoding/json"
)

type Message struct {
	Text        string   `json:"text"`
	Images      []string `json:"images"`
	ShowPreview bool     `json:"show_preview"`
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

type Action struct {
	Url      string `json:"url"`
	Interval int    `json:"interval"`
}

type ConfResponseResponse struct {
	StartUrl              string    `json:"start"`
	ButtonsHeader         string    `json:"buttons_header"`
	QueueUrl              string    `json:"queue_url"`
	QueueInterval         int       `json:"queue_interval"`
	ActionsUrl            string    `json:"actions_url"`
	Actions               []Action  `json:"actions"`
	ActionsInterval       int       `json:"actions_interval"`
	Commands              []Command `json:"commands"`
	ContextTTL            *uint32   `json:"context_ttl_days"`
	ContextMissingMessage string    `json:"no_context_message"`
	AdminPassword         string    `json:"admin_password"`
}

type Command struct {
	Command     string `json:"command"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

type Config struct {
	BotToken              string
	ApiToken              string
	BaseUrl               string
	StartUrl              string
	ButtonsHeader         string
	QueueUrl              string
	QueueInterval         int
	ActionsUrl            string
	Actions               []Action
	ActionsInterval       int
	Commands              []Command
	ContextTTL            uint32
	ContextMissingMessage string
	AdminPassword         string
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
	KeepOnClick     bool                   `json:"keep_on_click,omitempty"`
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
	State     string `json:"state"`
}

type ActionsCheckUserAvailable struct {
	TgChatID  int64 `json:"tg_chat_id"`
	Available bool  `json:"available"`
}

type ActionsDTO struct {
	CheckUserInChannel []ActionsCheckUserInChannel `json:"check_user_in_channel"`
	CheckUserAvailable []ActionsCheckUserAvailable `json:"check_user_available"`
	Mock               interface{}                 `json:"mock"`
}

type ActionsResponse struct {
	Status   string     `json:"status"`
	Response ActionsDTO `json:"response"`
}
