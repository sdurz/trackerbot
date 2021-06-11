package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type Chat struct {
	chatId        int64
	vehicle       string
	state         State
	statusMessage axon.O
	pinnedMessage axon.O
}

func NewChatStatus(bot *ubot.Bot, chatId int64) (result *Chat) {
	result = &Chat{
		chatId:  chatId,
		vehicle: "hike",
	}
	result.SetState(bot, &StateReady{
		StateBase: StateBase{
			parent: result,
		},
	},
		nil)
	return
}

func (s *Chat) StartBot(bot *ubot.Bot, message axon.O) (err error) {
	err = s.state.Start(bot, message)
	return
}

func (s *Chat) BeginTracking(bot *ubot.Bot, position *Position) (err error) {
	err = s.state.BeginTracking(bot, position)
	return
}

func (s *Chat) PauseTracking(bot *ubot.Bot) (err error) {
	err = s.state.PauseTracking(bot)
	return
}

func (s *Chat) ResumeTracking(bot *ubot.Bot) {
	s.state.ResumeTracking(bot)
}

func (s *Chat) EndTracking(bot *ubot.Bot) {
	s.state.EndTracking(bot)
}

func (s *Chat) UpdatePosition(bot *ubot.Bot, position *Position) {
	s.state.UpdateTracking(bot, position)
}

func (s *Chat) SendGPX(bot *ubot.Bot) (result []byte, err error) {
	if byteData, mapMatched, err := s.state.GetGPX(bot, s.vehicle); err == nil {
		if !mapMatched {
			bot.SendMessage(axon.O{
				"text":       "*Warning*: map matching failed, returning raw positions only",
				"parse_mode": "MarkdownV2",
			})
		}

		fileName := fmt.Sprintf("TelegramTrack-%v.gpx", time.Now().Format("20060102-150405"))
		if uploadFile, err := ubot.NewBytesUploadFile(fileName, byteData); err == nil {
			bot.SendDocument(axon.O{
				"chat_id":  s.chatId,
				"document": uploadFile,
			})
		}
	}
	return
}

func (status *Chat) SetState(bot *ubot.Bot, state State, maessage axon.O) (err error) {
	if err := state.EnterState(bot, nil); err == nil {
		status.state = state
	} else {
		log.Println("Error in EnterState: " + err.Error())
	}
	return
}

func (status *Chat) Callback(bot *ubot.Bot, data string) (result string) {
	result = ""
	switch data {
	case "pause tracking":
		status.state.PauseTracking(bot)
		result = "Tracking paused"
	case "resume tracking":
		status.state.ResumeTracking(bot)
		result = "Tracking resumed"
	case "end tracking":
		status.state.EndTracking(bot)
		result = "Tracking ended"
	case "set bike":
		status.vehicle = "bike"
		result = "Bike profile set"
	case "set hike":
		status.vehicle = "hike"
		result = "Hike profile set"
	case "set car":
		status.vehicle = "car"
		result = "Car profile set"
	default:
	}
	return
}
