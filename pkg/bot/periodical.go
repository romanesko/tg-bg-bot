package botlogic

import (
	"bodygraph-bot/pkg/api"
	"bodygraph-bot/pkg/common"
	"encoding/json"
	"log"
	"time"
)

func CheckMessagesToSend(url string, interval int) {
	println("Starting periodical message checker for every", interval, "seconds")
	for {
		fetchMessagesUrl(url)
		time.Sleep(time.Second * time.Duration(interval))

	}
}

func CheckActionsToProcess(url string, interval int) {
	println("Starting periodical actions checker for every", interval, "seconds")
	for {
		fetchActionsUrl(url)
		time.Sleep(time.Second * time.Duration(interval))

	}
}

var blocked []int64
var ok []int64
var unknown []int64

type SendMessageStatus int

// Define constants of type Status to represent the possible return values
const (
	OK SendMessageStatus = iota
	BLOCKED
	UNKNOWN
)

func fetchMessagesUrl(url string) {

	log.Println("Getting queue to send")

	var params = map[string]any{
		"success": ok,
		"blocked": blocked,
		"unknown": unknown,
	}
	bodyString, err := api.FetchUrlAbstract(url, params)
	if err != nil {
		log.Println("Error fetching data", err)
		return
	}

	result := common.QueueResponse{}
	err = json.Unmarshal(bodyString, &result)
	if err != nil {

		var generic interface{}
		_ = json.Unmarshal(bodyString, &generic)

		log.Println("Error decoding response", generic)
		return
	}

	ok = nil
	blocked = nil
	unknown = nil

	if result.Response.Items == nil || len(*result.Response.Items) == 0 {
		log.Println("Nothing to send")
		return
	}

	log.Println("Sending", len(*result.Response.Items), "messages")

	for _, item := range *result.Response.Items {
		state := sendQueueMessage(item)

		if state == OK {
			ok = append(ok, item.MessageID)
		} else if state == BLOCKED {
			blocked = append(blocked, item.MessageID)
		} else if state == UNKNOWN {
			unknown = append(unknown, item.MessageID)
		}
	}

}

var checkUserInChannel []common.ActionsCheckUserInChannel

func fetchActionsUrl(url string) {

	log.Println("Getting queue to send")

	req := common.ActionsDTO{}

	req.CheckUserInChannel = checkUserInChannel

	bodyString, err := api.FetchUrlAbstract(url, req)
	if err != nil {
		log.Println("Error fetching data", err)
		return
	}

	result := common.ActionsResponse{}
	err = json.Unmarshal(bodyString, &result)
	if err != nil {
		var generic interface{}
		_ = json.Unmarshal(bodyString, &generic)
		log.Println("Error decoding response", generic)
		return
	}

	log.Println("items in channel")

	var checked []common.ActionsCheckUserInChannel

	for _, item := range result.Response.CheckUserInChannel {
		item.Exists = CheckUserInChannel(item.TgChatID, item.TgChannel)
		checked = append(checked, item)
	}

	checkUserInChannel = checked

}

func sendQueueMessage(item common.QueueItem) SendMessageStatus {
	log.Println("Sending message id", item.MessageID, "to", item.TgChatID)
	err := SendMessageData(item.TgChatID, item.Data)
	if err != nil {
		log.Println("Error sending message", err)
		return UNKNOWN
	}
	return OK
}
