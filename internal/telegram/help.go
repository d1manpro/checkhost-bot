package telegram

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func (b *Bot) hHelpCmd(ctx *th.Context, u telego.Update) error {
	b.reply(ctx, u.Message, b.Cfg.Messages.Help)
	return nil
}
