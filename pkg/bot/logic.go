package botlogic

import (
	"bodygraph-bot/pkg/api"
	"bodygraph-bot/pkg/common"
	"bodygraph-bot/pkg/config"
	dbmodels "bodygraph-bot/pkg/repo/models"
	"fmt"
	"github.com/go-telegram/bot/models"
	"strconv"
	"strings"
)

type ContextReceiver func(chatID int64)

func getMenu(user dbmodels.User) error {
	common.FuncLog("getMenu")

	api.RefreshConfig()

	cfg := config.GetConfig()

	var menuItem = common.MenuItem{URL: cfg.Api.StartUrl, Label: "main menu"}
	return processOuter(int64(user.TgChatId), &menuItem)
}

//func getMenuItemByCommand(command string, user dbmodels.User) (common.MenuItem, error) {
//	common.FuncLog("getMenuItemByCommand", command)
//
//	if strings.HasPrefix(command, "/") {
//		command = command[1:]
//	}
//
//	menu, err := getMenu(user)
//	if err != nil {
//		return common.MenuItem{}, err
//	}
//	for _, item := range menu.Items {
//		if item.Command == command {
//			return item, nil
//		}
//	}
//
//	return common.MenuItem{}, nil
//
//}

//func AskSkill(chatID int64) error {
//	common.FuncLog("AskSkill", chatID)
//
//	buttons := [][]MessageButton{
//		[]MessageButton{NewButton("Новичок", "skill:beginner"), NewButton("Профи", "skill:professional")},
//	}
//
//	return SendMessageWithReplyMarkup(chatID, "Для начала скажи, новичок ли ты в Дизайне Человека? Или уже знаком с основными понятиями и более менее ориентируешься?\n\nЭто нужно, чтобы понять, использовать ли мне терминологию ДЧ или\nговорить бытовым языком.\n\nДа, я хоть и робот, но всегда стараюсь угодить каждому ;)", buttons, onSkillReceived)
//}
//
//func onSkillReceived(message *models.Message, callbackData string) error {
//	common.FuncLog("OnBirthdayReceived", message.Chat.ID, callbackData)
//	user, _ := repo.GetUser(message.Chat)
//
//	if callbackData == "skill:professional" {
//		user.Skill = "professional"
//	} else if callbackData == "skill:beginner" {
//		user.Skill = "beginner"
//	} else {
//		return SendMessage(message.Chat.ID, "Пожалуйста, воспользуйтесь кнопкой", onSkillReceived)
//	}
//	err := repo.UpdateUserSkill(user, user.Skill)
//	if err != nil {
//		return err
//	}
//	return SendStartMessage(message.Chat)
//}
//
//func AskBirthDate(chatID int64) error {
//	common.FuncLog("AskBirthDate", chatID)
//	return SendMessage(chatID, "Укажите, пожалуйста, дату вашего рождения и время. Время важно для точности, но не обязательно, если не знаете время, то укажите просто дату", OnBirthdayReceived)
//
//}
//
//func AskBirthPlace(chatID int64) error {
//	common.FuncLog("AskBirthPlace", chatID)
//	return SendMessage(chatID, "Напишите город вашего рождения", OnBirthPlaceReceived)
//}
//
//func OnBirthdayReceived(message *models.Message, callbackData string) error {
//	common.FuncLog("OnBirthdayReceived", message.Chat.ID, callbackData)
//
//	str := strings.Split(message.Text, " ")
//	dateString := str[0]
//
//	date, err := utils.ParseDateFromString(dateString)
//	if err != nil {
//		return SendMessage(message.Chat.ID, "Не удалось разобрать дату. Попробуйте указать в формате 1991-12-31 или 31.12.1991 или 31/12/1991", OnBirthdayReceived)
//	}
//
//	var newTime time.Time
//
//	if len(str) > 1 {
//		newTime, err = utils.ParseTimeFromString(str[1])
//		if err != nil {
//			return SendMessage(message.Chat.ID, "Не удалось разобрать время. Укажите в формате 15:04", OnBirthTimeReceived)
//		}
//	}
//
//	log.Println("date:", date)
//
//	user, _ := repo.GetUser(message.Chat)
//
//	if newTime.IsZero() {
//		log.Println("time is not set")
//	} else {
//		log.Printf("new time is %s\n", newTime)
//		date.Add(newTime.Sub(time.Now()))
//		date = time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), 0, 0, time.UTC)
//		user.DefaultProfile.BirthTimeSet = true
//	}
//
//	d := types.DateTime{}
//	if err := d.Scan(date); err != nil {
//		SendMessage(message.Chat.ID, "Не удалось записать в базу", nil)
//		return SendStartMessage(message.Chat)
//
//	}
//	user.DefaultProfile.Birthday = d
//	user.DefaultProfile.BirthDateSet = true
//
//	if err := repo.UpdateProfile(user.DefaultProfile); err != nil {
//		return SendErrorMessage(message.Chat.ID, err)
//	}
//
//	return AskBirthPlace(message.Chat.ID)
//
//}
//
//func OnBirthTimeReceived(message *models.Message, callbackData string) error {
//	common.FuncLog("OnBirthTimeReceived", message.Chat.ID, callbackData)
//
//	//user, _ := db.GetUser(message.Chat)
//
//	//userProfile, err := db.GetProfileById(user.ProfileId)
//	//if err != nil {
//	//	return SendErrorMessage(message.Chat.ID, err)
//	//}
//	//
//	//if callbackData == "notime" {
//	//	userProfile.TimeUnknown = true
//	//	err := db.UpdateProfile(userProfile)
//	//	if err != nil {
//	//		return SendErrorMessage(message.Chat.ID, err)
//	//	}
//	//	return AskBirthPlace(message.Chat.ID)
//	//}
//	//
//	//newTime, err := utils.ParseTimeFromString(message.Text)
//	//if err != nil {
//	//	return SendMessage(message.Chat.ID, "Не удалось разобрать время. Укажите в формате 15:04", OnBirthTimeReceived)
//	//}
//	//
//	//userProfile.BirthTime = &database.ShortTime{Time: newTime}
//	//err = db.UpdateProfile(userProfile)
//	//if err != nil {
//	//	return SendErrorMessage(message.Chat.ID, err)
//	//
//	//}
//
//	return AskBirthPlace(message.Chat.ID)
//}
//
//func OnBirthPlaceReceived(message *models.Message, callbackData string) error {
//	common.FuncLog("OnBirthPlaceReceived", message.Chat.ID, callbackData)
//	//user, _ := db.GetUser(message.Chat)
//	//log.Println("User: ", user)
//
//	if len(callbackData) > 0 {
//		return onCityConfirmed(message, callbackData)
//	}
//
//	cities, err := repo.FindCity(message.Text)
//	if err != nil {
//		return SendErrorMessage(message.Chat.ID, err)
//
//	}
//
//	if len(cities) == 0 {
//		return SendMessage(message.Chat.ID, "Не удалось найти такой город. Попробуйте указать ближайший крупный город в том же часовом моясе", OnBirthPlaceReceived)
//
//	}
//
//	var buttons [][]MessageButton
//
//	for _, city := range cities {
//		buttons = append(buttons, []MessageButton{NewButton(city.CityName, city.Id)})
//	}
//
//	if len(cities) > 1 {
//		return SendMessageWithReplyMarkup(message.Chat.ID, fmt.Sprintf("Выберите из указанных городов, или попробуйте ввести город ещё раз"), buttons, OnBirthPlaceReceived)
//	} else {
//		return SendMessageWithReplyMarkup(message.Chat.ID, fmt.Sprintf("Подтвердите город кнопкой, если он определён правильно, или попробуйте ввести город ещё раз"), buttons, OnBirthPlaceReceived)
//	}
//
//}
//
//func onCityConfirmed(message *models.Message, callbackData string) error {
//	common.FuncLog("onCityConfirmed", message.Chat.ID, callbackData)
//	user, _ := repo.GetUser(message.Chat)
//	city := &dbmodels.City{}
//	city.Id = callbackData
//
//	user.DefaultProfile.City = city
//	err := repo.UpdateProfile(user.DefaultProfile)
//	if err != nil {
//		return SendErrorMessage(message.Chat.ID, err)
//	}
//
//	//
//	//if callbackData == "" {
//	//	return OnBirthPlaceReceived(message, "")
//	//
//	//}
//	//
//	//cityId, err := strconv.Atoi(callbackData)
//	//if err != nil {
//	//	return SendErrorMessage(message.Chat.ID, err)
//	//}
//	//
//	//city, err := db.GetCity(cityId)
//	//if err != nil {
//	//	return SendErrorMessage(message.Chat.ID, err)
//	//
//	//}
//
//	//userProfile, err := db.GetProfileById(user.ProfileId)
//	//if err != nil {
//	//	return SendErrorMessage(message.Chat.ID, err)
//	//}
//	//
//	//userProfile.CityID = &cityId
//	//userProfile.CityName = &city.Name
//	//userProfile.CountryID = &city.CountryID
//	//userProfile.CountryName = &city.CountryName
//	//
//	//err = db.UpdateProfile(userProfile)
//	//if err != nil {
//	//	return SendErrorMessage(message.Chat.ID, err)
//	//
//	//}
//	//
//	//var fullBirthday = userProfile.BirthDate.Format("02.01.2006 в 15:04")
//	//
//	//if userProfile.TimeUnknown {
//	//	fullBirthday = userProfile.BirthDate.Format("02.01.2006 (время неизвестно)")
//	//}
//	//
//	//SendMessage(message.Chat.ID, fmt.Sprintf("Отлично, темперь зная, что вы родились в %s в городе %s мы сможем точно предоставить данные", fullBirthday, city.DisplayName), nil)
//	return SendStartMessage(message.Chat)
//}
//
//func SendSelectProfile(message *models.Message, callbackKey string) int {
//	//chatID := message.Chat.ID
//	//common.FuncLog("SendSelectProfile", chatID, callbackKey)
//	//profiles, err := db.GetProfiles(chatID)
//	//if err != nil {
//	//	return SendErrorMessage(chatID, err)
//	//}
//	//
//	//if len(profiles) == 1 {
//	//	return onProfileSelected(message, fmt.Sprintf("%s:%d", callbackKey, *profiles[0].ID))
//	//}
//	//
//	//buttons := [][]MessageButton{}
//	//for _, profile := range profiles {
//	//	label := fmt.Sprintf("%s (%s)", *profile.Name, (*profile.BirthDate).Format("2006-01-02"))
//	//	buttons = append(buttons, []MessageButton{NewButton(label, fmt.Sprintf("%s:%d", callbackKey, *profile.ID))})
//	//
//	//}
//	//
//	//return SendMessageWithReplyMarkup(chatID, "Выберите профиль: ", buttons, onProfileSelected)
//	return 0
//}

