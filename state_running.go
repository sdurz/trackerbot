package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
	"github.com/umahmood/haversine"
)

type StateRunning struct {
	parent    *ChatStatus
	positions []*Position
}

func (state *StateRunning) EnterState(bot *ubot.Bot, chatId int64) (err error) {
	messageId, _ := state.parent.statusMessage.GetInteger("message_id")
	_, err = bot.EditMessageText(axon.O{
		"chat_id":    state.parent.chatId,
		"message_id": messageId,
		"text":       "Tracking started",
		"parse_mode": "MarkdownV2",
		"reply_markup": axon.O{
			"keyboard": axon.A{
				axon.A{
					axon.O{
						"text":          "Pause",
						"callback_data": "pause",
					},
				},
			},
			"resize_keyboard": true,
		},
	})
	return
}

func (state *StateRunning) Start(bot *ubot.Bot, posiition *Position) (err error) {
	// no op
	return
}

func (state *StateRunning) Pause(bot *ubot.Bot) (err error) {
	state.parent.SetState(bot, &StatePaused{
		positions: state.positions,
	})
	return
}

func (state *StateRunning) Resume(bot *ubot.Bot) (err error) {
	// no op
	return
}

func (state *StateRunning) Update(bot *ubot.Bot, position *Position) (err error) {
	if position == nil {
		log.Fatalf("null position")
	}
	state.positions = append(state.positions, position)

	messageId, _ := state.parent.statusMessage.GetInteger("message_id")
	chatId, _ := state.parent.statusMessage.GetInteger("chat.id")
	bot.EditMessageText(axon.O{
		"chat_id":    chatId,
		"message_id": messageId,
		"parse_mode": "MarkdownV2",
		"text":       fmt.Sprintf("Current pace: **%s**\nUpdated: %s", state.GetCurrentPace(), time.Now().Format("15:04:05")),
	})
	return
}

func (state *StateRunning) Stop(bot *ubot.Bot) (err error) {
	messageId, _ := state.parent.statusMessage.GetInteger("message_id")
	bot.EditMessageText(axon.O{
		"chat_id":    state.parent.chatId,
		"message_id": messageId,
		"parse_mode": "MarkdownV2",
		"text":       fmt.Sprintf("Stopped at **%s**", time.Now().Format("15:04:05")),
	})
	state.parent.SetState(bot, &StateStopped{
		positions: state.positions,
	})
	return
}

func (state *StateRunning) GetGPX(ubot *ubot.Bot) (data []byte, err error) {
	return makeGpx(state.positions)
}

func (state *StateRunning) GetCurrentPace() (result *Pace) {
	ms := state.currentSpeed()
	if ms != -1 {
		seconds := 1000 / ms
		mins := math.Floor(seconds / 60)
		secs := math.Round(seconds - mins*60)
		result = &Pace{mins: int64(mins), secs: int64(secs)}
	}
	return
}

func (s *StateRunning) currentSpeed() (result float64) {
	if len(s.positions) < 2 {
		result = -1
		return
	}

	lastPos := s.positions[len(s.positions)-1]
	preLastPos := s.positions[len(s.positions)-2]
	meters := distance(preLastPos, lastPos) * 1000
	seconds := lastPos.when.Sub(preLastPos.when).Seconds()
	result = meters / seconds
	return
}

func distance(fromP, toP *Position) (result float64) {
	from := haversine.Coord{Lat: fromP.latitude, Lon: fromP.longitude}
	to := haversine.Coord{Lat: toP.latitude, Lon: toP.longitude}
	_, result = haversine.Distance(from, to)
	return
}
