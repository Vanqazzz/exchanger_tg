package main

import (
	"main/internal/config"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"

	"go.uber.org/zap"
)

func Init(log *zap.Logger) (*telego.Bot, error) {
	_ = godotenv.Load()
	cfg := config.Must(config.NewFromEnv())

	bot, err := telego.NewBot(cfg.Token)
	if err != nil {
		log.Error("new bot", zap.Error(err))
		os.Exit(1)
	}

	return bot, nil
}
