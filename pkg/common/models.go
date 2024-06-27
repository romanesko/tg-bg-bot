package common

import "encoding/json"

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
	Files         *[]struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"files"`
	Callback *TextCallback `json:"callback"`
}

type ConfMessageData struct {
	StartUrl      *string `json:"start"`
	ButtonsHeader string  `json:"buttons_header"`
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
