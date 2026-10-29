package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"

	"shopbot/internal/backend"
	"shopbot/internal/config"
	"shopbot/internal/mock"
	"shopbot/internal/ui"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal(err)
	}

	seed, err := mock.LoadSeed("data/seed.json")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var be backend.Backend = mock.New(seed)

	b, err := bot.New(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	ui.New(be).Register(b)

	log.Printf("Loaded %d categories, %d products", len(seed.Categories), len(seed.Products))
	log.Println("Bot started. Press Ctrl+C to stop.")
	b.Start(ctx)
}
