package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"regexp"
)

type Config struct {
	Bot struct {
		Token string `yaml:"token"`
	} `yaml:"bot"`
	Api struct {
		ConfigUrl string `yaml:"config_url"`
		Token     string `yaml:"token"`
		StartUrl  string `yaml:"start"`
		BaseUrl   string `yaml:"base_url"`
	} `yaml:"api"`
	Settings struct {
		ButtonsHeader string `yaml:"buttons_header"`
	}
}

var conf *Config

func Init() {
	file, err := os.OpenFile("config.yaml", os.O_RDONLY, 0600)

	if err != nil {
		log.Fatalf("error opening/creating file: %v", err)
	}
	defer file.Close()

	dec := yaml.NewDecoder(file)
	err = dec.Decode(&conf)
	if err != nil {
		panic(err)
	}

	if conf.Api.ConfigUrl == "" {
		log.Fatalf("config.yaml is missing api.config_url")
	}

	if conf.Api.Token == "" {
		log.Fatalf("config.yaml is missing api.token")
	}

	if conf.Api.BaseUrl == "" {
		re := regexp.MustCompile(`(https?://[^/]+)`)

		// Find the scheme and host using the regular expression
		match := re.FindStringSubmatch(conf.Api.ConfigUrl)

		if len(match) > 1 {
			// The scheme and host are in the first submatch
			schemeAndHost := match[1]
			conf.Api.BaseUrl = schemeAndHost
			log.Println("Base Url Host:", schemeAndHost)
		}
	}
}

func GetConfig() *Config {

	if conf == nil {
		Init()
	}

	return conf
}
