package main

import (
	"errors"

	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StateReady struct {
	parent *ChatStatus
}

func (state *StateReady) Start(bot *ubot.Bot, position *Position) (err error) {
	err = state.parent.SetState(
		bot,
		&StateRunning{
			parent:    state.parent,
			positions: []*Position{position},
		},
	)
	return
}

func (state *StateReady) EnterState(bot *ubot.Bot, chatId int64) (err error) {
	if _, err := bot.SendMessage(axon.O{
		"chat_id": chatId,
		"text":    "Hello! Share your position to start tracking",
		"reply_markup": axon.O{
			"remove_keyboard": true,
		},
	}); err == nil {
		state.parent.chatId = chatId
	}
	return
}

func (state *StateReady) Pause(bot *ubot.Bot) (err error) {
	// no op
	return
}

func (state *StateReady) Resume(bot *ubot.Bot) (err error) {
	// no op
	return
}

func (state *StateReady) Update(bot *ubot.Bot, position *Position) (err error) {
	// no op
	return
}

func (state *StateReady) Stop(bot *ubot.Bot) (err error) {
	// no op
	return
}

func (state *StateReady) GetGPX(ubot *ubot.Bot) (data []byte, err error) {
	// no op
	err = errors.New("no tracking data in idle state")
	return
}

func (state *StateReady) GetCurrentPace() (result *Pace) {
	return
}
