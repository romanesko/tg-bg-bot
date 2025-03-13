package botlogic

import (
	"bodygraph-bot/pkg/api"
	"bodygraph-bot/pkg/common"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot/models"
	"log"
	"sync"
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

	checked := make([]common.ActionsCheckUserInChannel, 0)

	channelGroup := make(map[string][]common.ActionsCheckUserInChannel)

	minChannelLength := 0
	minUserLength := 0

	for _, item := range result.Response.CheckUserInChannel {
		if _, ok := channelGroup[item.TgChannel]; !ok {
			channelGroup[item.TgChannel] = make([]common.ActionsCheckUserInChannel, 0)
			if len(item.TgChannel) > minChannelLength {
				minChannelLength = len(item.TgChannel)
			}
		}

		channelGroup[item.TgChannel] = append(channelGroup[item.TgChannel], item)
		chatIdStr := fmt.Sprintf("%d", item.TgChatID)
		if len(chatIdStr) > minUserLength {
			minUserLength = len(chatIdStr)
		}

	}

	for channelName, items := range channelGroup {

		channel, err := getChannelByName(channelName)
		if err != nil {
			log.Printf("CHECKING CHANNEL «%s»: ERROR: %s", channelName, common.UnwrapError(err))
			for _, item := range items {
				item.Exists = false
				item.State = "channel-not-found"
				checked = append(checked, item)
			}
			continue
		}

		log.Printf("CHECKING CHANNEL «%s»:", channel.Title)

		var wg sync.WaitGroup
		semaphore := make(chan struct{}, 50)
		mu := sync.Mutex{}

		for idx, item := range items {
			wg.Add(1)
			semaphore <- struct{}{}
			go func(index int, it common.ActionsCheckUserInChannel, channel models.ChatFullInfo) {
				defer wg.Done()
				item.Exists, item.State = CheckUserInChannel(item.TgChatID, channel, minUserLength)
				mu.Lock()
				checked = append(checked, item)
				mu.Unlock()
				<-semaphore
			}(idx, item, *channel)
		}
		wg.Wait()
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
