package config

import (
	"bodygraph-bot/pkg/common"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

type LocalConfig struct {
	Bot struct {
		Token string `yaml:"token"`
	} `yaml:"bot"`
	Api struct {
		ConfigUrl string `yaml:"config_url"`
		Token     string `yaml:"token"`
	}
}

var config *common.Config

func Init() {

	var localConf *LocalConfig

	file, err := os.OpenFile("config.yaml", os.O_RDONLY, 0600)

	if err != nil {
		log.Fatalf("error opening/creating file: %v", err)
	}
	defer file.Close()

	dec := yaml.NewDecoder(file)
	err = dec.Decode(&localConf)
	if err != nil {
		panic(err)
	}

	if localConf.Api.ConfigUrl == "" {
		log.Fatalf("config.yaml is missing api.config_url")
	}

	if localConf.Api.Token == "" {
		log.Fatalf("config.yaml is missing api.token")
	}

	err = fetchRemoteConfig(localConf.Api.ConfigUrl, localConf.Api.Token, localConf.Bot.Token)
	if err != nil {
		log.Fatalf("error fetching remote config: %v", err)
	}
}

func GetConfig() *common.Config {
	if config == nil {
		Init()
	}
	return config
}

func fetchUrl(url string, token string) ([]byte, error) {

	url = url + "?dkey=" + token

	color.Set(color.FgGreen)
	log.Println("API", url)
	color.Unset()

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) { Body.Close() }(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil

}

func fetchRemoteConfig(configUrl string, configToken string, botToken string) error {
	bodyString, err := fetchUrl(configUrl, configToken)
	if err != nil {
		return err
	}
	result := common.ConfResponse{}
	err = json.Unmarshal(bodyString, &result)
	if err != nil {
		var generic interface{}
		_ = json.Unmarshal(bodyString, &generic)
		return fmt.Errorf("error decoding response: %s", generic)
	}

	conf := common.Config{}
	conf.QueueInterval = result.Response.QueueInterval
	conf.ActionsInterval = result.Response.ActionsInterval

	conf.QueueUrl = result.Response.QueueUrl
	conf.ActionsUrl = result.Response.ActionsUrl

	conf.ButtonsHeader = result.Response.ButtonsHeader
	conf.StartUrl = result.Response.StartUrl
	conf.ApiToken = configToken
	conf.BotToken = botToken

	re := regexp.MustCompile(`(https?://[^/]+)`)
	match := re.FindStringSubmatch(configUrl)
	if len(match) > 1 {
		schemeAndHost := match[1]
		conf.BaseUrl = schemeAndHost
		log.Println("Base Url Host:", schemeAndHost)
	} else {
		return fmt.Errorf("Invalid configUrl: %s", configUrl)
	}

	config = &conf
	return nil
}
