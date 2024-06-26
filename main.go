package main

import (
	"bodygraph-bot/pkg/api"
	tbgot "bodygraph-bot/pkg/bot"
	"bodygraph-bot/pkg/repo"
	"bodygraph-bot/pkg/tasker"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"log"
	"os"
)

func main() {

	app := pocketbase.NewWithConfig(pocketbase.Config{DefaultDataDir: "./pb_data"})

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		return nil
	})
	repo.Init(app)
	api.RefreshConfig()

	go tbgot.Init()

	go tasker.CheckTasksToProcess()
	go tasker.CheckTasksToSend()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}
