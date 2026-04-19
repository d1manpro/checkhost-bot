package telegram

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/d1manpro/checkhost-bot/internal/chhost"
	"github.com/d1manpro/checkhost-bot/internal/config"
	"golang.org/x/net/proxy"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
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
	var httpClient *http.Client
	if cfg.Bot.Proxy != "" {
		u, err := url.Parse(cfg.Bot.Proxy)
		if err != nil {
			return nil, err
		}

		log.Info("Using SOCKS5 proxy", zap.String("host", u.Host))

		var auth *proxy.Auth
		if u.User != nil {
			password, _ := u.User.Password()
			auth = &proxy.Auth{
				User:     u.User.Username(),
				Password: password,
			}
		}

		dialer, err := proxy.SOCKS5("tcp", u.Host, auth, proxy.Direct)
		if err != nil {
			return nil, err
		}

		transport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}

		httpClient = &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}
	}

	var telegoBot *telego.Bot
	var err error

	if httpClient != nil {
		telegoBot, err = telego.NewBot(
			cfg.Bot.Token,
			telego.WithHTTPClient(httpClient),
		)
	} else {
		telegoBot, err = telego.NewBot(cfg.Bot.Token)
	}
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
	var err error
	var m string
	if !b.Cfg.Bot.Webhook.Enabled {
		err = b.startLongPolling(ctx)
		m = "long-polling"
	} else {
		err = b.startWebhook(ctx)
		m = "webhook"
	}

	if err == nil {
		_, err = b.Bot.SendMessage(ctx, tu.Message(
			tu.ID(b.Cfg.Bot.AdminID),
			"[system] bot running via "+m,
		).WithParseMode("HTML"))
		if err != nil {
			b.Log.Error("failed to send bot-start message", zap.Error(err))
		}
	}
	return err
}

func (b *Bot) startWebhook(ctx context.Context) error {
	fullWhURL := b.Cfg.Bot.Webhook.URL + b.Cfg.Bot.Webhook.Path

	srv := &fasthttp.Server{}

	updates, err := b.Bot.UpdatesViaWebhook(ctx,
		telego.WebhookFastHTTP(srv, b.Cfg.Bot.Webhook.Path, b.Bot.SecretToken()),
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

	b.setupMiddleware(bh)
	b.initHandlers(bh)

	go func() {
		err = bh.Start()
		if err != nil {
			b.Log.Fatal("failed to start BotHandler", zap.Error(err))
		}
	}()

	go func() {
		err := srv.ListenAndServe(":" + b.Cfg.Bot.Webhook.Port)
		if err != nil {
			b.Log.Fatal("server error", zap.Error(err))
		}
	}()
	b.Log.Info("Webhook server started", zap.String("url", fullWhURL), zap.String("port", b.Cfg.Bot.Webhook.Port))
	return nil
}

func (b *Bot) startLongPolling(ctx context.Context) error {
	updates, err := b.Bot.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{
		AllowedUpdates: []string{"message", "callback_query"},
		Timeout:        20,
	})
	if err != nil {
		return fmt.Errorf("failed to start long polling: %w", err)
	}

	bh, err := th.NewBotHandler(b.Bot, updates)
	if err != nil {
		return fmt.Errorf("failed to create BotHandler: %w", err)
	}

	b.setupMiddleware(bh)
	b.initHandlers(bh)

	go func() {
		if err := bh.Start(); err != nil {
			b.Log.Fatal("handler error", zap.Error(err))
		}
	}()

	b.Log.Info("Long polling started")
	return nil
}

func (b *Bot) setupMiddleware(bh *th.BotHandler) {
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
	bh.Handle(b.hNodesCmd, th.CommandEqual("nodes"))

	bh.Handle(b.handleAnyMessage, th.AnyMessage())
}

func (b *Bot) Stop(ctx context.Context) error {
	err := b.Bot.DeleteWebhook(ctx, &telego.DeleteWebhookParams{
		DropPendingUpdates: true,
	})
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}
	_, err = b.Bot.SendMessage(ctx, tu.Message(
		tu.ID(b.Cfg.Bot.AdminID),
		"[system] bot stopped",
	).WithParseMode("HTML"))
	if err != nil {
		b.Log.Error("failed to send bot-start message", zap.Error(err))
	}

	return nil
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
