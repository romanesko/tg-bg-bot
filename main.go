package main

import (
	botlogic "bodygraph-bot/pkg/bot"
	"bodygraph-bot/pkg/config"
	"bodygraph-bot/pkg/web"
	"github.com/fatih/color"
	"log"
	"os"
)

var Version string

func main() {

	println()
	color.Set(color.FgBlue)
	log.Printf(" Application started, ver. %s ", Version)
	color.Unset()
	println()

	if os.Getenv("TEST_WEBSERVER") != "" {
		go web.Init()
	}

	config.Init()

	botlogic.Init()

	log.Println("STARTED")

	//ch := make(chan bool)
	//<-ch

}
