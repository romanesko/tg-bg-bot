package database

import (
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

/*
type User struct {
	TgChatID   int64
	TgUsername string
	FirstName  string
	ProfileId  int64
}
*/

// ensures that the User struct satisfy the models.Model interface

type User struct {
	models.BaseModel
	TgChatId       int
	Username       string
	Skill          string
	DefaultProfile Profile
	Profiles       []Profile
}

func CastUserFromDbModel(dao *daos.Dao, record *models.Record) *User {
	user := &User{}
	user.Id = record.Id
	user.Skill = record.GetString("skill")
	user.TgChatId = record.GetInt("tg_chat_id")
	user.Username = record.GetString("username")

	if errs := dao.ExpandRecord(record, []string{"default_profile", "profiles"}, nil); len(errs) > 0 {
		println("errs", errs)
	}

	if defaultProfileRecord := record.ExpandedOne("default_profile"); defaultProfileRecord != nil {
		user.DefaultProfile = *CastProfileFromDbModel(dao, defaultProfileRecord)
	}

	if profilesRecord := record.ExpandedOne("profiles"); profilesRecord != nil {
		user.Profiles = []Profile{*CastProfileFromDbModel(dao, profilesRecord)}
	}

	return user
}
