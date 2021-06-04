package main

import (
	"time"

	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StatePaused struct {
	parent    *ChatStatus
	positions []*Position
}

func (state *StatePaused) EnterState(bot *ubot.Bot, chatId int64) (err error) {
	_, err = bot.SendMessage(axon.O{
		"chat_id": state.parent.chatId,
		"text":    "Tracking paused at " + time.Now().Format("15:04:05"),
		"reply_markup": axon.O{
			"keyboard": axon.A{
				axon.A{
					axon.O{
						"text": "Resume",
					},
				},
			},
			"resize_keyboard": true,
		},
	})
	return
}

func (state *StatePaused) Start(bot *ubot.Bot, position *Position) (err error) {
	bot.SendMessage(axon.O{
		"chat_id": state.parent.chatId,
		"text":    "Current tracking aborted, now restarting...",
	})
	err = state.parent.SetState(
		bot,
		&StateRunning{
			parent:    state.parent,
			positions: []*Position{position},
		},
	)
	return
}

func (state *StatePaused) Pause(bot *ubot.Bot) (err error) {
	return
}

func (state *StatePaused) Resume(bot *ubot.Bot) (err error) {
	err = state.parent.SetState(bot,
		&StateRerunning{
			StateRunning{
				parent:    state.parent,
				positions: state.positions,
			},
		})
	return
}

func (state *StatePaused) Update(bot *ubot.Bot, position *Position) (err error) {
	return
}

func (state *StatePaused) Stop(bot *ubot.Bot) (err error) {
	return
}

func (state *StatePaused) GetGPX(ubot *ubot.Bot) (data []byte, err error) {
	return makeGpx(state.positions)
}

func (state *StatePaused) GetCurrentPace() (result *Pace) {
	return
}
