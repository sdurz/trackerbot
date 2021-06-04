package main

import (
	"github.com/sdurz/ubot"
)

type State interface {
	EnterState(bot *ubot.Bot, chatId int64) (err error)
	Start(bot *ubot.Bot, position *Position) (err error)
	Pause(bot *ubot.Bot) (err error)
	Resume(bot *ubot.Bot) (err error)
	Update(bot *ubot.Bot, position *Position) (err error)
	Stop(bot *ubot.Bot) (err error)
	GetGPX(bot *ubot.Bot, matchType string) (data []byte, err error)
	GetCurrentPace() *Pace
}
