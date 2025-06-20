package main

import (
	"context"
	"main/internal/handlers"

	th "github.com/mymmrac/telego/telegohandler"
	"go.uber.org/zap"
)

func main() {

	log, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Error to enable logger")
		return
	}

	Bot, err := Init(log)
	if err != nil {
		log.Fatal("Fail to start bot")
		return
	}

	log.Info("Bot started")

	updates, _ := Bot.UpdatesViaLongPolling(context.Background(), nil)

	bh, _ := th.NewBotHandler(Bot, updates)

	defer func() { _ = bh.Stop() }()

	handlers.StartHandler(bh, log)

	_ = bh.Start()

}
