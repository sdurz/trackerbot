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
	updateTime := time.Now().Format("15:04:05")
	pinnedId, _ := state.parent.pinnedMessage.GetInteger("message_id")
	bot.EditMessageText(axon.O{
		"chat_id":    state.parent.chatId,
		"message_id": pinnedId,
		"text":       "State: **paused**, Pace: --:--",
	})

	_, err = bot.SendMessage(axon.O{
		"chat_id": state.parent.chatId,
		"text":    "Tracking paused at " + updateTime,
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

func (state *StatePaused) BeginTracking(bot *ubot.Bot, position *Position) (err error) {
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

func (state *StatePaused) PauseTracking(bot *ubot.Bot) (err error) {
	return
}

func (state *StatePaused) ResumeTracking(bot *ubot.Bot) (err error) {
	err = state.parent.SetState(bot,
		&StateRerunning{
			StateRunning{
				parent:    state.parent,
				positions: state.positions,
			},
		})
	return
}

func (state *StatePaused) UpdateTracking(bot *ubot.Bot, position *Position) (err error) {
	return
}

func (state *StatePaused) EndTracking(bot *ubot.Bot) (err error) {
	return
}

func (state *StatePaused) GetGPX(ubot *ubot.Bot, matchType string) (data []byte, err error) {
	return makeGpx(state.positions, matchType)
}

func (state *StatePaused) GetCurrentPace() (result *Pace) {
	return
}
