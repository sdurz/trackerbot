package main

import (
	"context"
	"time"

	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

func messagePosition(message axon.O) (result *Position, err error) {
	var (
		location axon.O
		date     int64
	)

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

func doWhileLockingChatId(bot *ubot.Bot, chatIdPropertyPath string, message axon.O, dispatchedFunc func(*Chat) error) (err error) {
	if chatId, err := message.GetInteger(chatIdPropertyPath); err == nil {
		err = dispatcher.AcquireAndExecute(chatId, func() error {
			var chat *Chat
			if statusI, ok := lrucache.Get(chatId); ok {
				chat = statusI.(*Chat)
			} else {
				chat = NewChatStatus(bot, chatId)
				lrucache.Add(chatId, chat)
			}
			return dispatchedFunc(chat)
		})
	}
	return
}

func MessagePositionHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	err = doWhileLockingChatId(bot, "chat.id", message, func(chat *Chat) (err error) {
		if position, err := messagePosition(message); err == nil {
			chat.BeginTracking(bot, position)
		}
		done = true
		return
	})
	return
}

func MessagePositionUpdateHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	err = doWhileLockingChatId(bot, "chat.id", message, func(chat *Chat) (err error) {
		if position, err := messagePosition(message); err == nil {
			chat.state.UpdateTracking(bot, position)
		}
		return
	})
	return
}

func GetGpxCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	err = doWhileLockingChatId(bot, "chat.id", message, func(chat *Chat) (err error) {
		chat.SendGPX(bot)
		done = true
		return
	})
	return
}

func GetHelpCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	err = doWhileLockingChatId(bot, "chat.id", message, func(chat *Chat) (err error) {
		bot.SendMessage(axon.O{
			"chat_id":    chat.chatId,
			"text":       helpMarkup,
			"parse_mode": "MarkdownV2",
		})
		done = true
		return
	})
	return
}

func StartCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	err = doWhileLockingChatId(bot, "chat.id", message, func(chat *Chat) (err error) {
		err = chat.StartBot(bot, message)
		done = true
		return
	})
	return
}

func StopCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	err = doWhileLockingChatId(bot, "chat.id", message, func(chat *Chat) (err error) {
		_, err = bot.SendMessage(axon.O{
			"chat_id": chat.chatId,
			"text":    "Goodbye!",
		})
		lrucache.Remove(chat.chatId)
		done = true
		return
	})
	return
}

func SetProfileMessageHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	err = doWhileLockingChatId(bot, "chat.id", message, func(chat *Chat) (err error) {
		bot.SendMessage(axon.O{
			"chat_id": chat.chatId,
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
		done = true
		return
	})
	return
}

func PauseTrackingCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	err = doWhileLockingChatId(bot, "chat.id", message, func(chat *Chat) (err error) {
		err = chat.PauseTracking(bot)
		done = true
		return
	})
	return
}

func ResumeTrackingCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	err = doWhileLockingChatId(bot, "chat.id", message, func(chat *Chat) (err error) {
		chat.ResumeTracking(bot)
		done = true
		return
	})
	return
}

func EndTrackingCommandHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	err = doWhileLockingChatId(bot, "chat.id", message, func(chat *Chat) (err error) {
		chat.EndTracking(bot)
		done = true
		return
	})
	return
}

func CommandMessageHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	err = doWhileLockingChatId(bot, "chat.id", message, func(chat *Chat) (err error) {
		text, _ := message.GetString("text")
		messageId, _ := message.GetInteger("message_id")
		switch text {
		case "Pause":
			err = chat.PauseTracking(bot)
		case "Resume":
			chat.ResumeTracking(bot)
		case "End":
			chat.EndTracking(bot)
		case "Get GPX":
			_, err = chat.SendGPX(bot)
		}
		bot.DeleteMessage(axon.O{
			"chat_id":    chat.chatId,
			"message_id": messageId,
		})
		done = true
		return
	})
	return
}

func CallbackQueryHandler(ctx context.Context, bot *ubot.Bot, message axon.O) (done bool, err error) {
	err = doWhileLockingChatId(bot, "message.chat.id", message, func(chat *Chat) error {
		data, _ := message.GetString("data")
		result := chat.Callback(bot, data)
		if callbackQueryId, err := message.GetString("id"); err == nil {
			cbBody := axon.O{
				"callback_query_id": callbackQueryId,
			}
			if result != "" {
				cbBody["text"] = result
			}
			bot.AnswerCallbackQuery(cbBody)
		}
		return nil
	})
	return
}
