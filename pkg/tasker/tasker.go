package tasker

import (
	"bodygraph-bot/pkg/api"
	botlogic "bodygraph-bot/pkg/bot"
	"bodygraph-bot/pkg/repo"
	"log"
	"time"
)

func CheckTasksToProcess() {
	for {
		time.Sleep(time.Second * 5)

		if !repo.RepoIsRunning() {
			log.Println("Repo is not running, skipping CheckTasksToProcess")
			continue
		}

		tasks := repo.GetTasksToProcess()

		if len(tasks) == 0 {
			continue
		}

		log.Println("Found", len(tasks), "tasks to process")

		for _, task := range tasks {
			log.Println("Checking task to process", "chatId", task.TgChatId, "created", task.Created)

			fetchedData, err := api.FetchUrl(task.Request, task.Params)
			if err != nil {
				log.Println("Error fetching data", err)
				continue
			}
			if fetchedData.Status == "complete" {
				log.Println("Task completed successfully")
				task.Response = &fetchedData.Response
				task.ResponseReady = true
				err = repo.UpdateTask(task)
				if err != nil {
					log.Println("Error updating task", err)
				}
			}

		}

	}
}

func CheckTasksToSend() {
	for {
		time.Sleep(time.Second * 5)

		if !repo.RepoIsRunning() {
			log.Println("Repo is not running, skipping CheckTasksToSend ")
			continue
		}

		if !botlogic.BotIsRunning() {
			log.Println("Bot is not running, skipping CheckTasksToSend")
			continue
		}

		tasks := repo.GetTasksToSend()

		if len(tasks) == 0 {
			continue
		}

		log.Println("Found", len(tasks), "tasks to send")

		for _, task := range tasks {
			log.Println("Checking task to send", "chatId", task.TgChatId, "created", task.Created)
			if task.Response == nil {
				log.Println("OOPS: Task has no response, but in marked as ready")
				continue
			}

			err := botlogic.SendMessageData(int64(task.TgChatId), *task.Response)
			if err != nil {
				log.Println("Error sending message", err)
				continue
			}

			task.SentToUser = true
			err = repo.UpdateTask(task)
			if err != nil {
				log.Println("Error updating task", err)
			}
			_ = botlogic.DeleteMessage(int64(task.TgChatId), task.TgMessageId)

		}

	}
}
