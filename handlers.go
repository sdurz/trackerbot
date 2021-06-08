package main

import (
	"context"
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
		status.BeginTracking(bot, position)
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
		status.state.UpdateTracking(bot, position)
	}
	done = true
	return
}

func GetGpxCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	status := findOrCreateStatus(bot, chatId)
	status.SendGPX(bot)
	done = true
	return
}

func GetHelpCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	bot.SendMessage(axon.O{
		"chat_id":    chatId,
		"text":       helpMarkup,
		"parse_mode": "MarkdownV2",
	})
	done = true
	return
}

func StartCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	status := findOrCreateStatus(bot, chatId)
	status.StartBot(bot, message)
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

func SetProfileCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	bot.SendMessage(axon.O{
		"chat_id": chatId,
		"text":    "Choose your vehicle",
		"reply_markup": axon.O{
			"inline_keyboard": axon.A{
				axon.A{
					axon.O{
						"text":          "🚶🏽 Hike",
						"callback_data": "set hike",
					},
				},
				axon.A{
					axon.O{
						"text":          "🚴 Bike",
						"callback_data": "set bike",
					},
				},
				axon.A{
					axon.O{
						"text":          "🚗 Car",
						"callback_data": "set car",
					},
				},
			},
		},
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

func EndTrackingCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	chatId, _ := message.GetInteger("chat.id")
	status := findOrCreateStatus(bot, chatId)
	status.EndTracking(bot)
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
		case "End":
			status.EndTracking(bot)
		case "Get GPX":
			status.SendGPX(bot)
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
	chatId, _ := message.GetInteger("message.chat.id")
	data, _ := message.GetString("data")
	status := findOrCreateStatus(bot, chatId)
	result := status.Callback(bot, data)

	if callbackQueryId, err := message.GetString("id"); err == nil {
		cbBody := axon.O{
			"callback_query_id": callbackQueryId,
		}
		if result != "" {
			cbBody["text"] = result
		}
		bot.AnswerCallbackQuery(cbBody)
	}
	return
}
