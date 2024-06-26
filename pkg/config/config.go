package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Bot struct {
		Token string `json:"token"`
	} `json:"bot"`
	Api struct {
		ConfigUrl string `json:"config_url"`
		Token     string `json:"token"`
		StartUrl  string `json:"start"`
	} `json:"api"`
	Settings struct {
		ButtonsHeader string `json:"buttons_header"`
	}
}

var conf *Config

func Init() {
	// open config file config.json
	jsonFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	// Read the file contents
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	conf = &Config{}
	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

}

func GetConfig() *Config {

	if conf == nil {
		Init()
	}

	return conf
}
