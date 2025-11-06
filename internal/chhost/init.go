package chhost

import "github.com/d1manpro/checkhost"

type ChHost struct {
	Client *checkhost.Client
}

func New() *ChHost {
	return &ChHost{
		Client: checkhost.New(),
	}
}
