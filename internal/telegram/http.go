package telegram

import (
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"go.uber.org/zap"
)

func (b *Bot) hHttpCmd(ctx *th.Context, u telego.Update) error {
	msg := b.reply(ctx, u.Message, "Processing...")
	defer b.delete(ctx, msg)

	args, err := b.parseArgs(u.Message.Text)
	if err != nil {
		var text string
		switch err {
		case ErrEmpty:
			text = "You need to specify target (a valid IP address or domain) in the first argument of the command"
		case ErrInvalidTarget:
			text = "Specified target isn't valid IP address or domain"
		}
		b.reply(ctx, u.Message, text)
		return nil
	}

	res, link, err := b.CH.CheckHttp(args.Target, args.MaxNodes, args.Nodes)
	if err != nil {
		b.reply(ctx, u.Message, "An error occurred while retrieving the check result. You can report error using /report")
		b.Log.Error("failed to check ping", zap.String("target", args.Target), zap.String("link", link), zap.Error(err))
		return nil
	}

	var textResult string
	for n, d := range res {
		textResult += fmt.Sprintf(`
Node <code>%s</code>
    Code: <b>%s</b>
    Message: <b>%s</b>
    Time: <b>%.3f</b> ms
    IP: <code>%s</code>
`, strings.TrimSuffix(n, ".node.check-host.net"), d.Code, d.Message, d.Time*1000, d.IP)

		if len(textResult) > 4000 {
			b.reply(ctx, u.Message, fmt.Sprintf(`HTTP check <code>%s</code>
<a href="%s">Check result</a>
%s`, args.Target, link, textResult))

			textResult = ""
		}
	}

	b.reply(ctx, u.Message, fmt.Sprintf(`HTTP check <code>%s</code>
<a href="%s">Check result</a>
%s`, args.Target, link, textResult))

	return nil
}
