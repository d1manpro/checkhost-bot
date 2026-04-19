package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init() (*zap.Logger, func()) {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:     "time",
		LevelKey:    "level",
		MessageKey:  "msg",
		EncodeTime:  zapcore.TimeEncoderOfLayout("2006.01.02 15:04:05.000"),
		EncodeLevel: zapcore.CapitalLevelEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	os.MkdirAll("logs", 0755)

	logFile, err := os.OpenFile("logs/latest.log", os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("cannot open log file: %v", err))
	}

	consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
	fileCore := zapcore.NewCore(encoder, zapcore.AddSync(logFile), zapcore.InfoLevel)
	logger := zap.New(zapcore.NewTee(consoleCore, fileCore))

	cleanup := func() {
		fmt.Println("Saving log...")

		logger.Sync()
		logFile.Close()

		ts := time.Now().Format("2006-01-02_15-04-05")
		dstName := filepath.Join("logs", fmt.Sprintf("%s.log.gz", ts))

		src, err := os.Open("logs/latest.log")
		if err != nil {
			fmt.Printf("cannot open latest.log: %v\n", err)
			return
		}
		defer src.Close()

		dst, err := os.Create(dstName)
		if err != nil {
			fmt.Printf("cannot create %s: %v\n", dstName, err)
			return
		}
		defer dst.Close()

		gz := gzip.NewWriter(dst)
		_, err = io.Copy(gz, src)
		if err != nil {
			fmt.Printf("error compressing log: %v\n", err)
		}
		gz.Close()

		fmt.Printf("Saved log into logs/%s.log.gz\n", ts)
	}

	log = logger

	return logger, cleanup
}

func Get() *zap.Logger {
	if log == nil {
		panic("logger not loaded")
	}
	return log
}
