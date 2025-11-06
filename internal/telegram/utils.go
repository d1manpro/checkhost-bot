package telegram

import (
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"go.uber.org/zap"
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

func parseCommandArgs(m *telego.Message) (args []string) {
	if m.Text == "" {
		return nil
	}

	parts := strings.Fields(m.Text)
	if len(parts) == 0 {
		return nil
	}

	if len(parts) > 1 {
		args = parts[1:]
	}

	return args
}
