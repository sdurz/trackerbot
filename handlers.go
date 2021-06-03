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

func MessagePositionHandler(ctx context.Context, b *ubot.Bot, message axon.O) (done bool, err error) {
	var (
		chatId   int64
		position *Position
		status   *ChatStatus
	)

	if chatId, position, err = messagePosition(message); err == nil {
		status = NewChatStatus(chatId)
		lrucache.Add(chatId, status)

		if err = status.StartTracking(b, position); err != nil {
			log.Println(err)
			lrucache.Remove(chatId)
		}
	}
	return
}

func MessagePositionUpdateHandler(ctx context.Context, b *ubot.Bot, message axon.O) (done bool, err error) {
	var (
		chatId   int64
		position *Position
		status   *ChatStatus
	)

	if chatId, position, err = messagePosition(message); err == nil {
		if cached, ok := lrucache.Get(chatId); ok {
			status = cached.(*ChatStatus)
			status.state.Update(b, position)
		}
	}
	return
}

func GetGpxCommandHandler(ctx context.Context, b *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	if statusI, ok := lrucache.Get(chatId); ok {
		status := statusI.(*ChatStatus)
		status.SendGPX(b)
	}
	return
}

func PauseTrackingCommandHandler(ctx context.Context, b *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	if cached, ok := lrucache.Get(chatId); ok {
		status := cached.(*ChatStatus)
		status.PauseTracking(b)
	}
	return
}

func ResumeTrackingCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	if cached, ok := lrucache.Get(chatId); ok {
		status := cached.(*ChatStatus)
		status.ResumeTracking(bot)
	}
	return
}

func CallbackQueryHandler(ctx context.Context, b *ubot.Bot, message axon.O) (done bool, err error) {
	log.Println("callback query")
	chatId, _ := message.GetInteger("chat.id")
	if cached, ok := lrucache.Get(chatId); ok {
		status := cached.(*ChatStatus)
		data, _ := message.GetString("data")
		status.Callback(b, data)
	}
	return
}
