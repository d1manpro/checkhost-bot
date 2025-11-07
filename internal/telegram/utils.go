package telegram

import (
	"errors"
	"regexp"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"go.uber.org/zap"
)

type cmdArgs struct {
	Target   string
	Nodes    []string
	MaxNodes int
}

var (
	ErrEmpty         = errors.New("empty args")
	ErrInvalidTarget = errors.New("invalid target")
)

func (b *Bot) reply(ctx *th.Context, msg *telego.Message, text string) *telego.Message {
	msgThreadID := msg.MessageThreadID
	if msg.ReplyToMessage != nil {
		if msg.ReplyToMessage.MessageID == msg.MessageThreadID {
			msgThreadID = 0
		}
	}

	msg, err := b.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:          tu.ID(msg.Chat.ID),
		Text:            text,
		ParseMode:       "HTML",
		MessageThreadID: msgThreadID,
		LinkPreviewOptions: &telego.LinkPreviewOptions{
			IsDisabled: true,
		},
	})

	if err != nil {
		b.Log.Error("failed to send message", zap.Int64("chatID", msg.Chat.ID), zap.Error(err))
	}
	return msg
}

func (b *Bot) edit(ctx *th.Context, msg *telego.Message, text string) {
	_, err := b.Bot.EditMessageText(ctx, &telego.EditMessageTextParams{
		MessageID: msg.MessageID,
		ChatID:    msg.Chat.ChatID(),
		Text:      text,
		ParseMode: "HTML",
	})
	if err != nil {
		b.Log.Error("failed to edit message", zap.Int64("chatID", msg.Chat.ID), zap.Int("messageID", msg.MessageID), zap.Error(err))
	}
}

func (b *Bot) delete(ctx *th.Context, msg *telego.Message) {
	err := b.Bot.DeleteMessage(ctx, &telego.DeleteMessageParams{MessageID: msg.MessageID, ChatID: msg.Chat.ChatID()})
	if err != nil {
		b.Log.Error("failed to delete message", zap.Int64("chatID", msg.Chat.ID), zap.Int("messageID", msg.MessageID), zap.Error(err))
	}
}

func (b *Bot) parseArgs(msg string) (*cmdArgs, error) {
	parts := strings.Fields(msg)
	if len(parts) == 0 {
		return nil, ErrEmpty
	}
	args := parts[1:]

	if len(args) == 0 {
		return nil, ErrEmpty
	} else {
		if (!valid.IsIPv4(args[0]) && !valid.IsDNSName(args[0])) || !strings.Contains(args[0], ".") {
			return nil, ErrInvalidTarget
		}
	}

	var maxNodes int
	var nodes []string

	if len(args) > 1 {
		for _, n := range args[1:] {
			n := completeNode(n)
			if n != "" {
				nodes = append(nodes, n)
			}
		}
	}

	return &cmdArgs{
		Target:   args[0],
		Nodes:    nodes,
		MaxNodes: maxNodes,
	}, nil
}

func completeNode(n string) string {
	re := regexp.MustCompile(`^[a-z]{2}[0-9](?:\.node\.check-host\.net)?$`)
	if !re.MatchString(n) {
		return ""
	}
	if len(n) == 3 {
		n = n + ".node.check-host.net"
	}
	return n
}
