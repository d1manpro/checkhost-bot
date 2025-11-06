package chhost

import (
	"errors"
	"time"

	"github.com/d1manpro/checkhost"
)

var ErrInvalidNodes = errors.New("invalid check-host nodes")

type ChHost struct {
	Client   *checkhost.Client
	Timeout  time.Duration
	Interval time.Duration
}

func New() *ChHost {
	return &ChHost{
		Client:   checkhost.New(),
		Timeout:  15 * time.Second,
		Interval: 1 * time.Second,
	}
}
