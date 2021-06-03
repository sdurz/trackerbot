package main

import (
	"errors"

	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StateInitial struct {
	parent *ChatStatus
}

func (state *StateInitial) Start(bot *ubot.Bot, position *Position) (err error) {
	err = state.parent.SetState(
		bot,
		&StateRunning{
			parent:    state.parent,
			positions: []*Position{position},
		},
	)
	return
}

func (state *StateInitial) EnterState(bot *ubot.Bot, chatId int64) (err error) {
	if message, err := bot.SendMessage(axon.O{
		"chat_id": chatId,
		"text":    "Hello! Share your position to start tracking",
		"reply_markup": axon.O{
			"remove_keyboard": true,
		},
	}); err == nil {
		state.parent.statusMessage = message
		state.parent.chatId = chatId
	}
	return
}

func (state *StateInitial) Pause(bot *ubot.Bot) (err error) {
	// no op
	return
}

func (state *StateInitial) Resume(bot *ubot.Bot) (err error) {
	// no op
	return
}

func (state *StateInitial) Update(bot *ubot.Bot, position *Position) (err error) {
	// no op
	return
}

func (state *StateInitial) Stop(bot *ubot.Bot) (err error) {
	// no op
	return
}

func (state *StateInitial) GetGPX(ubot *ubot.Bot) (data []byte, err error) {
	// no op
	err = errors.New("no tracking data in idle state")
	return
}

func (state *StateInitial) GetCurrentPace() (result *Pace) {
	return
}
