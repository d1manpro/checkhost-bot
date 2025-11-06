package telegram

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func (b *Bot) handleStartCommand(ctx *th.Context, u telego.Update) error {
	b.reply(ctx, u.Message, `Bot is started`)
	return nil
}
