package main

import (
	botlogic "bodygraph-bot/pkg/bot"
	"bodygraph-bot/pkg/config"
	"bodygraph-bot/pkg/web"
	"log"
	"os"
)

func main() {

	if os.Getenv("TEST_WEBSERVER") != "" {
		go web.Init()
	}

	config.Init()

	botlogic.Init()

	log.Println("STARTED")

	//ch := make(chan bool)
	//<-ch

}
