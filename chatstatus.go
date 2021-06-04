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
	state         State
	statusMessage axon.O
	matchMode     string
}

func NewChatStatus(bot *ubot.Bot, chatId int64) (result *ChatStatus) {
	result = &ChatStatus{
		chatId:    chatId,
		matchMode: "hike",
	}
	result.SetState(bot, &StateReady{result})
	return
}

func (s *ChatStatus) StartTracking(bot *ubot.Bot, position *Position) (err error) {
	err = s.state.Start(bot, position)
	return
}

func (s *ChatStatus) PauseTracking(bot *ubot.Bot) (err error) {
	err = s.state.Pause(bot)
	return
}

func (s *ChatStatus) ResumeTracking(bot *ubot.Bot) {
	s.state.Resume(bot)
}

func (s *ChatStatus) StopTracking(bot *ubot.Bot) {
	s.state.Stop(bot)
}

func (s *ChatStatus) Append(bot *ubot.Bot, position *Position) {
	s.state.Update(bot, position)
}

func (s *ChatStatus) SendGPX(bot *ubot.Bot, matchType string) (result []byte, err error) {
	if byteData, err := s.state.GetGPX(bot, matchType); err == nil {
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
	case "pause":
		status.state.Pause(bot)
	case "stop":
		status.state.Stop(bot)
	case "resume":
		status.state.Resume(bot)
	default:
	}
}
