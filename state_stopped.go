package main

import (
	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StateStopped struct {
	StateBase
	positions []*Position
}

func (state *StateStopped) EnterState(bot *ubot.Bot, message axon.O) (err error) {
	pinnedId, _ := state.parent.pinnedMessage.GetInteger("message_id")
	bot.EditMessageText(axon.O{
		"chat_id":    state.parent.chatId,
		"message_id": pinnedId,
		"text":       "State: **ended**",
	})
	_, err = bot.SendMessage(axon.O{
		"chat_id": state.parent.chatId,
		"text":    "Tracking complete! Share your position again to restart tracking",
		"reply_markup": axon.O{
			"remove_keyboard": true,
		},
	})
	state.parent.SendGPX(bot)
	return
}

func (state *StateStopped) BeginTracking(bot *ubot.Bot, position *Position) (err error) {
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

func (state *StateStopped) GetGPX(bot *ubot.Bot, matchType string) (result []byte, mapMatched bool, err error) {
	return makeGpx(state.positions, state.parent.vehicle)
}
