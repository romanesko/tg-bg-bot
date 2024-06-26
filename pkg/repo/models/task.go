package database

import (
	"bodygraph-bot/pkg/common"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type Task struct {
	models.BaseModel
	TgChatId      int                 `db:"tg_chat_id" json:"tg_chat_id"`
	Request       string              `db:"request" json:"request"`
	Params        map[string]any      `db:"params" json:"params"`
	Response      *common.MessageData `db:"response" json:"response"`
	ResponseReady bool                `db:"response_ready" json:"response_ready"`
	SentToUser    bool                `db:"sent_to_user" json:"sent_to_user"`
	TgMessageId   int                 `db:"tg_message_id" json:"tg_message_id"`
}

func CastTaskFromDbModel(dao *daos.Dao, record *models.Record) *Task {
	task := &Task{}
	task.Id = record.Id
	task.TgChatId = record.GetInt("tg_chat_id")
	task.Request = record.GetString("request")
	_ = record.UnmarshalJSONField("params", &task.Params)
	task.TgMessageId = record.GetInt("tg_message_id")
	task.ResponseReady = record.GetBool("response_ready")
	task.SentToUser = record.GetBool("sent_to_user")
	response := &common.MessageData{}
	_ = record.UnmarshalJSONField("response", &response)
	task.Response = response
	return task

}
