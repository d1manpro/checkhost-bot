package chhost

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/d1manpro/checkhost"
	"go.uber.org/zap"
)

var ErrInvalidNodes = errors.New("invalid check-host nodes")

type ChHost struct {
	Client   *checkhost.Client
	Timeout  time.Duration
	Interval time.Duration
}

func New(log *zap.Logger) *ChHost {
	return &ChHost{
		Client: checkhost.New(checkhost.WithHTTPClient(&http.Client{
			Timeout: 10 * time.Second,
			Transport: &loggingRoundTripper{
				logger: log,
				rt: &http.Transport{
					DialContext: (&net.Dialer{
						Timeout:   5 * time.Second,
						KeepAlive: 30 * time.Second,
					}).DialContext,
					TLSHandshakeTimeout: 5 * time.Second,
					MaxIdleConns:        100,
					IdleConnTimeout:     90 * time.Second,
				},
			},
		})),
		Timeout:  15 * time.Second,
		Interval: 1 * time.Second,
	}
}

type loggingRoundTripper struct {
	logger *zap.Logger
	rt     http.RoundTripper
}

func (l *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := l.rt.RoundTrip(req)
	duration := time.Since(start)

	if err != nil {
		l.logger.Error("HTTP request failed",
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		return nil, err
	}

	l.logger.Info("HTTP request",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Int("status", resp.StatusCode),
		zap.Duration("duration", duration),
	)
	return resp, nil
}
