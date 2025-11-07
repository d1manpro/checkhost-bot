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

func (c *ChHost) CheckHttp(hostname string, maxNodes int, nodes []string) (checkhost.HTTPResult, string, error) {
	req, err := c.Client.CheckHTTP(checkhost.RequestData{
		Host:     hostname,
		MaxNodes: maxNodes,
		Nodes:    nodes,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to request check: %w", err)
	}

	res, err := c.Client.WaitHTTPResult(req.RequestID, c.Timeout, c.Interval)
	if err != nil {
		return nil, req.PermanentLink, fmt.Errorf("failed to parse check result: %w", err)
	}

	return res, req.PermanentLink, nil
}

func (c *ChHost) CheckTCP(hostname string, maxNodes int, nodes []string) (checkhost.TCPResult, string, error) {
	req, err := c.Client.CheckTCP(checkhost.RequestData{
		Host:     hostname,
		MaxNodes: maxNodes,
		Nodes:    nodes,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to request check: %w", err)
	}

	res, err := c.Client.WaitTCPResult(req.RequestID, c.Timeout, c.Interval)
	if err != nil {
		return nil, req.PermanentLink, fmt.Errorf("failed to parse check result: %w", err)
	}

	return res, req.PermanentLink, nil
}

func (c *ChHost) CheckUDP(hostname string, maxNodes int, nodes []string) (checkhost.UDPResult, string, error) {
	req, err := c.Client.CheckUDP(checkhost.RequestData{
		Host:     hostname,
		MaxNodes: maxNodes,
		Nodes:    nodes,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to request check: %w", err)
	}

	res, err := c.Client.WaitUDPResult(req.RequestID, c.Timeout, c.Interval)
	if err != nil {
		return nil, req.PermanentLink, fmt.Errorf("failed to parse check result: %w", err)
	}

	return res, req.PermanentLink, nil
}
