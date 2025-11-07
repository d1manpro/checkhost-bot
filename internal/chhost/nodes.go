package chhost

import "github.com/d1manpro/checkhost"

func (ch *ChHost) GetNodes() (map[string]checkhost.Node, error) {
	return ch.Client.GetNodes()
}
