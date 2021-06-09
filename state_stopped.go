package main

import (
	"errors"

	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StateStopped struct {
	StateBase
	positions     []*Position
	downloadCount int64
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

func (state *StateStopped) GetGPX(bot *ubot.Bot, matchType string) (reesult []byte, err error) {
	if state.downloadCount < 10 {
		reesult, err = makeGpx(state.positions, matchType)
		state.downloadCount++
		if state.downloadCount == 10 {
			bot.SendMessage(axon.O{
				"chat_id": state.parent.chatId,
				"text":    "Maximum no. of downloads reached",
				"reply_markup": axon.O{
					"remove_keyboard": true,
				},
			})
		}
	} else {
		bot.SendMessage(axon.O{
			"chat_id": state.parent.chatId,
			"text":    "Max no. of downloads exceeded",
		})
		err = errors.New("to many downloads")
	}
	return
}
