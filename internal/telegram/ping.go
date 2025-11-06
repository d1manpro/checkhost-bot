package telegram

import (
	"fmt"
	"regexp"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"go.uber.org/zap"
)

func (b *Bot) hPingCmd(ctx *th.Context, u telego.Update) error {
	msg := b.reply(ctx, u.Message, "Processing...")
	defer b.delete(ctx, msg)

	args := parseCommandArgs(u.Message)
	if len(args) == 0 {
		b.reply(ctx, u.Message, "You need to specify target (a valid IP address or domain) in the first argument of the command.")
		return nil
	} else {
		if !valid.IsIP(args[0]) && !valid.IsDNSName(args[0]) {
			b.reply(ctx, u.Message, "Specified target isn't valid IP address or domain")
			return nil
		}
	}

	var maxNodes int
	var err error
	var nodes []string
	re := regexp.MustCompile(`^[a-z]{2}[0-9](?:\.node\.check-host\.net)?$`)

	if len(args) > 1 {
		for _, n := range args[1:] {
			if !re.MatchString(n) {
				b.reply(ctx, u.Message, "One of specified check-host nodes is invalid domain name. Check available nodes: /nodes")
				return nil
			}
			if len(n) == 3 {
				n = n + ".node.check-host.net"
			}
			nodes = append(nodes, n)
		}
	}

	res, link, err := b.CH.CheckPing(args[0], maxNodes, nodes)
	if err != nil {
		b.reply(ctx, u.Message, "An error occurred while retrieving the check result. You can report error using /report")
		b.Log.Error("failed to check ping", zap.String("target", args[0]), zap.String("link", link), zap.Error(err))
		return nil
	}

	var textResult string
	for n, d := range res {
		textResult += fmt.Sprintf(`
Node <code>%s</code>
    Result: <b>%s</b>
    Time: <b>%.3f</b> ms
    IP: <code>%s</code>
`, strings.TrimSuffix(n, ".node.check-host.net"), d.Status, d.Time*1000, d.IP)

		if len(textResult) > 4000 {
			b.reply(ctx, u.Message, fmt.Sprintf(`Ping server <b><code>%s</code></b>
<a href="%s">Check result</a>
%s`, args[0], link, textResult))

			textResult = ""
		}
	}

	b.reply(ctx, u.Message, fmt.Sprintf(`Ping server <b><code>%s</code></b>
<a href="%s">Check result</a>
%s`, args[0], link, textResult))

	return nil
}
