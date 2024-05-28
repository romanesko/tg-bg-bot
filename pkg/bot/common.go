package bot

import "github.com/go-telegram/bot/models"

var userContext = make(map[int64]func(message *models.Message, callbackData string) int)

func SetUserContext(userId int64, contextValue func(message *models.Message, callbackData string) int) {
	userContext[userId] = contextValue
}

func GetUserContext(userId int64) func(message *models.Message, callbackData string) int {
	return userContext[userId]
}
