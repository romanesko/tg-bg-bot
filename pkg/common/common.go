package common

import (
	"bodygraph-bot/pkg/kvstore"
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-telegram/bot/models"
	"log"
	"regexp"
	"strings"
)

var userContext = make(map[int64]*MessageData)
var UserInfoMap = make(map[int64]*UserInfo)

type UserInfo struct {
	FirstName string
	LastName  string
	Username  string
	Start     string
}

func SetUserInfo(chat *models.Chat) {
	var u = &UserInfo{}
	u.FirstName = chat.FirstName
	u.LastName = chat.LastName
	u.Username = chat.Username
	UserInfoMap[chat.ID] = u
}

func SetUserStartCommandParams(chatId int64, start string) {
	if UserInfoMap[chatId] == nil {
		log.Printf("FAILED to set referral for missing chatId in UserInfoMap %d\n", chatId)
		return
	}

	UserInfoMap[chatId].Start = start

}

func GetUserInfo(userId int64) UserInfo {
	var ui = UserInfoMap[userId]
	if ui == nil {
		return UserInfo{}
	}
	return *ui
}

func SetUserContext(userId int64, contextValue *MessageData) {
	FuncLog("SetUserContext", userId, contextValue)

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(contextValue)
	if err != nil {
		log.Printf("Encoding error: %v", err)
		return
	}

	kvstore.Write(fmt.Sprintf("USER:%d", userId), buffer.Bytes(), 60*60*24*30)

	userContext[userId] = contextValue
}

func GetUserContext(userId int64) *MessageData {

	serializedData := kvstore.Read(fmt.Sprintf("USER:%d", userId))
	if serializedData == nil {
		return nil
	}

	var buffer bytes.Buffer

	buffer = *bytes.NewBuffer(serializedData) // Recreate buffer from bytes
	var decoded MessageData
	decoder := gob.NewDecoder(&buffer)
	err := decoder.Decode(&decoded)

	if err != nil {
		log.Printf("Decoding error: %v", err)
	}

	return &decoded

	//return userContext[userId]
}

func FuncLog(s ...any) {
	color.Set(color.FgBlue)
	str := ""
	for _, v := range s {
		str += fmt.Sprintf("%v ", v)
	}
	log.Printf("→ %s\n", str)
	color.Unset()
}

func ErrorLog(s any) {
	color.Set(color.FgRed)
	log.Printf("✖ %s\n", s)
	color.Unset()
}

func RemoveUnwantedTags(input string) string {
	allowedTags := `(?i)(&lt;\/?(b|strong|i|em|u|ins|s|strike|del|a|code|pre)\b[^&]*&gt;)`

	escaped := strings.ReplaceAll(input, "<", "&lt;")
	escaped = strings.ReplaceAll(escaped, ">", "&gt;")

	re := regexp.MustCompile(allowedTags)

	// Revert allowed tags back to their original form
	result := re.ReplaceAllStringFunc(escaped, func(tag string) string {
		tag = strings.ReplaceAll(tag, "&lt;", "<")
		tag = strings.ReplaceAll(tag, "&gt;", ">")
		return tag
	})

	return result
}

func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
