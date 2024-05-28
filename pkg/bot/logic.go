package bot

import (
	"bodygraph-bot/pkg/database"
	"bodygraph-bot/pkg/utils"
	"fmt"
	"github.com/go-telegram/bot/models"
	"log"
	"strconv"
	"strings"
	"time"
)

var db = database.NewDb()

type ContextReceiver func(chatID int64)

func SendStartMessage(chat models.Chat) int {
	log.Println("→ SendStartMessage")
	user, err := db.GetUser(chat)
	if err != nil {
		log.Println(err)
		return SendErrorMessage(chat.ID, err)

	}

	if user.BirthDate == nil {
		return AskBirthDate(chat.ID)

	}

	if user.CityID == nil {
		return AskBirthPlace(chat.ID)

	}

	SetUserContext(chat.ID, nil)

	var buttons [][]MessageButton

	buttons = ButtonsSet(
		NewButton("ℹ️ Общая информация", "info"),
		NewButton("🫣 Прогноз на день", "prognoz"),
		NewButton("🍔 Пресональная диета", "diet"),
		NewButton("💳 Баланс", "balance"),
	)
	buttons = append(buttons, []MessageButton{NewButton("⚙️Настройки", "settings"), NewButton("❓ Помощь", "help")})

	return SendMessageWithReplyMarkup(chat.ID, "Главное меню", buttons, onMenuSelected)

	//sendMessage(chat.ID, "Всё ок")

}

func AskBirthDate(chatID int64) int {
	log.Println("→ AskBirthDate")
	return SendMessage(chatID, "Укажите, пожалуйста, дату вашего рождения", OnBirthdayReceived)

}

func AskBirthTime(chatID int64) int {
	log.Println("→ AskBirthTime")

	var buttons [][]MessageButton

	buttons = ButtonsSet(NewButton("Мне неизвестно время рождения 😩", "notime"))

	return SendMessageWithReplyMarkup(chatID,
		"Укажите время вашего рождения (это важно для точности). Если вы не знаете, и нет возможности уточнить, но нажмите кнопку ниже",
		buttons, OnBirthTimeReceived)
}

func AskBirthPlace(chatID int64) int {
	log.Println("→ AskBirthPlace")
	return SendMessage(chatID, "Напишите город вашего рождения", OnBirthPlaceReceived)
}

func OnBirthdayReceived(message *models.Message, callbackData string) int {
	log.Println("→ OnBirthdayReceived")
	date, err := utils.ParseDateFromString(message.Text)
	if err != nil {
		return SendMessage(message.Chat.ID, "Не удалось разобрать дату. Попробуйте указать в формате 1991-12-31 или 31.12.1991 или 31/12/1991", OnBirthdayReceived)
	}

	user, _ := db.GetUser(message.Chat)
	user.BirthDate = &date
	err = db.UpdateUser(&user)
	if err != nil {
		return SendErrorMessage(message.Chat.ID, err)
	}

	return AskBirthTime(message.Chat.ID)
}

