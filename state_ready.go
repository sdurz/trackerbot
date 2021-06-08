package main

import (
	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type StateReady struct {
	StateBase
}

func (state *StateReady) EnterState(bot *ubot.Bot, message axon.O) (err error) {
	return
}

func (state *StateReady) Start(bot *ubot.Bot, message axon.O) (err error) {
	err = state.parent.SetState(bot, &StateStarted{
		StateBase: StateBase{
			parent: state.parent,
		},
	}, message)
	return
}

func (state *StateReady) BeginTracking(bot *ubot.Bot, position *Position) (err error) {
	state.parent.SetState(bot, &StateTracking{
		StateBase: StateBase{
			parent: state.parent,
		},
	},
		nil)
	return
}
