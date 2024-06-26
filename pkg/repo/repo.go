package repo

import (
	database "bodygraph-bot/pkg/repo/models"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
	pbmodels "github.com/pocketbase/pocketbase/models"
	"log"
)

type Type struct {
	app *pocketbase.PocketBase
}

var repo *Type

func Init(a *pocketbase.PocketBase) {
	repo = &Type{app: a}
}

func GetTasksToSend() []database.Task {
	var tasks []database.Task

	dao := repo.app.Dao()

	records, err := dao.FindRecordsByExpr("tasks",
		dbx.HashExp{"response_ready": true, "sent_to_user": false},
	)

	if err != nil {
		log.Println(err)
		return tasks
	}

	if len(records) == 0 {
		return tasks
	}

	for _, record := range records {
		tasks = append(tasks, *database.CastTaskFromDbModel(record))
	}

	return tasks
}

func GetTasksToProcess() []database.Task {
	var tasks []database.Task

	dao := repo.app.Dao()

	records, err := dao.FindRecordsByExpr("tasks",
		dbx.HashExp{"response_ready": false, "sent_to_user": false},
	)

	if err != nil {
		log.Println(err)
		return tasks
	}

	if len(records) == 0 {
		return tasks
	}

	for _, record := range records {
		tasks = append(tasks, *database.CastTaskFromDbModel(record))
	}

	return tasks
}

func AddTask(tgChatId int, request string, params interface{}, sentMessageId int) error {
	collection, err := repo.app.Dao().FindCollectionByNameOrId("tasks")
	if err != nil {
		return err
	}

	record := pbmodels.NewRecord(collection)

	form := forms.NewRecordUpsert(repo.app, record)

	// or form.LoadRequest(r, "")
	err = form.LoadData(map[string]any{
		"tg_chat_id":    tgChatId,
		"request":       request,
		"params":        params,
		"tg_message_id": sentMessageId,
	})
	if err != nil {
		return err
	}

	// validate and submit (internally it calls app.Dao().SaveRecord(record) in a transaction)
	if err := form.Submit(); err != nil {
		return err
	}
	return nil
}

func UpdateTask(task database.Task) error {
	record, err := repo.app.Dao().FindRecordById("tasks", task.Id)
	if err != nil {
		return err
	}

	form := forms.NewRecordUpsert(repo.app, record)

	// or form.LoadRequest(r, "")
	err = form.LoadData(map[string]any{
		"tg_chat_id":     task.TgChatId,
		"request":        task.Request,
		"params":         task.Params,
		"tg_message_id":  task.TgMessageId,
		"response":       task.Response,
		"response_ready": task.ResponseReady,
		"sent_to_user":   task.SentToUser,
	})
	if err != nil {
		return err
	}

	//log.Println(form.Data())

	// validate and submit (internally it calls app.Dao().SaveRecord(record) in a transaction)
	if err := form.Submit(); err != nil {
		return err
	}
	return nil
}

func IsRunning() bool {
	return repo != nil && repo.app != nil && repo.app.Dao() != nil
}
