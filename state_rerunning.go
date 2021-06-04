package main

import (
	"time"

	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StateRerunning struct {
	StateRunning
}

func (state *StateRerunning) EnterState(bot *ubot.Bot, chatId int64) (err error) {
	if message, err := bot.SendMessage(axon.O{
		"chat_id":    state.parent.chatId,
		"text":       "Tracking **restarted** at " + time.Now().Format("15:04:05"),
		"parse_mode": "MarkdownV2",
		"reply_markup": axon.O{
			"keyboard": axon.A{
				axon.A{
					axon.O{
						"text": "Stop",
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
