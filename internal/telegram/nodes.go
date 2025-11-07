package telegram

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"go.uber.org/zap"
)

func (b *Bot) hNodeCmd(ctx *th.Context, u telego.Update) error {
	msg := b.reply(ctx, u.Message, "Processing...")
	defer b.delete(ctx, msg)

	parts := strings.Fields(u.Message.Text)
	if len(parts) == 0 {
		return nil
	}
	args := parts[1:]

	invalidNodeText := "You need to specify check-host node name in the first argument of the command. Avaliable formats: <code>us1</code>, <code>us1.node.check-host.net</code>"

	if len(args) != 1 {
		b.reply(ctx, u.Message, invalidNodeText)
		return nil
	}

	node := completeNode(args[0])
	if node == "" {
		b.reply(ctx, u.Message, invalidNodeText)
		return nil
	}

	nodes, err := b.CH.GetNodes()
	if err != nil {
		b.reply(ctx, u.Message, "An error occurred while retrieving the check result. You can report error using /report")
		b.Log.Error("failed to get nodes", zap.Error(err))
		return nil
	}

	if nodes[node].IP == "" {
		b.reply(ctx, u.Message, fmt.Sprintf("Node %s not found", node))
	}

	b.reply(ctx, u.Message, fmt.Sprintf(`Node info <code>%s</code>

ASN: %s
IP: %s
Location: %s
`, node, nodes[node].ASN, nodes[node].IP, strings.Join(nodes[node].Location, ", ")))

	return nil
}

func (b *Bot) hNodesCmd(ctx *th.Context, u telego.Update) error {
	msg := b.reply(ctx, u.Message, "Processing...")
	defer b.delete(ctx, msg)

	nodes, err := b.CH.GetNodes()
	if err != nil {
		b.reply(ctx, u.Message, "An error occurred while retrieving the check result. You can report error using /report")
		b.Log.Error("failed to get nodes", zap.Error(err))
		return nil
	}

	var textNodes string
	var page int = 1

	keys := make([]string, 0, len(nodes))
	for k := range nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := nodes[k]
		textNodes += fmt.Sprintf(
			"\n<code>%s</code> - %s, %s (<code>%s</code>)",
			strings.TrimSuffix(k, ".node.check-host.net"),
			v.Location[1],
			v.Location[2],
			v.IP,
		)
		if len(textNodes) > 2000 {
			b.reply(ctx, u.Message, fmt.Sprintf("Avaliable check-host nodes - page %d: %s", page, textNodes))
			page += 1
			textNodes = ""
		}
	}

	b.reply(ctx, u.Message, fmt.Sprintf("Avaliable check-host nodes - page %d: %s", page, textNodes))

	return nil
}
