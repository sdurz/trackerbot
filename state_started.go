package main

import (
	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StateStarted struct {
	StateBase
}

func (state *StateStarted) EnterState(bot *ubot.Bot, message axon.O) (err error) {
	userName, _ := message.GetString("")
	bot.SendMessage(axon.O{
		"chat_id": state.parent.chatId,
		"text":    "Hello " + userName + "! Share your position to start tracking",
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
			StateBase: state.StateBase,
			positions: []*Position{position},
		},
		nil,
	)
	return
}
