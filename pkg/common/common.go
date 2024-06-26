package common

import (
	"fmt"
	"github.com/fatih/color"
	"log"
)

var userContext = make(map[int64]*MessageData)

func SetUserContext(userId int64, contextValue *MessageData) {
	FuncLog("SetUserContext", userId, contextValue)
	userContext[userId] = contextValue
}

func GetUserContext(userId int64) *MessageData {
	return userContext[userId]
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
