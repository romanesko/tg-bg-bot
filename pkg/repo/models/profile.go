package database

import (
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

/*
type DefaultProfile struct {
		ID          *int64     `json:"id"`
	Name        *string    `json:"name"`
	BirthDate   *ShortDate `json:"birth_date"`
	BirthTime   *ShortTime `json:"birth_time"`
	CountryID   *int       `json:"country_id"`
	CountryName *string    `json:"country_name"`
	CityID      *int       `json:"city_id"`
	CityName    *string    `json:"city_name"`
	TimeUnknown bool       `json:"time_unknown"`
	CityUnknown bool       `json:"city_unknown"`
}
*/

// ensures that the Profile struct satisfy the models.Model interface

type Profile struct {
	models.BaseModel
	Name         string
	Birthday     types.DateTime
	City         *City
	BirthDateSet bool
	BirthTimeSet bool
}

func CastProfileFromDbModel(dao *daos.Dao, record *models.Record) *Profile {
	profile := &Profile{}
	profile.Id = record.Id
	profile.Name = record.GetString("name")
	profile.Birthday = record.GetDateTime("birthday")
	profile.BirthDateSet = record.GetBool("birth_date_set")
	profile.BirthTimeSet = record.GetBool("birth_time_set")

	if errs := dao.ExpandRecord(record, []string{"city"}, nil); len(errs) > 0 {
		println("errs", errs)
	}

	if cityRecord := record.ExpandedOne("city"); cityRecord != nil {
		profile.City = CastCityFromDbModel(dao, cityRecord)
	}

	return profile
}

func ContextToJson(user *User, context interface{}) map[string]any {
	json := make(map[string]any)
	json["tg_chat_id"] = user.TgChatId
	json["context"] = context
	return json
}
