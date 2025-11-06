package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/d1manpro/checkhost-bot/internal/chhost"
	"github.com/d1manpro/checkhost-bot/internal/config"
	"github.com/d1manpro/checkhost-bot/internal/telegram"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	log := setupLogger()
	defer log.Sync()
	log.Info("\n\nStarting...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}
	log.Info("Config succesfilly loaded")

	ch := chhost.New()

	tgBot, err := telegram.NewBot(log, cfg, ch)
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

func setupLogger() *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:     "time",
		LevelKey:    "level",
		MessageKey:  "msg",
		EncodeTime:  zapcore.TimeEncoderOfLayout("2006.01.02 15:04:05.000"),
		EncodeLevel: zapcore.CapitalLevelEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)

	var cores []zapcore.Core
	cores = append(cores, consoleCore)

	logFile, err := os.OpenFile("bot.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("cannot open log file: %v", err))
	}
	fileCore := zapcore.NewCore(encoder, zapcore.AddSync(logFile), zapcore.InfoLevel)
	cores = append(cores, fileCore)

	return zap.New(zapcore.NewTee(cores...))
}
