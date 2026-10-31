package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"

	"tg-bot/internal/config"
	"tg-bot/internal/mock"
	"tg-bot/internal/ui"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	seed, err := mock.LoadSeed("data/seed.json")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	be, err := mock.New(seed, "data/state.json")
	if err != nil {
		log.Fatal(err)
	}

	b, err := bot.New(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	ui.New(be).Register(b)

	log.Printf("Loaded %d categories, %d products", len(seed.Categories), len(seed.Products))
	log.Println("Bot started. Press Ctrl+C to stop.")
	b.Start(ctx)
}
