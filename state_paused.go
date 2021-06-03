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
	messageId, _ := state.parent.statusMessage.GetInteger("message_id")
	_, err = bot.EditMessageText(axon.O{
		"chat_id":    state.parent.chatId,
		"message_id": messageId,
		"text":       "Tracking started",
		"parse_mode": "MarkdownV2",
		"reply_markup": axon.O{
			"keyboard": axon.A{
				axon.A{
					axon.O{
						"text":          "Resume",
						"callback_data": "resume",
					},
				},
			},
			"resize_keyboard": true,
		},
	})
	return
}

func (state *StatePaused) Start(bot *ubot.Bot, position *Position) (err error) {
	return
}

func (state *StatePaused) Pause(bot *ubot.Bot) (err error) {
	return
}

func (state *StatePaused) Resume(bot *ubot.Bot) (err error) {
	messageId, _ := state.parent.statusMessage.GetInteger("message_id")
	if _, err = bot.EditMessageText(axon.O{
		"chat_id":    state.parent.chatId,
		"message_id": messageId,
		"parse_mode": "MarkdownV2",
		"text":       "Tracking restarted at" + time.Now().Format(time.RFC3339),
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
	}); err != nil {
		return
	}
	err = state.parent.SetState(bot,
		&StateRunning{
			positions: state.positions,
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
