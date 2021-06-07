package main

import (
	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StateStarted struct {
	parent *ChatStatus
}

func (state *StateStarted) EnterState(bot *ubot.Bot, chatId int64) (err error) {
	if _, err := bot.SendMessage(axon.O{
		"chat_id": chatId,
		"text":    "Hello! Share your position to start tracking",
		"reply_markup": axon.O{
			"remove_keyboard": true,
		},
	}); err == nil {
		state.parent.chatId = chatId
	}
	return
}

func (state *StateStarted) BeginTracking(bot *ubot.Bot, position *Position) (err error) {
	err = state.parent.SetState(
		bot,
		&StateTracking{
			StateBase: StateBase{
				parent: state.parent,
			},
			positions: []*Position{position},
		},
	)
	return
}
