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
	if message, err := bot.SendMessage(axon.O{
		"chat_id":    state.parent.chatId,
		"text":       "Tracking **started** at " + time.Now().Format("15:04:05"),
		"parse_mode": "MarkdownV2",
		"reply_markup": axon.O{
			"keyboard": axon.A{
				axon.A{
					axon.O{
						"text": "Stop",
					},
				},
				axon.A{
					axon.O{
						"text": "Pause",
					},
				},
			},
			"resize_keyboard": true,
		},
	}); err == nil {
		state.parent.statusMessage = message
	}
	return
}

func (state *StateRunning) Start(bot *ubot.Bot, position *Position) (err error) {
	bot.SendMessage(axon.O{
		"chat_id": state.parent.chatId,
		"text":    "Current tracking aborted, now restarting...",
	})
	err = state.parent.SetState(
		bot,
		&StateRunning{
			parent:    state.parent,
			positions: []*Position{position},
		},
	)
	return
}

func (state *StateRunning) Pause(bot *ubot.Bot) (err error) {
	state.parent.SetState(bot, &StatePaused{
		parent:    state.parent,
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

	chatId, _ := state.parent.statusMessage.GetInteger("chat.id")
	bot.SendMessage(axon.O{
		"chat_id": chatId,
		"text":    fmt.Sprintf("Current pace: **%s**\nUpdated: %s", state.GetCurrentPace(), time.Now().Format("15:04:05")),
	})
	return
}

func (state *StateRunning) Stop(bot *ubot.Bot) (err error) {
	bot.SendMessage(axon.O{
		"chat_id":    state.parent.chatId,
		"parse_mode": "MarkdownV2",
		"text":       fmt.Sprintf("Stopped at **%s**", time.Now().Format("15:04:05")),
	})
	state.parent.SetState(bot, &StateStopped{
		parent:    state.parent,
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
