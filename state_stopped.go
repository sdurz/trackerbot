package main

import (
	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StateStopped struct {
	parent    *ChatStatus
	positions []*Position
}

func (state *StateStopped) EnterState(bot *ubot.Bot, chatId int64) (err error) {
	_, err = bot.EditMessageText(axon.O{
		"chat_id":    chatId,
		"text":       "Tracking complete! Share your position to start tracking again",
		"parse_mode": "MarkdownV2",
		"reply_markup": axon.O{
			"remove_keyboard": true,
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

func (state *StateStopped) StartTracking(bot *ubot.Bot, position *Position) (err error) {
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
	if statusMessage, err := bot.SendMessage(axon.O{
		"chat_id":    state.parent.chatId,
		"text":       "Start tracking",
		"parse_mode": "MarkdownV2",
		"reply_markup": axon.O{
			"keyboard": axon.A{
				axon.A{
					axon.O{
						"text":          "Pause",
						"callback_data": "pause",
					},
				},
			},
			"resize_keyboard": true,
		},
	}); err == nil {
		state.parent.statusMessage = statusMessage
		state.parent.SetState(
			bot,
			&StateRunning{
				parent:    state.parent,
				positions: []*Position{position},
			})
	}
	return
}

func (state *StateStopped) GetKeyboard(bot *ubot.Bot) axon.O {
	return axon.O{}
}

func (state *StateStopped) GetGPX(bot *ubot.Bot) (data []byte, err error) {
	return makeGpx(state.positions)
}

func (state *StateStopped) GetCurrentPace() (result *Pace) {
	return
}
