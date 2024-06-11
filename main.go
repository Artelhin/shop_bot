package main

import (
	"shop_bot/bot"
	"shop_bot/log"

	"shop_bot/config"
)

func main() {
	log.Info("startup...")
	defer log.Info("app closed")

	cfg, err := config.FromFile("./config.yaml")
	if err != nil {
		log.Error("can't configure: %s", err)
		return
	}
	log.Info("config: %v", cfg)

	err = bot.RunBot(cfg)
	if err != nil {
		log.Error("bot failed", err)
	}
}
