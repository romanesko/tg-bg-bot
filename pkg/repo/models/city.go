package database

import (
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

/*
type DefaultCity struct {
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

// ensures that the City struct satisfy the models.Model interface

type City struct {
	models.BaseModel
	CityId      int    `db:"city_id" json:"city_id"`
	CityName    string `db:"city_name" json:"city_name"`
	CountryId   int    `db:"country_id" json:"country_id"`
	CountryName string `db:"country_name" json:"country_name"`
}

func CastCityFromDbModel(dao *daos.Dao, record *models.Record) *City {
	city := &City{}
	city.Id = record.Id
	city.CityId = record.GetInt("city_id")
	city.CityName = record.GetString("city_name")
	city.CountryId = record.GetInt("country_id")
	city.CountryName = record.GetString("country_name")
	return city
}
