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

type StateTracking struct {
	StateBase
	positions []*Position
}

func (state *StateTracking) EnterState(bot *ubot.Bot, message axon.O) (err error) {
	bot.UnpinAllChatMessages(axon.O{
		"chat_id": state.parent.chatId,
	})
	_, err = bot.SendMessage(axon.O{
		"chat_id":    state.parent.chatId,
		"text":       "Tracking **started** at " + time.Now().Format("15:04:05"),
		"parse_mode": "MarkdownV2",
		"reply_markup": axon.O{
			"keyboard": axon.A{
				axon.A{
					axon.O{
						"text": "End",
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
	})

	if pinnedMessage, err := bot.SendMessage(axon.O{
		"chat_id": state.parent.chatId,
		"text":    "State: **started**, Pace: --:--",
	}); err == nil {
		messageId, _ := pinnedMessage.GetInteger("message_id")
		bot.PinChatMessage(axon.O{
			"chat_id":    state.parent.chatId,
			"message_id": messageId,
		})
		state.parent.pinnedMessage = pinnedMessage
	}
	return
}

func (state *StateTracking) BeginTracking(bot *ubot.Bot, position *Position) (err error) {
	bot.SendMessage(axon.O{
		"chat_id": state.parent.chatId,
		"text":    "Current tracking aborted, now restarting...",
	})
	err = state.parent.SetState(
		bot,
		&StateTracking{
			StateBase: StateBase{
				parent: state.parent,
			},
			positions: []*Position{position},
		},
		nil)
	return
}

func (state *StateTracking) PauseTracking(bot *ubot.Bot) (err error) {
	state.parent.SetState(bot, &StatePaused{
		StateBase: StateBase{
			parent: state.parent,
		},
		positions: state.positions,
	},
		nil,
	)
	return
}

func (state *StateTracking) UpdateTracking(bot *ubot.Bot, position *Position) (err error) {
	if position == nil {
		log.Fatalf("null position")
	}
	state.positions = append(state.positions, position)

	currentPace := state.GetCurrentPace()
	pinnedId, _ := state.parent.pinnedMessage.GetInteger("message_id")
	bot.EditMessageText(axon.O{
		"chat_id":    state.parent.chatId,
		"message_id": pinnedId,
		"text":       fmt.Sprintf("State: **tracking**, Pace: %s", currentPace),
	})
	return
}

func (state *StateTracking) EndTracking(bot *ubot.Bot) (err error) {
	bot.SendMessage(axon.O{
		"chat_id":    state.parent.chatId,
		"parse_mode": "MarkdownV2",
		"text":       fmt.Sprintf("Stopped at **%s**", time.Now().Format("15:04:05")),
	})
	state.parent.SetState(bot, &StateStopped{
		StateBase: StateBase{
			parent: state.parent,
		},
		positions: state.positions,
	},
		nil)
	return
}

func (state *StateTracking) GetGPX(ubot *ubot.Bot, matchType string) (data []byte, err error) {
	return makeGpx(state.positions, matchType)
}

func (state *StateTracking) GetCurrentPace() (result *Pace) {
	ms := state.currentSpeed()
	if ms != -1 {
		seconds := 1000 / ms
		mins := math.Floor(seconds / 60)
		secs := math.Round(seconds - mins*60)
		result = &Pace{mins: int64(mins), secs: int64(secs)}
	}
	return
}

func (s *StateTracking) currentSpeed() (result float64) {
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