type CallbackData struct {
	Key   string
	Value int64
}

func splitCallback(callbackData string) CallbackData {
	x := strings.Split(callbackData, ":")
	if len(x) == 1 {
		return CallbackData{Key: x[0], Value: 0}
	}
	if len(x) == 2 {
		val, err := strconv.ParseInt(x[1], 10, 64)
		if err != nil {
			return CallbackData{Key: "", Value: 0}
		}
		return CallbackData{Key: x[0], Value: val}
	}
	return CallbackData{Key: "", Value: 0}
}

func onProfileSelected(message *models.Message, callbackData string) error {
	common.FuncLog("onProfileSelected", message.Chat.ID, callbackData)
	//if callbackData == "" {
	//	return SendStartMessage(message.Chat)
	//}
	//
	//cd := splitCallback(callbackData)
	//
	//switch cd.Key {
	//case "info", "prognoz", "diet":
	//	j, err := db.GetProfileByIdJson(cd.Value)
	//	if err != nil {
	//		return SendErrorMessage(message.Chat.ID, err)
	//	}
	//	text := fmt.Sprintf("Тут вызываем `https://bodygraph.online/api_v1/%s.php POST запросом с параметрами:\n<code>%s</code>\nи транслирую ответ пользователю\n\n<b>главное меню</b>: /start", cd.Key, j)
	//	return SendHtmlMessage(message.Chat.ID, text, nil)
	//}

	return SendStartMessage(message.Chat)

}

