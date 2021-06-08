package main

import (
	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StateStarted struct {
	StateBase
	parent *ChatStatus
}

func (state *StateStarted) EnterState(bot *ubot.Bot, message axon.O) (err error) {
	bot.SendMessage(axon.O{
		"chat_id": state.parent.chatId,
		"text":    "Hello! Share your position to start tracking",
		"reply_markup": axon.O{
			"remove_keyboard": true,
		},
	})
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
		nil,
	)
	return
}
