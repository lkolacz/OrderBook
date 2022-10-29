package core

import (
	"github.com/lkolacz/OrderBook/rest/config"
)

type AppCore interface {
	// ProcessErrand is an function that storag incoming orders and process it
	ProcessErrand(string) error
}

type app struct {
	cfg *config.Base
}

func New(cfg *config.Base) (*app, error) {

	return &app{
		cfg: cfg,
	}, nil
}
