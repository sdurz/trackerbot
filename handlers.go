package main

import (
	"context"
	"log"
	"time"

	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

func messagePosition(message axon.O) (chatId int64, result *Position, err error) {
	var (
		location axon.O
		date     int64
	)

	if chatId, err = message.GetInteger("chat.id"); err != nil {
		return
	}
	if location, err = message.GetObject("location"); err != nil {
		return
	}
	if date, err = message.GetInteger("edit_date"); err != nil {
		if date, err = message.GetInteger("date"); err != nil {
			return
		}
	}
	result = &Position{
		time.Unix(date, 0),
		location["latitude"].(float64),
		location["longitude"].(float64),
	}
	return
}

func findOrCreateStatus(bot *ubot.Bot, chatId int64) (result *ChatStatus) {
	if statusI, ok := lrucache.Get(chatId); ok {
		result = statusI.(*ChatStatus)
	} else {
		result = NewChatStatus(bot, chatId)
		lrucache.Add(chatId, result)
	}
	return
}

func MessagePositionHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	var (
		chatId   int64
		position *Position
	)

	if chatId, position, err = messagePosition(message); err == nil {
		status := findOrCreateStatus(bot, chatId)
		status.StartTracking(bot, position)
	}
	done = true
	return
}

func MessagePositionUpdateHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	var (
		chatId   int64
		position *Position
	)

	if chatId, position, err = messagePosition(message); err == nil {
		status := findOrCreateStatus(bot, chatId)
		status.state.Update(bot, position)
	}
	done = true
	return
}

func GetGpxCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	status := findOrCreateStatus(bot, chatId)
	status.SendGPX(bot, "")
	done = true
	return
}

func StartCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	findOrCreateStatus(bot, chatId)
	done = true
	return
}

func StopCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	bot.SendMessage(axon.O{
		"chat_id": chatId,
		"text":    "Goodbye!",
	})
	lrucache.Remove(chatId)
	done = true
	return
}

func PauseTrackingCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	status := findOrCreateStatus(bot, chatId)
	status.PauseTracking(bot)
	done = true
	return
}

func ResumeTrackingCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	status := findOrCreateStatus(bot, chatId)
	status.ResumeTracking(bot)
	done = true
	return
}

func CommandMessageHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	if cached, ok := lrucache.Get(chatId); ok {
		status := cached.(*ChatStatus)
		text, _ := message.GetString("text")
		messageId, _ := message.GetInteger("message_id")
		switch text {
		case "Pause":
			status.PauseTracking(bot)
		case "Resume":
			status.ResumeTracking(bot)
		case "Stop":
			status.StopTracking(bot)
		case "Get GPX":
			status.SendGPX(bot, "")
		case btnBikeGPX:
			status.SendGPX(bot, "bike")
		case btnHikeGPX:
			status.SendGPX(bot, "hike")
		case btnCarGPX:
			status.SendGPX(bot, "car")
		}
		bot.DeleteMessage(axon.O{
			"chat_id":    chatId,
			"message_id": messageId,
		})
	}
	done = true
	return
}

func CallbackQueryHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	log.Println("callback query")
	chatId, _ := message.GetInteger("chat.id")
	data, _ := message.GetString("data")
	status := findOrCreateStatus(bot, chatId)
	status.Callback(bot, data)
	return
}
