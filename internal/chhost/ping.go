package chhost

import (
	"fmt"

	"github.com/d1manpro/checkhost"
)

func (c *ChHost) CheckPing(hostname string, maxNodes int, nodes []string) (checkhost.PingResult, string, error) {
	req, err := c.Client.CheckPing(checkhost.RequestData{
		Host:     hostname,
		MaxNodes: maxNodes,
		Nodes:    nodes,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to request check: %w", err)
	}

	res, err := c.Client.WaitPingResult(req.RequestID, c.Timeout, c.Interval)
	if err != nil {
		return nil, req.PermanentLink, fmt.Errorf("failed to parse check result: %w", err)
	}

	return res, req.PermanentLink, nil
}