func SendProfilesList(chat models.Chat) int {
	common.FuncLog("SendProfilesList", chat.ID)
	//user, _ := db.GetUser(chat)
	//
	//profiles, err := db.GetProfiles(user.TgChatID)
	//if err != nil {
	//	return SendErrorMessage(chat.ID, err)
	//}
	//
	//fmt.Println("profiles", profiles)
	//
	//row1 := []MessageButton{NewButton("➕Добавить", "profiles:add")}
	//
	//if len(profiles) > 1 {
	//	row1 = append(row1, NewButton("🗑️ Удалить", "profiles:del"))
	//}
	//
	//buttons := [][]MessageButton{row1}
	//buttons = append(buttons, []MessageButton{NewButton("⤴️ В главное меню", "mainmenu")})
	//
	//text := "Профили"
	//
	//for _, profile := range profiles {
	//	text += fmt.Sprintf("\n - %s (%s)", *profile.Name, (*profile.BirthDate).Format("2006-01-02"))
	//}
	//
	//return SendMessageWithReplyMarkup(chat.ID, text, buttons, onProfileActionSelected)
	return 0
}

//func onProfileActionSelected(message *models.Message, callbackData string) int {
//	common.FuncLog("onProfilesSelected", message.Chat.ID, callbackData)
//	user, _ := db.GetUser(message.Chat)
//
//	switch callbackData {
//	case "profiles:add":
//		return SendMarkDownMessage(message.Chat.ID, fmt.Sprintf("Напишите имя, дату рождения и город\nПримеры:\n`Людмил Огурченко 2018-01-01 02:20 Москва`\n`Андрей 2019-02-03 Химки`\n`Анастасия 2019-02-03`"), onNewProfile)
//
//	case "profiles:del":
//		var buttons [][]MessageButton
//		profiles, err := db.GetProfiles(user.TgChatID)
//		if err != nil {
//			return SendErrorMessage(message.Chat.ID, err)
//		}
//		for _, profile := range profiles {
//			if *profile.ID != user.ProfileId {
//				label := fmt.Sprintf("%s (%s)", *profile.Name, (*profile.BirthDate).Format("2006-01-02"))
//				buttons = append(buttons, []MessageButton{NewButton(label, fmt.Sprintf("%d", *profile.ID))})
//			}
//		}
//
//		return SendMessageWithReplyMarkup(message.Chat.ID, "Выберите профиль для удаления.\n\nглавное меню: /start", buttons, onProfileDelete)
//	}
//	return SendStartMessage(message.Chat)
//
//}

