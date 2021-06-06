package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/sdurz/axon"
	"github.com/sdurz/ubot"
)

var (
	apiKey         string
	webhookUrl     string
	graphHopperUrl string
	serverBind     string
	signals        chan os.Signal

	lrucache *lru.Cache
)

func init() {
	flag.StringVar(&apiKey, "apiKey", "", "api key")
	flag.StringVar(&webhookUrl, "webhookUrl", "", "webhook url")
	flag.StringVar(&graphHopperUrl, "graphHopperUrl", "http://localhost:8989", "graphhopper url")
	flag.StringVar(&serverBind, "serverBind", "", "server:port")

	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var err error
	if lrucache, err = lru.New(20000); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var (
		wg  sync.WaitGroup
		bot *ubot.Bot
	)

	flag.Parse()
	bot = ubot.NewBot(&ubot.Configuration{
		APIToken:   apiKey,
		WebhookUrl: webhookUrl,
		ServerBind: serverBind,
	})
	bot.AddMessageHandler(ubot.MessageHasCommand("start"), StartCommandHandler)
	bot.AddMessageHandler(ubot.MessageHasCommand("stop"), StopCommandHandler)
	bot.AddMessageHandler(ubot.And(ubot.MessageIsPrivate, ubot.MessageHasLocation), MessagePositionHandler)
	bot.AddEditedMessageHandler(ubot.And(ubot.MessageIsPrivate, ubot.MessageHasLocation), MessagePositionUpdateHandler)
	bot.AddMessageHandler(ubot.MessageHasCommand("getgpx"), GetGpxCommandHandler)
	bot.AddMessageHandler(ubot.MessageHasCommand("pause"), PauseTrackingCommandHandler)
	bot.AddMessageHandler(ubot.MessageHasCommand("resume"), ResumeTrackingCommandHandler)
	bot.AddMessageHandler(ubot.MessageHasCommand("end"), EndTrackingCommandHandler)
	bot.AddMessageHandler(ubot.MessageHasCommand("setprofile"), SetProfileCommandHandler)
	bot.AddMessageHandler(ubot.MessageHasCommand("help"), GetHelpCommandHandler)
	bot.AddMessageHandler(ubot.MessageIsPrivate, CommandMessageHandler)
	bot.AddCallbackQueryHandler(ubot.Always, CallbackQueryHandler)

	ctx, cancel := context.WithCancel(context.Background())
	updatesSource := ubot.ServerSource
	if serverBind == "" {
		if rr, err := bot.DeleteWebhook(axon.O{"drop_pending_updates": true}); err != nil || !rr {
			log.Fatal("can't delete webhooks")
		}
		updatesSource = ubot.GetUpdatesSource
	}

	go bot.Forever(ctx, &wg, updatesSource)
	wg.Add(1)
	<-signals

	cancel()
	wg.Wait()
	log.Println("done with main")
}