func OnBirthTimeReceived(message *models.Message, callbackData string) int {
	log.Println("→ OnBirthTimeReceived")

	user, _ := db.GetUser(message.Chat)

	if callbackData == "notime" {
		user.TimeUnknown = true
		err := db.UpdateUser(&user)
		if err != nil {
			return SendErrorMessage(message.Chat.ID, err)
		}
		return AskBirthPlace(message.Chat.ID)
	}

	t, err := utils.ParseTimeFromString(message.Text)
	if err != nil {
		return SendMessage(message.Chat.ID, "Не удалось разобрать время. Укажите в формате 15:04", OnBirthTimeReceived)
	}

	if user.BirthDate == nil {
		log.Println("User has no birth date")
		return 1
	}

	newTime := time.Date(user.BirthDate.Year(), user.BirthDate.Month(), user.BirthDate.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	user.BirthDate = &newTime
	err = db.UpdateUser(&user)
	if err != nil {
		return SendErrorMessage(message.Chat.ID, err)

	}

	return AskBirthPlace(message.Chat.ID)
}

func OnBirthPlaceReceived(message *models.Message, callbackData string) int {
	log.Println("→ OnBirthPlaceReceived")
	//user, _ := db.GetUser(message.Chat)
	//log.Println("User: ", user)
	fmt.Printf("%s", message.Text)

	if len(callbackData) > 0 {
		return onCityConfirmed(message, callbackData)
	}

	cities, err := db.FindCity(strings.TrimSpace(strings.ToLower(message.Text)))
	if err != nil {
		log.Println(err)
		return SendErrorMessage(message.Chat.ID, err)

	}

	if len(cities) == 0 {
		return SendMessage(message.Chat.ID, "Не удалось найти такой город. Попробуйте указать ближайший крупный город в том же часовом моясе", OnBirthPlaceReceived)

	}

	var buttons [][]MessageButton

	for _, city := range cities {
		buttons = append(buttons, []MessageButton{NewButton(city.DisplayName, strconv.Itoa(city.ID))})
	}

	if len(cities) > 1 {
		return SendMessageWithReplyMarkup(message.Chat.ID, fmt.Sprintf("Выберите из указанных городов, или попробуйте ввести город ещё раз"), buttons, OnBirthPlaceReceived)
	} else {
		return SendMessageWithReplyMarkup(message.Chat.ID, fmt.Sprintf("Подтвердите город кнопкой, если он опреедлён правильно, или попробуйте ввести город ещё раз"), buttons, OnBirthPlaceReceived)
	}

}

func onCityConfirmed(message *models.Message, callbackData string) int {
	log.Println("→ onCityConfirmed")
	user, _ := db.GetUser(message.Chat)

	if callbackData == "" {
		return OnBirthPlaceReceived(message, "")

	}

	cityId, err := strconv.Atoi(callbackData)
	if err != nil {
		log.Println(err)
		return SendErrorMessage(message.Chat.ID, err)
	}

	city, err := db.GetCity(cityId)
	if err != nil {
		log.Println(err)
		return SendErrorMessage(message.Chat.ID, err)

	}

	user.CityID = &cityId
	user.CityName = &city.Name
	user.CountryID = &city.CountryID
	user.CountryName = &city.CountryName

	err = db.UpdateUser(&user)
	if err != nil {
		return SendErrorMessage(message.Chat.ID, err)

	}

	var fullBirthday = user.BirthDate.Format("02.01.2006 в 15:04")

	if user.TimeUnknown {
		fullBirthday = user.BirthDate.Format("02.01.2006 (время неизвестно)")
	}

	SendMessage(message.Chat.ID, fmt.Sprintf("Отлично, темперь зная, что вы родились в %s в городе %s мы сможем точно предоставить данные", fullBirthday, city.DisplayName), nil)
	return SendStartMessage(message.Chat)
}

func onMenuSelected(message *models.Message, callbackData string) int {
	log.Println("→ onMenuSelected")
	user, _ := db.GetUser(message.Chat)

	switch callbackData {
	case "info":
		SendMarkDownMessage(message.Chat.ID, fmt.Sprintf("Тут я тупо вызываю `https://bodygraph.online/api_v1/info.php?dkey=<bot_key>&city=%d&bd=%s` и транслирую ответ пользователю", *user.CityID, user.BirthDate.Format("2006-01-02 15:04")), nil)
		break
	case "prognoz":
		SendMarkDownMessage(message.Chat.ID, fmt.Sprintf("И здесь будет просто вызов `https://bodygraph.online/api_v1/prognoz.php?dkey=<bot_key>&city=%d&bd=%s` и транслирую ответ пользователю", *user.CityID, user.BirthDate.Format("2006-01-02 15:04")), nil)
		break
	case "diet":
		SendMarkDownMessage(message.Chat.ID, fmt.Sprintf("По диете так же на отдельный урл `https://bodygraph.online/api_v1/diet.php?dkey=<bot_key>&city=%d&bd=%s` и транслирую ответ пользователю", *user.CityID, user.BirthDate.Format("2006-01-02 15:04")), nil)
		break
	case "else":
		SendMarkDownMessage(message.Chat.ID, fmt.Sprintf("Тут наверное тоже что-то будет"), nil)
		break
	case "help":
		SendMarkDownMessage(message.Chat.ID, fmt.Sprintf("Тут дёрнем урл с помощью `https://bodygraph.online/api_v1/help.php` и транслирую ответ пользователю", *user.CityID, user.BirthDate.Format("2006-01-02 15:04")), nil)
		break
	case "settings":
		return SendSettingsMenu(message.Chat.ID)

	case "balance":
		return SendBalanceMenu(message.Chat.ID)

	}

	return SendStartMessage(message.Chat)

}

func SendSettingsMenu(chatID int64) int {
	log.Println("→ SendSettingsMenu")

	buttons := ButtonsSet(
		NewButton("💀 Сбросить данные о рождении", "settings:reset"),
		NewButton("🫥 Что-то ещё", "settings:else"),
		NewButton("⤴️ В главное меню", "mainmenu"),
	)

	return SendMessageWithReplyMarkup(chatID, "Настройки", buttons, onSettingsSelected)
}

func onSettingsSelected(message *models.Message, callbackData string) int {
	log.Println("→ onSettingsSelected")
	user, _ := db.GetUser(message.Chat)

	switch callbackData {
	case "settings:reset":
		user.BirthDate = nil
		user.CityID = nil
		user.CityName = nil
		user.CountryID = nil
		user.CountryName = nil
		err := db.UpdateUser(&user)
		if err != nil {
			return SendErrorMessage(message.Chat.ID, err)

		}
		SendMessage(message.Chat.ID, "Данные о рождении сброшены", nil)
		return SendStartMessage(message.Chat)

	case "balance":
		return SendMarkDownMessage(message.Chat.ID, fmt.Sprintf("Просто что-то ещё"), nil)

	case "mainmenu":
		return SendStartMessage(message.Chat)

	}

	log.Println("Как мы тут оказались? 🤔")
	return 1

}

func SendBalanceMenu(chatID int64) int {
	log.Println("→ SendBalanceMenu")

	buttons := ButtonsSet(NewButton("➕ Пополнить", "add"), NewButton("⤴️ В главное меню", "mainmenu"))

	return SendMessageWithReplyMarkup(chatID, "Ваш баланс 1 000 000 рублей", buttons, onBalanceSelected)

}

func onBalanceSelected(message *models.Message, callbackData string) int {
	log.Println("→ onBalanceSelected")
	//user, _ := db.GetUser(message.Chat)
	//log.Println("User: ", user)

	switch callbackData {
	case "add":
		return SendMarkDownMessage(message.Chat.ID, fmt.Sprintf("Тут будем пополнять баланс"), nil)
	case "mainmenu":
		return SendStartMessage(message.Chat)
	}

	log.Println("Как мы тут оказались? 🤔")
	return 1
}