//func onProfileDelete(message *models.Message, callbackData string) int {
//	common.FuncLog("onProfileDelete", message.Chat.ID, callbackData)
//	//user, _ := db.GetUser(message.Chat)
//
//	profileId, err := strconv.ParseInt(callbackData, 10, 64)
//	if err != nil {
//		return SendErrorMessage(message.Chat.ID, err)
//	}
//
//	err = db.DeleteProfile(message.Chat.ID, profileId)
//	if err != nil {
//		return SendErrorMessage(message.Chat.ID, err)
//	}
//	return SendProfilesList(message.Chat)
//}
//
//func onNewProfile(message *models.Message, callbackData string) int {
//	common.FuncLog("onNewProfile", message.Chat.ID, callbackData)
//	//user, _ := db.GetUser(message.Chat)
//
//	println("пробуем разобрать", message.Text)
//	name, date, t, cityName, err := utils.ParseNameDateTimeCityString(message.Text)
//
//	if err != nil {
//		common.ErrorLog(err)
//		return SendMessage(message.Chat.ID, "Не удалось разобрать данные. Попробуйте ещё раз.\n/start для выхода", onNewProfile)
//	}
//
//	text := "Добавлен профиль " + name + ": дата рождения " + date.Format("2006-01-02")
//
//	if t != nil {
//		text += " время: " + t.Format("15:04")
//	}
//
//	city := &database.City{}
//	if cityName != "" {
//		text += " город: " + cityName
//
//		cities, err := db.FindCity(strings.ToLower(cityName))
//		if err != nil {
//			return SendErrorMessage(message.Chat.ID, err)
//		}
//
//		if len(cities) == 0 {
//			return SendMessage(message.Chat.ID, "Не удалось найти такой город. Попробуйте указать ближайший крупный город в том же часовом моясе", onNewProfile)
//		}
//
//		if len(cities) > 1 {
//			citiesString := ""
//			for _, city := range cities {
//				citiesString += fmt.Sprintf("%s\n", city.Name)
//			}
//
//			return SendMessage(message.Chat.ID, "Найдено несколько вариантов городов, попробуйте создать профиль ещё раз указав город из списка:\n"+citiesString, onNewProfile)
//		}
//		city = &cities[0]
//	}
//
//	userInfo := database.UserInfo{Name: &name, BirthDate: date, BirthTime: t, CityID: &city.ID, CityName: &city.Name, CountryID: &city.CountryID, CountryName: &city.CountryName,
//		TimeUnknown: t == nil, CityUnknown: &city.ID == nil}
//
//	db.AddProfile(message.Chat.ID, userInfo)
//
//	return SendProfilesList(message.Chat)
//
//}

