package telegram

import (
	"context"
	"fmt"

	"github.com/d1manpro/checkhost-bot/internal/chhost"
	"github.com/d1manpro/checkhost-bot/internal/config"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Bot struct {
	Bot *telego.Bot
	Log *zap.Logger
	Cfg *config.Config
	CH  *chhost.ChHost
}

type ErrData struct {
	UserID int64
	Text   string
	Error  error
}

func NewBot(log *zap.Logger, cfg *config.Config, ch *chhost.ChHost) (*Bot, error) {
	telegoBot, err := telego.NewBot(cfg.Token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Bot: telegoBot,
		Log: log,
		Cfg: cfg,
		CH:  ch,
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	fullWhURL := b.Cfg.Webhook.URL + b.Cfg.Webhook.Path

	srv := &fasthttp.Server{}

	updates, err := b.Bot.UpdatesViaWebhook(ctx,
		telego.WebhookFastHTTP(srv, b.Cfg.Webhook.Path, b.Bot.SecretToken()),
		telego.WithWebhookSet(ctx,
			&telego.SetWebhookParams{
				URL:            fullWhURL,
				SecretToken:    b.Bot.SecretToken(),
				AllowedUpdates: []string{"message", "callback_query"},
			},
		),
	)
	if err != nil {
		return fmt.Errorf("failed to setup webhook: %w", err)
	}

	bh, err := th.NewBotHandler(b.Bot, updates)
	if err != nil {
		return fmt.Errorf("failed to create BotHandler: %w", err)
	}

	bh.Use(func(ctx *th.Context, u telego.Update) error {
		if u.Message != nil {
			b.Log.Info("handling message", zap.Int("updateID", u.UpdateID), zap.String("text", u.Message.Text))
		} else if u.CallbackQuery != nil {
			b.Log.Info("handling callback-query", zap.Int("updateID", u.UpdateID), zap.String("data", u.CallbackQuery.Data))
		}
		if u.Message != nil && u.Message.Chat.Type != "private" {
			b.Log.Info("chat_type != private", zap.Int("updateID", u.UpdateID), zap.String("text", u.Message.Text))
			b.reply(ctx, u.Message, "Bot avaliable only in private chats")
			return nil
		}
		return ctx.Next(u)
	})

	b.initHandlers(bh)

	go func() {
		err = bh.Start()
		if err != nil {
			b.Log.Fatal("failed to start BotHandler", zap.Error(err))
		}
	}()

	go func() {
		err := srv.ListenAndServe(":" + b.Cfg.Webhook.Port)
		if err != nil {
			b.Log.Fatal("server error", zap.Error(err))
		}
	}()
	b.Log.Info("Webhook server started", zap.String("url", fullWhURL), zap.String("port", b.Cfg.Webhook.Port))
	return nil
}

func (b *Bot) Stop(ctx context.Context) error {
	err := b.Bot.DeleteWebhook(ctx, &telego.DeleteWebhookParams{
		DropPendingUpdates: true,
	})
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}
	return nil
}

func (b *Bot) initHandlers(bh *th.BotHandler) {
	bh.Handle(b.hStartCmd, th.CommandEqual("start"))
	bh.Handle(b.hHelpCmd, th.CommandEqual("help"))

	bh.Handle(b.hPingCmd, th.CommandEqual("ping"))
	bh.Handle(b.hHttpCmd, th.CommandEqual("http"))
	bh.Handle(b.hTcpCmd, th.CommandEqual("tcp"))
	bh.Handle(b.hUdpCmd, th.CommandEqual("udp"))
	bh.Handle(b.hDnsCmd, th.CommandEqual("dns"))

	bh.Handle(b.hNodeCmd, th.CommandEqual("node"))

	bh.Handle(b.handleAnyMessage, th.AnyMessage())
}

func (b *Bot) handleAnyMessage(ctx *th.Context, u telego.Update) error {
	return nil
}

func (b *Bot) hStartCmd(ctx *th.Context, u telego.Update) error {
	b.reply(ctx, u.Message, b.Cfg.Messages.Start)
	return nil
}

func (b *Bot) hHelpCmd(ctx *th.Context, u telego.Update) error {
	b.reply(ctx, u.Message, b.Cfg.Messages.Help)
	return nil
}
