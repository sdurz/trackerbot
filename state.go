package main

import (
	"github.com/sdurz/ubot"
)

type State interface {
	EnterState(bot *ubot.Bot, chatId int64) (err error)
	BeginTracking(bot *ubot.Bot, position *Position) (err error)
	PauseTracking(bot *ubot.Bot) (err error)
	ResumeTracking(bot *ubot.Bot) (err error)
	UpdateTracking(bot *ubot.Bot, position *Position) (err error)
	EndTracking(bot *ubot.Bot) (err error)
	GetGPX(bot *ubot.Bot, matchType string) (data []byte, err error)
	GetCurrentPace() *Pace
}
