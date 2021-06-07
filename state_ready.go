package main

import (
	"github.com/sdurz/ubot"
)

type StateReady struct {
	StateBase
}

func (state *StateReady) EnterState(bot *ubot.Bot, chatId int64) (err error) {
	// do nothing
	return
}

func (state *StateReady) BeginTracking(bot *ubot.Bot, position *Position) (err error) {
	state.parent.SetState(bot, &StateTracking{
		StateBase: StateBase{
			parent: state.parent,
		},
	})
	return
}
