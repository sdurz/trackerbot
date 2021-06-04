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

func (state *StateStopped) GetGPX(bot *ubot.Bot) (data []byte, err error) {
	if state.downloadCount < 3 {
		data, err = makeGpx(state.positions)
	} else {
		bot.SendMessage(axon.O{
			"text": "Download exceeded",
		})
		err = errors.New("too many downloads")
	}
	return
}

func (state *StateStopped) GetCurrentPace() (result *Pace) {
	return
}