func onSettingsSelected(message *models.Message, callbackData string) error {
	common.FuncLog("onSettingsSelected", message.Chat.ID, callbackData)
	//user, _ := db.GetUser(message.Chat)
	//
	//switch callbackData {
	//case "settings:reset":
	//
	//	userProfile, err := db.GetProfileById(user.ProfileId)
	//	if err != nil {
	//		return SendErrorMessage(message.Chat.ID, err)
	//	}
	//
	//	profileId := userProfile.ID
	//
	//	userProfile = database.UserInfo{ID: profileId}
	//
	//	err = db.UpdateProfile(userProfile)
	//	if err != nil {
	//		return SendErrorMessage(message.Chat.ID, err)
	//
	//	}
	//	SendMessage(message.Chat.ID, "Данные о рождении сброшены", nil)
	//	return SendStartMessage(message.Chat)
	//
	//case "balance":
	//	return SendMarkDownMessage(message.Chat.ID, fmt.Sprintf("Просто что-то ещё"), nil)
	//
	//case "mainmenu":
	//	return SendStartMessage(message.Chat)
	//
	//}

	return fmt.Errorf("Как мы тут оказались? 🤔")

}

//func SendBalanceMenu(chatID int64) error {
//	common.FuncLog("SendBalanceMenu", chatID)
//
//	buttons := ButtonsSet(NewButton("➕ Пополнить", "add"), NewButton("⤴️ В главное меню", "mainmenu"))
//
//	return SendMessageWithReplyMarkup(chatID, "Ваш баланс 1 000 000 рублей", buttons, onBalanceSelected)
//
//}
//
//func onBalanceSelected(message *models.Message, callbackData string) error {
//	common.FuncLog("onBalanceSelected", message.Chat.ID, callbackData)
//	//user, _ := db.GetUser(message.Chat)
//	//log.Println("User: ", user)
//
//	switch callbackData {
//	case "add":
//		return SendMarkDownMessage(message.Chat.ID, fmt.Sprintf("Тут будем пополнять баланс"), nil)
//	case "mainmenu":
//		return SendStartMessage(message.Chat)
//	}
//
//	return fmt.Errorf("Как мы тут оказались? 🤔")
//}
