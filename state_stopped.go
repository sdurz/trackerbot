package main

import (
	"errors"

	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StateStopped struct {
	parent        *ChatStatus
	positions     []*Position
	downloadCount int64
}

func (state *StateStopped) EnterState(bot *ubot.Bot, chatId int64) (err error) {
	_, err = bot.SendMessage(axon.O{
		"chat_id": state.parent.chatId,
		"text":    "Tracking complete! Share your position to start tracking again",
		"reply_markup": axon.O{
			"keyboard": axon.A{
				axon.A{
					axon.O{
						"text": "Get GPX",
					},
				},
			},
		},
	})
	return
}

func (state *StateStopped) Pause(bot *ubot.Bot) (err error) {
	return
}

func (state *StateStopped) Resume(bot *ubot.Bot) (err error) {
	// no op
	return
}

func (state *StateStopped) Update(bot *ubot.Bot, position *Position) (err error) {
	return
}

func (state *StateStopped) Stop(bot *ubot.Bot) (err error) {
	// no op
	return
}

func (state *StateStopped) Start(bot *ubot.Bot, position *Position) (err error) {
	err = state.parent.SetState(
		bot,
		&StateRunning{
			parent:    state.parent,
			positions: []*Position{position},
		},
	)
	return
}

func (state *StateStopped) GetGPX(bot *ubot.Bot) (reesult []byte, err error) {
	if state.downloadCount < 3 {
		reesult, err = makeGpx(state.positions)
		state.downloadCount++
		if state.downloadCount == 3 {
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

func (state *StateStopped) GetCurrentPace() (result *Pace) {
	return
}
