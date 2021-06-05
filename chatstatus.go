package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

type ChatStatus struct {
	chatId        int64
	vehicle       string
	state         State
	statusMessage axon.O
}

func NewChatStatus(bot *ubot.Bot, chatId int64) (result *ChatStatus) {
	result = &ChatStatus{
		chatId:  chatId,
		vehicle: "hike",
	}
	result.SetState(bot, &StateReady{result})
	return
}

func (s *ChatStatus) StartTracking(bot *ubot.Bot, position *Position) (err error) {
	err = s.state.BeginTracking(bot, position)
	return
}

func (s *ChatStatus) PauseTracking(bot *ubot.Bot) (err error) {
	err = s.state.PauseTracking(bot)
	return
}

func (s *ChatStatus) ResumeTracking(bot *ubot.Bot) {
	s.state.ResumeTracking(bot)
}

func (s *ChatStatus) EndTracking(bot *ubot.Bot) {
	s.state.EndTracking(bot)
}

func (s *ChatStatus) Append(bot *ubot.Bot, position *Position) {
	s.state.UpdateTracking(bot, position)
}

func (s *ChatStatus) SendGPX(bot *ubot.Bot) (result []byte, err error) {
	if byteData, err := s.state.GetGPX(bot, s.vehicle); err == nil {
		fileName := fmt.Sprintf("TelegramTrack-%v.gpx", time.Now().Format("20060102-150405"))
		uploadFile, _ := ubot.NewBytesUploadFile(fileName, byteData)
		bot.SendDocument(axon.O{
			"chat_id":  s.chatId,
			"document": uploadFile,
		})
	}
	return
}

func (status *ChatStatus) SetState(bot *ubot.Bot, state State) (err error) {
	if err := state.EnterState(bot, status.chatId); err == nil {
		status.state = state
	} else {
		log.Println("Error in EnterState: " + err.Error())
		log.Println("state not changed")
	}
	return
}

func (status *ChatStatus) Callback(bot *ubot.Bot, data string) {
	switch data {
	case "pause tracking":
		status.state.PauseTracking(bot)
	case "stop tracking":
		status.state.EndTracking(bot)
	case "resume tracking":
		status.state.ResumeTracking(bot)
	case "set bike":
		status.vehicle = "bike"
	case "set hike":
		status.vehicle = "hike"
	default:
	}
}
