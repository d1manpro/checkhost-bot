package telegram

import (
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"go.uber.org/zap"
)

func (b *Bot) hTcpCmd(ctx *th.Context, u telego.Update) error {
	msg := b.reply(ctx, u.Message, "Processing...")
	defer b.delete(ctx, msg)

	args, err := b.parseArgs(u.Message.Text)
	if err != nil {
		var text string
		switch err {
		case ErrEmpty:
			text = b.Cfg.Messages.Usage["tcp"]
		case ErrInvalidTarget:
			text = "Specified target isn't valid IP address or domain"
		}
		b.reply(ctx, u.Message, text)
		return nil
	}

	res, link, err := b.CH.CheckTCP(args.Target, args.MaxNodes, args.Nodes)
	if err != nil {
		b.reply(ctx, u.Message, "An error occurred while retrieving the check result. You can report error using /report")
		b.Log.Error("failed to check tcp", zap.String("target", args.Target), zap.String("link", link), zap.Error(err))
		return nil
	}

	var textResult string
	for n, d := range res {
		if d.Error != "" {
			d.Error = fmt.Sprintf("\n    Error: <b>%s</b>", d.Error)
		}
		textResult += fmt.Sprintf(`
Node <code>%s</code>
    Time: <b>%.3f</b> ms
    IP: <code>%s</code>%s
`, strings.TrimSuffix(n, ".node.check-host.net"), d.Time*1000, d.IP, d.Error)

		if len(textResult) > 4000 {
			b.reply(ctx, u.Message, fmt.Sprintf(`TCP check <code>%s</code>
<a href="%s">Check result</a>
%s`, args.Target, link, textResult))

			textResult = ""
		}
	}

	b.reply(ctx, u.Message, fmt.Sprintf(`TCP check <code>%s</code>
<a href="%s">Check result</a>
%s`, args.Target, link, textResult))

	return nil
}
