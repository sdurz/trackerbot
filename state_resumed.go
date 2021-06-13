package main

import (
	"time"

	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StateResumed struct {
	StateTracking
}

func (state *StateResumed) EnterState(bot *ubot.Bot, message axon.O) (err error) {
	pinnedId, _ := state.parent.pinnedMessage.GetInteger("message_id")
	bot.EditMessageText(axon.O{
		"chat_id":    state.parent.chatId,
		"message_id": pinnedId,
		"text":       "State: tracking",
		"parse_mode": "MarkdownV2",
	})
	if message, err := bot.SendMessage(axon.O{
		"chat_id":    state.parent.chatId,
		"text":       "Tracking *resumed* at " + time.Now().Format("15:04:05"),
		"parse_mode": "MarkdownV2",
		"reply_markup": axon.O{
			"keyboard": axon.A{
				axon.A{
					axon.O{
						"text": "End",
					},
				},
				axon.A{
					axon.O{
						"text": "Pause",
					},
				},
			},
			"resize_keyboard": true,
		},
	}); err == nil {
		state.parent.statusMessage = message
	}
	return
}
