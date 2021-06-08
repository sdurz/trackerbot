package main

import (
	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StateBase struct {
	parent *ChatStatus
}

func (state *StateBase) Start(bot *ubot.Bot, message axon.O) (err error) {
	// no op
	return
}

func (state *StateBase) BeginTracking(bot *ubot.Bot, position *Position) (err error) {
	return
}

func (state *StateBase) PauseTracking(bot *ubot.Bot) (err error) {
	// no op
	return
}

func (state *StateBase) ResumeTracking(bot *ubot.Bot) (err error) {
	// no op
	return
}

func (state *StateBase) UpdateTracking(bot *ubot.Bot, position *Position) (err error) {
	// no op
	return
}

func (state *StateBase) EndTracking(bot *ubot.Bot) (err error) {
	// no op
	return
}

func (state *StateBase) GetGPX(ubot *ubot.Bot, matchType string) (data []byte, err error) {
	// no op
	return
}

func (state *StateBase) GetCurrentPace() (result *Pace) {
	return
}
