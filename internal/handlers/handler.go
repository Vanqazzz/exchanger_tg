package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"main/internal/config"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"go.uber.org/zap"
)

type App struct {
	bh  *th.BotHandler
	log *zap.Logger
	cfg *config.Config
}

func StartHandler(bh *th.BotHandler, log *zap.Logger) {

	// Start handler
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			fmt.Sprint("Привіт! Це бот з актуальним курсом валют."),
		).WithReplyMarkup(
			tu.Keyboard(tu.KeyboardRow(

				tu.KeyboardButton("🇺🇦 UAH"),
				tu.KeyboardButton("🇨🇿 CZK"),
				tu.KeyboardButton("💰 Crypto"),
			),
			).WithResizeKeyboard()))

		return nil
	}, th.CommandEqual("start"))

	// UAH handler ...
	bh.Handle(func(ctx *th.Context, update telego.Update) error {

		e, err := getExchangeRates("usd.min", log)
		if err != nil {
			log.Error("Fail to get rates", zap.Error(err))
			return err
		}

		temp := e["usd"].(map[string]interface{})

		USD_Rate := temp["uah"].(float64)

		e, err = getExchangeRates("eur.min", log)
		if err != nil {
			log.Error("Fail to get rates", zap.Error(err))
			return err
		}

		temp = e["eur"].(map[string]interface{})

		EUR_Rate := temp["uah"].(float64)

		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			fmt.Sprintf("\n🇺🇦 Курс UAH:\n\n🇺🇸 USD: %f\n🇪🇺 EUR: %f", USD_Rate, EUR_Rate),
		))
		log.Info("Rates printed")
		return nil

	}, th.TextEqual("🇺🇦 UAH"))

	// CZK handler ...
	bh.Handle(func(ctx *th.Context, update telego.Update) error {

		e, err := getExchangeRates("usd.min", log)
		if err != nil {
			log.Error("Fail to get rates", zap.Error(err))
			return err
		}
		temp := e["usd"].(map[string]interface{})

		USD_Rate := temp["czk"].(float64)

		e, err = getExchangeRates("eur.min", log)
		if err != nil {
			log.Error("Fail to get rates", zap.Error(err))
			return err

		}
		temp = e["eur"].(map[string]interface{})
		EUR_Rate := temp["czk"].(float64)

		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			fmt.Sprintf("\n🇨🇿 Курс CZK:\n\n🇺🇸 USD: %f\n🇪🇺 EUR: %f", USD_Rate, EUR_Rate),
		))
		log.Info("Rates printed")
		return nil

	}, th.TextEqual("🇨🇿 CZK"))

	// Crypto handler ...
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		// btc
		e, err := getExchangeRates("btc.min", log)
		if err != nil {
			log.Error("Fail to get rates", zap.Error(err))
			return err
		}

		temp := e["btc"].(map[string]interface{})

		BTC_TO_USD_Rate := temp["usd"].(float64)

		// eth
		e, err = getExchangeRates("eth.min", log)
		if err != nil {
			log.Error("Fail to get rates", zap.Error(err))
			return err
		}

		temp = e["eth"].(map[string]interface{})
		ETH_TO_USD_Rate := temp["usd"].(float64)

		// usdt
		e, err = getExchangeRates("usdt.min", log)
		if err != nil {
			log.Error("Fail to get rates", zap.Error(err))
			return err
		}

		temp = e["usdt"].(map[string]interface{})

		USDT_TO_USD_Rate := temp["usd"].(float64)

		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			fmt.Sprintf("\n⭐ Bitcoin       💲 USD: %f\n⭐ Ethereum 💲 USD: %f\n⭐ USDT          💲 USD: %f\n", BTC_TO_USD_Rate, ETH_TO_USD_Rate, USDT_TO_USD_Rate),
		))
		log.Info("Rates printed")
		return nil

	}, th.TextEqual("💰 Crypto"))

	// Update handler
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Дані оновлено.",
		))

		return nil
	}, th.CommandEqual("restart"))

	// Unknown command
	bh.Handle(func(ctx *th.Context, update telego.Update) error {

		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Unknown command, use /start",
		))
		return nil
	}, th.AnyCommand())

}

// getExchangeRates ...
func getExchangeRates(currency string, log *zap.Logger) (map[string]interface{}, error) {
	_ = godotenv.Load()
	cfg := config.Must(config.NewFromEnv())

	api := cfg.Api

	url := api + currency + ".json"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error Get response", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error to read body response", zap.Error(err))
		return nil, err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		log.Fatal("Fail unmarshal", zap.Error(err))
		return nil, err
	}

	return raw, nil
}
