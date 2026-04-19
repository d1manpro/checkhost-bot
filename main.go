package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/d1manpro/checkhost-bot/internal/chhost"
	"github.com/d1manpro/checkhost-bot/internal/config"
	"github.com/d1manpro/checkhost-bot/internal/logger"
	"github.com/d1manpro/checkhost-bot/internal/telegram"
	"go.uber.org/zap"
)

func main() {
	cfgPath := flag.String("config", "config/", "path to config directory")
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	err := config.Load(*cfgPath, *debug)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	log, cleanup := logger.Init()
	defer cleanup()
	log.Info("\n\nStarting...\n")

	ch := chhost.New(log)

	tgBot, err := telegram.NewBot(log, config.Get(), ch)
	if err != nil {
		log.Fatal("failed to create telegram bot", zap.Error(err))
	}

	ctx := context.Background()
	err = tgBot.Start(ctx)
	if err != nil {
		log.Fatal("failed to start telegram bot", zap.Error(err))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Info("Stopping...")

	tgBot.Stop(ctx)
	log.Info("Telegram bot stopped")

	log.Info("Done.")
}
