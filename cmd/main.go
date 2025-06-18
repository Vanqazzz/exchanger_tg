package main

import (
	"context"
	"fmt"
	"main/internal/config"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Must(config.NewFromEnv())

	l, err := zap.NewDevelopment()

	bot, err := telego.NewBot(cfg.Token)
	if err != nil {
		l.Error("new bot", zap.Error(err))
		os.Exit(1)
	}

	l.Info("Bot started")

	updates, _ := bot.UpdatesViaLongPolling(context.Background(), nil)

	bh, _ := th.NewBotHandler(bot, updates)

	defer func() { _ = bh.Stop() }()

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			fmt.Sprintf("Hello %s!", update.Message.From.FirstName),
		))
		return nil
	}, th.CommandEqual("start"))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {

		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Unknown command, use /start",
		))
		return nil
	}, th.AnyCommand())

	_ = bh.Start()
}
