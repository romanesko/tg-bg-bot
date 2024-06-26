package repo

import (
	"bodygraph-bot/pkg/api"
	database "bodygraph-bot/pkg/repo/models"
	"fmt"
	"github.com/go-telegram/bot/models"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	pbmodels "github.com/pocketbase/pocketbase/models"
	"github.com/twinj/uuid"
	"log"
	"strings"
)

type RepoType struct {
	app *pocketbase.PocketBase
}

var repo *RepoType

func Init(a *pocketbase.PocketBase) {
	repo = &RepoType{app: a}
}

func InsertCitiesFromApi() error {

	cities, err := api.GetCities()
	if err != nil {
		return fmt.Errorf("InsertCitiesFromApi: %w", err)
	}

	dao := repo.app.Dao()

	collection, err := dao.FindCollectionByNameOrId("cities")
	if err != nil {
		log.Fatal(err)
	}

	for _, city := range cities {

		record, err := dao.FindFirstRecordByData("cities", "city_id", city.CityID)
		if err != nil {
			record = pbmodels.NewRecord(collection)
		}

		record.Load(map[string]any{
			"city_id":             city.CityID,
			"city_name":           city.CityName,
			"country_id":          city.CountryID,
			"country_name":        city.CountryName,
			"city_name_lowercase": strings.ToLower(city.CityName),
		})

		if err := dao.SaveRecord(record); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

/*
birthday date,
		birthtime time,
		time_unknown boolean not null DEFAULT false,
		country_id INTEGER,
		country_name INTEGER,
		city_id INTEGER,
		city_name TEXT,
*/

func findUserByChatId(dao *daos.Dao, id int64) (*database.User, error) {

	record, err := dao.FindFirstRecordByData("users", "tg_chat_id", id)

	if err != nil {
		return nil, err
	}

	return database.CastUserFromDbModel(dao, record), nil
}

func CreateNewUser(chat *models.Chat) (*database.User, error) {

	userCollection, err := repo.app.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		return nil, err
	}
	profileCollection, err := repo.app.Dao().FindCollectionByNameOrId("profiles")
	if err != nil {
		return nil, err
	}

	defaultProfileRecord := pbmodels.NewRecord(profileCollection)
	defaultProfileRecord.Set("name", chat.FirstName)
	defaultProfileRecord.Set("birth_day_set", false)
	defaultProfileRecord.Set("birth_time_set", false)
	if err := repo.app.Dao().SaveRecord(defaultProfileRecord); err != nil {
		return nil, err
	}

	profilesList := []string{defaultProfileRecord.Id}

	userRecord := pbmodels.NewRecord(userCollection)
	userRecord.Set("tg_chat_id", chat.ID)

	var userName = chat.Username
	if userName == "" {
		userName = fmt.Sprintf("chat%d", chat.ID)
	}

	userRecord.Set("name", chat.FirstName)
	userRecord.Set("default_profile", defaultProfileRecord.Id)
	userRecord.Set("profiles", profilesList)

	userRecord.SetUsername(userName)
	u := uuid.NewV4()
	userRecord.SetPassword(u.String())

	if err := repo.app.Dao().SaveRecord(userRecord); err != nil {
		return nil, err
	}

	return database.CastUserFromDbModel(repo.app.Dao(), userRecord), nil
}

func UpdateUserSkill(user database.User, skill string) error {
	record, err := repo.app.Dao().FindRecordById("users", user.Id)
	if err != nil {
		return err
	}
	form := forms.NewRecordUpsert(repo.app, record)
	_ = form.LoadData(map[string]any{
		"skill": skill,
	})
	if err := form.Submit(); err != nil {
		return err
	}
	return nil
}

func GetUserByChatId(chatId int64) (database.User, error) {
	user, err := findUserByChatId(repo.app.Dao(), chatId)
	if err != nil {
		return database.User{}, err
	}
	return *user, err
}

func GetUser(chat models.Chat) (database.User, error) {
	user, err := findUserByChatId(repo.app.Dao(), chat.ID)

	if err != nil {
		user, err = CreateNewUser(&chat)
		if err != nil {
			return database.User{}, err
		}
	}
	//println("err", err)
	//println("record", record)

	return *user, nil
}

func FindCity(name string) ([]database.City, error) {

	cities := []database.City{}

	log.Printf("Looking for city %s\n", name)

	repo.app.Dao().DB().
		Select("id", "city_id", "city_name", "country_id", "country_name").
		From("cities").
		AndWhere(dbx.Like("city_name_lowercase", strings.TrimSpace(strings.ToLower(name)))).
		Limit(100).
		OrderBy("created ASC").
		All(&cities)

	return cities, nil
}

//func AddProfile(tgUserId int64, info UserInfo) error {
//
//	return nil
//}
//
//func DeleteProfile(tgUserId int64, ProfileId int64) error {
//
//	return nil
//}

func UpdateProfile(profile database.Profile) error {
	record, err := repo.app.Dao().FindRecordById("profiles", profile.Id)
	if err != nil {
		return err
	}

	form := forms.NewRecordUpsert(repo.app, record)

	if profile.City == nil {
		profile.City = &database.City{}
	}

	// or form.LoadRequest(r, "")
	form.LoadData(map[string]any{
		"name":           profile.Name,
		"birthday":       profile.Birthday,
		"birth_date_set": profile.BirthDateSet,
		"birth_time_set": profile.BirthTimeSet,
		"city":           profile.City.Id,
	})

	// validate and submit (internally it calls app.Dao().SaveRecord(record) in a transaction)
	if err := form.Submit(); err != nil {
		return err
	}

	return nil
}

func GetTasksToSend() []database.Task {
	tasks := []database.Task{}

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
		tasks = append(tasks, *database.CastTaskFromDbModel(dao, record))
	}

	return tasks
}

func GetTasksToProcess() []database.Task {
	tasks := []database.Task{}

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
		tasks = append(tasks, *database.CastTaskFromDbModel(dao, record))
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
	form.LoadData(map[string]any{
		"tg_chat_id":    tgChatId,
		"request":       request,
		"params":        params,
		"tg_message_id": sentMessageId,
	})

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
	form.LoadData(map[string]any{
		"tg_chat_id":     task.TgChatId,
		"request":        task.Request,
		"params":         task.Params,
		"tg_message_id":  task.TgMessageId,
		"response":       task.Response,
		"response_ready": task.ResponseReady,
		"sent_to_user":   task.SentToUser,
	})

	//log.Println(form.Data())

	// validate and submit (internally it calls app.Dao().SaveRecord(record) in a transaction)
	if err := form.Submit(); err != nil {
		return err
	}
	return nil
}

func RepoIsRunning() bool {
	return repo != nil && repo.app != nil && repo.app.Dao() != nil
}
